package workerpool

import (
	"log"
	"sync"
)

type WorkerPool struct {
	queue   chan func()
	workers int
	wg      sync.WaitGroup
}

// NewWorkerPool создает новый пул воркеров.
func NewWorkerPool(workers int, queueSize int) *WorkerPool {
	wp := &WorkerPool{
		queue:   make(chan func(), queueSize),
		workers: workers,
	}

	for i := 0; i < workers; i++ {
		go wp.worker()
	}

	return wp
}

// worker запускает воркер, который обрабатывает задачи из очереди.
func (wp *WorkerPool) worker() {
	for task := range wp.queue {
		task()
		wp.wg.Done() // Уведомляем, что задача выполнена
	}
}

// Submit добавляет задачу в очередь для обработки.
func (wp *WorkerPool) Submit(task func()) {
	wp.wg.Add(1)
	select {
	case wp.queue <- task:
		// Успешно добавлено в очередь
	default:
		log.Println("Queue is full, dropping task")
		wp.wg.Done() // Уведомляем, что задача не была выполнена
	}
}

// Wait дожидается завершения всех задач.
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

// Stop останавливает пул воркеров.
func (wp *WorkerPool) Stop() {
	close(wp.queue)
}
