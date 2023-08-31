package main_thread

import "sync"

// 主队列大小
const mainQSize = 2048

// 主队列
var mainQ = make(chan func(), mainQSize)

// 开始标记
var started = false
var startLocker = &sync.Mutex{}

// Process 处理任务,
// 只将任务添加到队列而不是马上执行...
func Process(task func()) {
	if nil == task {
		return
	}

	mainQ <- task

	if !started {
		startLocker.Lock()
		defer startLocker.Unlock()

		if !started {
			started = true
			go execute()
		}
	}
}

// 执行 task
func execute() {
	for {
		task := <-mainQ

		if nil != task {
			task()
		}
	}
}
