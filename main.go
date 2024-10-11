package main

import (
	"fmt"
	"sync"
	"time"
)

// 单个任务结构体
type Task struct {
	id      int
	mu      sync.Mutex
	cond    *sync.Cond
	running bool
	stop    bool
}

// 创建新的任务
func NewTask(id int) *Task {
	t := &Task{
		id:      id,
		running: false,
		stop:    false,
	}
	t.cond = sync.NewCond(&t.mu)
	return t
}

// 启动单个任务
func (t *Task) Start() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.running {
		fmt.Printf("Task %d is already running\n", t.id)
		return
	}
	t.running = true

	go func() {
		fmt.Printf("Task %d thread started\n", t.id)
		for {
			t.mu.Lock()

			for !t.stop && !t.running {
				t.cond.Wait() // 线程等待信号
			}

			if t.stop {
				fmt.Printf("Task %d thread stopped\n", t.id)
				t.mu.Unlock()
				break
			}

			fmt.Printf("Task %d is running...\n", t.id)
			t.mu.Unlock()

			// 模拟任务执行
			time.Sleep(1 * time.Second)
		}
	}()
}

// 暂停单个任务
func (t *Task) Pause() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.running {
		fmt.Printf("Task %d is not running\n", t.id)
		return
	}

	t.running = false
	fmt.Printf("Task %d thread paused\n", t.id)
}

// 恢复单个任务
func (t *Task) Resume() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.running {
		fmt.Printf("Task %d is already running\n", t.id)
		return
	}

	t.running = true
	t.cond.Signal() // 发送信号，通知任务线程恢复
	fmt.Printf("Task %d thread resumed\n", t.id)
}

// 停止单个任务
func (t *Task) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.running && !t.stop {
		fmt.Printf("Task %d is not running\n", t.id)
		return
	}

	t.stop = true
	t.cond.Signal() // 唤醒线程并关闭它
	fmt.Printf("Task %d is being stopped\n", t.id)
}

// TaskManager 负责管理多个任务
type TaskManager struct {
	tasks []*Task
	mu    sync.Mutex
}

// 创建新的 TaskManager
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: []*Task{},
	}
}

// 添加任务到管理器
func (tm *TaskManager) AddTask(task *Task) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.tasks = append(tm.tasks, task)
}

// 启动所有任务
func (tm *TaskManager) StartAll() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, task := range tm.tasks {
		task.Start()
	}
}

// 暂停所有任务
func (tm *TaskManager) PauseAll() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, task := range tm.tasks {
		task.Pause()
	}
}

// 恢复所有任务
func (tm *TaskManager) ResumeAll() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, task := range tm.tasks {
		task.Resume()
	}
}

// 停止所有任务
func (tm *TaskManager) StopAll() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, task := range tm.tasks {
		task.Stop()
	}
}

func main() {
	// 创建TaskManager
	manager := NewTaskManager()

	// 创建多个任务
	task1 := NewTask(1)
	task2 := NewTask(2)
	task3 := NewTask(3)

	// 将任务添加到管理器
	manager.AddTask(task1)
	manager.AddTask(task2)
	manager.AddTask(task3)

	// 启动所有任务
	manager.StartAll()

	// 让任务运行一段时间
	time.Sleep(3 * time.Second)

	// 暂停所有任务
	manager.PauseAll()

	// 等待暂停一段时间
	time.Sleep(2 * time.Second)

	// 恢复所有任务
	manager.ResumeAll()

	// 再运行一段时间
	time.Sleep(3 * time.Second)

	// 停止所有任务
	manager.StopAll()

	// 等待所有任务停止
	time.Sleep(1 * time.Second)
}
