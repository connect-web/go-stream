package queue

import (
	"log"
	"sync"
	"time"

	pb "github.com/connect-web/go-stream/api" // Replace with your actual import path
)

type TaskQueue struct {
	tasks chan *pb.ScrapeRequest
	wg    sync.WaitGroup
}

func NewTaskQueue(bufferSize int) *TaskQueue {
	return &TaskQueue{
		tasks: make(chan *pb.ScrapeRequest, bufferSize),
	}
}

func (q *TaskQueue) Enqueue(task *pb.ScrapeRequest) {
	q.tasks <- task
	q.wg.Add(1)
}

func (q *TaskQueue) StartWorker() {
	go func() {
		for task := range q.tasks {
			log.Printf("Processing task: %s", task.Url)
			time.Sleep(1 * time.Second) // Simulated processing time
			log.Printf("Completed task: %s", task.Url)
			q.wg.Done()
		}
	}()
}

func (q *TaskQueue) Wait() {
	q.wg.Wait()
	close(q.tasks)
}
