package go_workerpool

import (
	"log"
)

// Job interface which will be used to create a new job
type Job interface {
	Work() error
}

// Worker is the structure for worker
type Worker struct {
	id         int
	jobQueue   chan Job
	workerPool chan chan Job
	quitChan   chan bool
	started    bool
}

// NewWorker return a new instance of worker
func NewWorker(id int, workerPool chan chan Job) *Worker {
	return &Worker{
		id:         id,
		jobQueue:   make(chan Job),
		workerPool: workerPool,
		quitChan:   make(chan bool),
		started:    false,
	}
}

// Start worker
func (w *Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.workerPool <- w.jobQueue
			w.started = false

			select {
			case job := <-w.jobQueue:
				w.started = true

				if err := job.Work(); err != nil {
					log.Printf("error running worker %d: %s\n", w.id, err.Error())
				}

				w.started = false

			case <-w.quitChan:
				log.Printf("worker %d stopping\n", w.id)

				w.started = false

				return
			}
		}
	}()
}

// Stop worker
func (w *Worker) Stop() {
	go func() {
		w.quitChan <- true
	}()
}

// ID return worker id
func (w *Worker) ID() int {
	return w.id
}

// Started return worker status
func (w *Worker) Started() bool {
	return w.started
}