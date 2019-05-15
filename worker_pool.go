package main

import (
	"io"
	"log"
	"net"
	"sync"
)

//reactor 模式中 处理具体业务事件的 协程池

type ConnHandler func(conn net.Conn)

type WorkerPool struct {
	workers     int
	maxTasks    int
	taskQueue   chan net.Conn
	connHandler ConnHandler

	mu     sync.Mutex
	closed bool
	done   chan struct{}
}

func NewWorkerPool(workerNum int, maxTaskNum int, connHandler ConnHandler) *WorkerPool {
	if workerNum < 1 || maxTaskNum < 1 ||
		connHandler == nil {
		panic("NewWokerPool Param error")

	}
	return &WorkerPool{
		workers:     workerNum,
		maxTasks:    t,
		taskQueue:   make(chan net.Conn, maxTaskNum),
		connHandler: connHandler,

		done: make(chan struct{}),
	}
}

func (p *WorkerPool) Close() {
	p.mu.Lock()
	p.closed = true
	close(p.done)
	close(p.taskQueue)
	p.mu.Unlock()
}

func (p *WorkerPool) addTask(conn net.Conn) {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()

	p.taskQueue <- conn
}

func (p *WorkerPool) start() {
	for i := 0; i < p.workers; i++ {
		go p.startWorker()
	}
}

func (p *WorkerPool) startWorker() {
	for {
		select {
		case <-p.done:
			return
		case conn := <-p.taskQueue:
			if conn != nil {
				func() {
					defer func() {
						if err := recover(); err != nil {
							//log.
						}

					}()
					p.connHandler(conn)

				}()
			}
		}
	}
}

func handleConn(conn net.Conn) {
	_, err := io.CopyN(conn, conn, 8)
	if err != nil {
		if err := epoller.Remove(conn); err != nil {
			log.Printf("failed to remove %v", err)
		}
		conn.Close()
	}
	opsRate.Mark(1)
}
