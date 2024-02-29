package core

import (
	"errors"
	"sync"
)

type ThreadPool struct {
	num_threads  int
	channel_size int

	jobs []Job

	wg sync.WaitGroup
}

func (this *ThreadPool) Init(num_threads int) {
	if num_threads <= 0 {
		err := errors.New("num threads <= 0")
		panic(err)
	}

	this.num_threads = num_threads

	this.jobs = make([]Job, 0)
}

func (this *ThreadPool) Enque(job Job) {
	this.wg.Add(1)
	this.jobs = append(this.jobs, job)
}

func (this *ThreadPool) Start() {
	for _, job := range this.jobs {
		go this.Dispatch(job)
	}
	this.wg.Wait()
}

func (this *ThreadPool) Dispatch(job Job) {
	defer this.wg.Done()
	job.Execute()
}
