package limited

import (
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestPool_PushTask(t *testing.T) {
	pool := NewPool(3, 5, 3)
	task := func(id int) {
		time.Sleep(time.Millisecond * time.Duration(10-id))
		t.Log("task: ", id)
	}
	task1 := func(id int) {
		time.Sleep(time.Millisecond * time.Duration(10-id))
		panic("task1")
	}
	for i := 0; i < 10; i++ {
		go func(i1 int) {
			var err error
			if i1%3 == 0 {
				err = pool.PushTask(task1, 0)
			} else {
				err = pool.PushTask(task, 0)
			}
			assert.NoError(t, err)
		}(i)
	}
	time.Sleep(time.Second * 3)
}
