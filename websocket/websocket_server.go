package websocket

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"syscall"

	"github.com/cclehui/connection_platform/log"
	"github.com/gorilla/websocket"
)

type WSServer struct {
	workerPool *WorkerPool

	address string

	conCount int64
}

func NewServer(address string, workerPool *WorkerPool) *WSServer {

	if workerPool == nil {
		workerPool = NewWorkerPool(DEFAULT_WORKER_NUM, DEFAULT_MAX_TASK_NUM, testHandleConn)
	}

	return &WSServer{workerPool: workerPool, address: address}

}

func (wss *WSServer) ws(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warnf("connection upgrade fail: %v", err)
		return
	}

	n := atomic.AddInt64(&wss.conCount, 1)
	log.Debugf("new connection conn num: %v", n)

	//把 新连接的业务处理加到 workerPool中去
	wss.workerPool.AddTask(conn)

}

func testHandleConn(conn *websocket.Conn) {
	_, reader, err := conn.NextReader()

	if err != nil {
		log.Warnf("failed to testHandleConn NextReader  %v", err)
		conn.Close()
		return
	}

	io.Copy(ioutil.Discard, reader)
}

func (wss *WSServer) Start() {
	// Increase resources limitations
	setSysResLimit()

	//启动 worker 协程池
	wss.workerPool.Start()

	//set http handler
	http.HandleFunc("/ws", wss.ws)

	if err := http.ListenAndServe(wss.address, nil); err != nil {
		log.Fatal(fmt.Sprintf("%s", err))
	}
}

// Increase resources limitations
func setSysResLimit() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
}
