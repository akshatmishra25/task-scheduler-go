package main

import (
	"fmt"
	"sync"
	"time"
)

type TaskScheduler struct {
	maxConcurrentTasks int
	runningTasks       int
	pendingTasks       chan func()
	mu                 sync.Mutex
}

func NewTaskScheduler(maxConcurrentTasks int) *TaskScheduler {
	return &TaskScheduler{
		maxConcurrentTasks: maxConcurrentTasks,
		pendingTasks:       make(chan func(), 100),
	}
}

func (ts *TaskScheduler) Schedule(task func()) {
	ts.pendingTasks <- task
	ts.runNextTask()
}

func (ts *TaskScheduler) runNextTask() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if ts.runningTasks >= ts.maxConcurrentTasks {
		return
	}

	select {
	case task := <-ts.pendingTasks:
		ts.runningTasks++
		go ts.executeTask(task)
	default:
	}
}

func (ts *TaskScheduler) executeTask(task func()) {
	defer func() {
		ts.mu.Lock()
		ts.runningTasks--
		ts.mu.Unlock()
		ts.runNextTask()
	}()

	task()
}

func main() {
	scheduler := NewTaskScheduler(3)

	task1 := func() {
		fmt.Println("Task 1 executing")
		time.Sleep(2 * time.Second)
	}

	task2 := func() {
		fmt.Println("Task 2 executing")
		time.Sleep(2 * time.Second)
	}

	task3 := func() {
		fmt.Println("Task 3 executing")
		time.Sleep(2 * time.Second)
	}

	task4 := func() {
		fmt.Println("Task 4 executing")
		time.Sleep(2 * time.Second)
	}

	// Schedule tasks
	scheduler.Schedule(task1)
	scheduler.Schedule(task2)
	scheduler.Schedule(task3)
	scheduler.Schedule(task4)

	// Wait for all tasks to complete
	time.Sleep(10 * time.Second)
}
