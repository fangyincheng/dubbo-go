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
	"runtime/debug"
	"sync"
	"time"
)

import (
	perrors "github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go/common/logger"
)

type Task func(workerID int)

type Worker struct {
	ID    int
	taskQ chan Task
	done  chan struct{}
	wg    sync.WaitGroup
	once  sync.Once
}

func NewWorker(id int, k *Pool) *Worker {
	w := &Worker{
		ID:    id,
		taskQ: make(chan Task),
		done:  make(chan struct{}),
	}

	w.wg.Add(1)
	go w.work(k)

	return w
}

func (w *Worker) work(k *Pool) {
	defer w.wg.Done()
	for {
		select {
		case t, ok := <-w.taskQ: // get task
			k.workerQ <- w
			if ok {
				func() {
					defer func() {
						if r := recover(); r != nil {
							logger.Errorf("worker ID:%d, panic error:%#v, debug stack:%s", w.ID, r, string(debug.Stack()))
						}
					}()

					t(w.ID) // call fn and send result
				}()
			} else {
				logger.Warnf("worker %d done channel closed, so it exits now with {its taskQ len = %d}",
					w.ID, len(w.taskQ))
				return
			}

		case <-w.done:
			logger.Warnf("worker %d done channel closed, so it exits now with {its taskQ len = %d}",
				w.ID, len(w.taskQ))
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

type Pool struct {
	workers   []*Worker
	workerNum int
	workerQ   chan *Worker
	taskQ     chan Task
	wg        sync.WaitGroup
	done      chan struct{}
	once      sync.Once
}

func NewPool(workerNum, maxWorkerNum int, taskNum int) *Pool {
	if maxWorkerNum == 0 {
		maxWorkerNum = workerNum
	}
	p := &Pool{
		workerNum: workerNum,
		workerQ:   make(chan *Worker, maxWorkerNum),
		taskQ:     make(chan Task, taskNum),
		done:      make(chan struct{}),
	}

	for i := 0; i < workerNum; i++ {
		w := NewWorker(i, p)
		p.workers = append(p.workers, w)
		p.workerQ <- w
	}

	p.wg.Add(1)
	go p.run()

	return p
}

func (p *Pool) run() {
	defer p.wg.Done()
	for {
		select {
		case t, ok := <-p.taskQ:
			if !ok {
				logger.Warn("pool taskQ has been closed")
				return
			}

			if err := p.PushTask(t, 0); err != nil {
				logger.Warnf("[Pool.run] a task is discard, error: %v", err)
			}

		case <-p.done:
			logger.Warnf("pool exit now while its task queue size = %d.", len(p.taskQ))
			return
		}
	}
}

func (p *Pool) PushTask(t Task, timeout time.Duration) error {
	if timeout == 0 {
		timeout = time.Second * 2
	}

	select {
	case w, ok := <-p.workerQ:
		if !ok {
			logger.Warn("pool workerQ has been closed")
			return perrors.New("pool workerQ has been closed")
		}
		w.taskQ <- t
		return nil
	case p.taskQ <- t:
		return nil
	case <-p.done:
		return perrors.New("Keeper has stopped!")
	case <-time.After(timeout):
		return perrors.New("Wait timeout.Queue is full.")
	}
}

func (p *Pool) Stop() {
	select {
	case <-p.done:
		return
	default:
		p.once.Do(func() {
			close(p.done)              // stop to get new task
			for i := range p.workers { // stop all workers
				p.workers[i].close()
			}
		})
	}
}

func (p *Pool) Close() {
	p.Stop()
	p.wg.Wait()
}
