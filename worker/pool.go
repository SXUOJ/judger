package worker

import (
	"errors"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

var (
	WorkPool             *Pool
	ErrPollAlreadyClosed = errors.New("Pool already closed")
)

type Pool struct {
	capacity       uint64
	runningWorkers uint64
	state          int64
	tasks          chan Task
	close          chan bool
	PanicHandler   func(interface{})
}

const (
	STOPED  = 0
	RUNNING = 1
)

func init() {
	WorkPool = NewPool(10)
}

func NewPool(capacity uint64) *Pool {
	return &Pool{
		capacity: capacity,
		state:    RUNNING,
		tasks:    make(chan Task, capacity),
		close:    make(chan bool),
	}
}

func (p *Pool) run() {
	p.incRunning()

	go func() {
		defer func() {
			p.decRunning()
			if r := recover(); r != nil {
				if p.PanicHandler != nil {
					p.PanicHandler(r)
				} else {
					logrus.Info("Worker panic: ", r)
				}
			}
		}()

		for {
			select {
			case task, ok := <-p.tasks:
				if !ok {
					return
				}
				task.Run()
			case <-p.close:
				return
			}
		}
	}()
}

func (p *Pool) Put(task Task) error {
	if p.state == STOPED {
		return ErrPollAlreadyClosed
	}

	if p.GetRunningWorkers() < p.GetCap() {
		p.run()
	}
	p.tasks <- task

	return nil
}

func (p *Pool) Close() {
	p.state = STOPED // 设置 state 为已停止

	// 阻塞等待所有任务被 worker 消费
	for len(p.tasks) > 0 {
	}

	p.close <- true // 发送销毁 worker 信号
	close(p.tasks)  // 关闭任务队列
}

func (p *Pool) incRunning() {
	atomic.AddUint64(&p.runningWorkers, 1)
}

func (p *Pool) decRunning() {
	atomic.AddUint64(&p.runningWorkers, ^uint64(0))
}

func (p *Pool) GetRunningWorkers() uint64 {
	return atomic.LoadUint64(&p.runningWorkers)
}

func (p *Pool) GetCap() uint64 {
	return atomic.LoadUint64(&p.capacity)
}
