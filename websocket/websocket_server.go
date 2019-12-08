package websocket

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/cclehui/connection_platform/internal"
	"github.com/cclehui/connection_platform/log"
	"github.com/gorilla/websocket"
)

type WSServer struct {
	workerPool *WorkerPool
	eventPoll  *internal.Poll

	address string

	conCount int64

	connections map[int]*websocket.Conn

	mu sync.Mutex
}

func NewServer(address string, workerPool *WorkerPool) *WSServer {

	if workerPool == nil {
		workerPool = NewWorkerPool(DEFAULT_WORKER_NUM, DEFAULT_MAX_TASK_NUM, testHandleConn)
	}

	return &WSServer{
		workerPool:  workerPool,
		eventPoll:   internal.OpenPoll(),
		address:     address,
		connections: make(map[int]*websocket.Conn),
		mu:          sync.Mutex{},
	}

}

func (wss *WSServer) newConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warnf("connection upgrade fail: %v", err)
		return
	}

	wss.mu.Lock()
	defer wss.mu.Unlock()
	wss.connections[socketFD(conn.UnderlyingConn())] = conn

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

	//启动epoll 事件监听
	go wss.startEventPoll()

	//启动 worker 协程池
	wss.workerPool.Start()

	//set http handler
	http.HandleFunc("/ws", wss.newConnection)

	if err := http.ListenAndServe(wss.address, nil); err != nil {
		log.Fatal(fmt.Sprintf("%s", err))
	}
}

func (wss *WSServer) startEventPoll() {

	if wss.eventPoll == nil {
		panic("event poll not inited")
	}

	for {

	}

}

func epollEventHandler(fd int, note interface{}) {

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

func socketFD(conn net.Conn) int {
	//tls := reflect.TypeOf(conn.UnderlyingConn()) == reflect.TypeOf(&tls.Conn{})
	// Extract the file descriptor associated with the connection
	//connVal := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn").Elem()
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	//if tls {
	//  tcpConn = reflect.Indirect(tcpConn.Elem())
	//}
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")

	return int(pfdVal.FieldByName("Sysfd").Int())
}
