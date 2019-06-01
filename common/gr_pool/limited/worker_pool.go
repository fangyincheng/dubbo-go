// Copyright 2016-2019 Alex Stocks
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package limited

import (
	"container/heap"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

import (
	perrors "github.com/pkg/errors"
)

import (
	"github.com/dubbo/go-for-apache-dubbo/common/logger"
)

type Task func(workerID int)

type Worker struct {
	ID      int
	taskQ   chan Task
	pending int // count of pending tasks
	fin     int // count of finished tasks
	index   int // The index of the item in the heap.
	done    chan struct{}
	wg      sync.WaitGroup
	once    sync.Once
}

func NewWorker(id int, k *Keeper) *Worker {
	w := &Worker{
		ID:    id,
		taskQ: make(chan Task, 64),
		done:  make(chan struct{}),
	}

	w.wg.Add(1)
	go w.work(k)

	return w
}

func (w *Worker) work(k *Keeper) {
	defer w.wg.Done()
	for {
		select {
		case t, ok := <-w.taskQ: // get task from balancer
			if ok {
				func() {
					defer func() { // panic执行之前会保证defer被执行
						if r := recover(); r != nil {
							logger.Warnf("worker ID:%d, panic error:%#v, debug stack:%s", w.ID, r, string(debug.Stack()))
						}
					}()

					t(w.ID) // call fn and send result
				}()
				atomic.AddInt64(&k.finTaskNum, 1)
				k.workerQ <- w // we've finished w request
			} else {
				logger.Warnf("worker %d done channel closed, so it exits now with {its taskQ len = %d, pending = %d}",
					w.ID, len(w.taskQ), w.pending)
				return
			}

		case <-w.done:
			logger.Warnf("worker %d done channel closed, so it exits now with {its taskQ len = %d, pending = %d}",
				w.ID, len(w.taskQ), w.pending)
			return
		}
	}
}

func (w *Worker) stop() {
	select {
	case <-w.done:
		return

	default:
		w.once.Do(func() {
			close(w.done)
		})
	}
}

func (w *Worker) close() {
	w.stop()
	w.wg.Wait()
}

type WorkerPool []*Worker

func (p WorkerPool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p WorkerPool) Len() int {
	return len(p)
}

func (p WorkerPool) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].index = i
	p[j].index = j
}

func (p *WorkerPool) Push(w interface{}) {
	n := len(*p)
	item := w.(*Worker)
	item.index = n
	*p = append(*p, item)
}

func (p *WorkerPool) Pop() interface{} {
	arr := *p
	n := len(arr)
	item := arr[n-1]
	item.index = -1 // for safety
	*p = arr[0 : n-1]

	return item
}

type Keeper struct {
	workers    WorkerPool
	workerNum  int
	workerQ    chan *Worker
	taskQ      chan Task
	inTaskNum  int64
	finTaskNum int64
	wg         sync.WaitGroup
	done       chan struct{}
	once       sync.Once
}

func NewKeeper(workerNum, maxWorkerNum int, taskNum int) *Keeper {
	if maxWorkerNum == 0 {
		maxWorkerNum = workerNum
	}
	k := &Keeper{
		workers:   make(WorkerPool, 0, maxWorkerNum),
		workerNum: workerNum,
		workerQ:   make(chan *Worker, 128),
		taskQ:     make(chan Task, taskNum),
		done:      make(chan struct{}),
	}

	heap.Init(&k.workers)
	for i := 0; i < workerNum; i++ {
		heap.Push(&k.workers, NewWorker(i, k))
	}

	k.wg.Add(1)
	go k.run()

	return k
}

func (k *Keeper) run() {
	defer k.wg.Done()
	for {
		select {
		case t, ok := <-k.taskQ:
			if !ok {
				logger.Warn("keeper taskQ has been closed")
				return
			}

			if l := k.workers.Len(); l > 0 {
				worker := heap.Pop(&k.workers).(*Worker)
				worker.taskQ <- t
				worker.pending++
				heap.Push(&k.workers, worker)
			} else {
				// because Push & Pop is just invoked here. It is impossible that the workers is empty
				logger.Warn("failed to run task, this is impossible")
				k.taskQ <- t
			}
		case worker, ok := <-k.workerQ:
			if !ok {
				logger.Warn("keeper workerQ has been closed")
				return
			}

			worker.pending--
			worker.fin++
			heap.Remove(&k.workers, worker.index)
			heap.Push(&k.workers, worker)

		case <-k.done:
			logger.Warnf("keeper exit now while its task queue size = %d.", len(k.taskQ))
			return
		}
	}
}

func (k *Keeper) PushTask(t Task, timeout time.Duration) error {
	if timeout == 0 {
		timeout = time.Second * 2
	}
	select {
	case k.taskQ <- t:
		atomic.AddInt64(&k.inTaskNum, 1)
		return nil
	case <-k.done:
		return perrors.New("Keeper has stopped!")
	case <-time.After(timeout):
		return perrors.New("Wait timeout.Queue is full.")
	}
}

func (k *Keeper) PendingTaskNum() int {
	return int(atomic.LoadInt64(&k.inTaskNum) - atomic.LoadInt64(&k.finTaskNum))
}

func (k *Keeper) Stop() {
	select {
	case <-k.done:
		return
	default:
		k.once.Do(func() {
			close(k.done)              // stop to get new task
			for i := range k.workers { // stop all workers
				k.workers[i].close()
			}
		})
	}
}

func (k *Keeper) Close() {
	k.Stop()
	k.wg.Wait()
}
