package websocket

import (
	"sync"

	"github.com/cclehui/connection_platform/log"
	"github.com/gorilla/websocket"
)

//reactor 模式中 处理具体业务事件的 协程池

const (
	DEFAULT_WORKER_NUM   = 5
	DEFAULT_MAX_TASK_NUM = 1000
)

type ConnHandler func(conn *websocket.Conn)

type WorkerPool struct {
	workers     int
	maxTasks    int
	taskQueue   chan *websocket.Conn
	connHandler ConnHandler

	runing bool

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
		maxTasks:    maxTaskNum,
		taskQueue:   make(chan *websocket.Conn, maxTaskNum),
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

func (p *WorkerPool) AddTask(conn *websocket.Conn) {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()

	p.taskQueue <- conn
}

func (p *WorkerPool) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.runing {
		return
	}

	for i := 0; i < p.workers; i++ {
		go p.startWorker()
	}

	p.runing = true
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
							log.Warnf("connection handler error:%v", err)
						}
					}()

					p.connHandler(conn)

				}()
			}
		}
	}
}
