package core

import (
	"sync"
)

type ThreadPool struct {
	jobs []Job

	wg sync.WaitGroup
}

func (this *ThreadPool) Init() {
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
