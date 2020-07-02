package worker

import (
	"runtime/debug"

	"github.com/astaxie/beego/logs"
)

type Worker struct {
	id         int
	ctrl       chan bool
	jobChannel chan Job
}

type Job interface {
	Do(int) bool
}

func NewWorker(id int, channelSize int) *Worker {
	return &Worker{
		id:         id,
		ctrl:       make(chan bool),
		jobChannel: make(chan Job, channelSize),
	}
}

func (this *Worker) GetChannel() *chan Job {
	return &this.jobChannel
}

func (this *Worker) Run() {
	go func() {
		for {
			select {
			case job := <-this.jobChannel:
				func() {
					defer func() {
						if r := recover(); r != nil {
							logs.Error("job panic, %s, stack: %s", r, string(debug.Stack()))
						}
					}()
					job.Do(this.id)
				}()
			case <-this.ctrl:
				return
			}
		}
	}()
}

func (this *Worker) Stop() {
	close(this.ctrl)
}
