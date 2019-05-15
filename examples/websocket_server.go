package main

import (
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/cclehui/connection_platform/log"
	"github.com/cclehui/connection_platform/websocket"
	gwebsocket "github.com/gorilla/websocket"
)

func handleConn(conn *gwebsocket.Conn) {

	messageType, reader, err := conn.NextReader()

	if err != nil {
		log.Warnf("failed to handleConn NextReader  %v", err)
		conn.Close()
		return
	}

	writer, _ := conn.NextWriter(messageType)

	io.Copy(writer, reader)
	//io.Copy(ioutil.Discard, reader)
}

func main() {

	log.SetDefaultLogger(log.NewLogger(log.LEVEL_DEBUG, os.Stdout, log.DefaultLogFlag))

	// Enable pprof hooks
	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("Pprof failed: %v", err)
		}
	}()

	workerPool := websocket.NewWorkerPool(2, 1000, handleConn)

	server := websocket.NewServer(":8972", workerPool)
	server.Start()

}
