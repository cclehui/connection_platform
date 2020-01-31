package server

import (
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

var grpcClientConn map[string]*grpc.ClientConn = make(map[string]*grpc.ClientConn)
var grpcClientsOnce map[string]*sync.Once = make(map[string]*sync.Once)
var clientsMu sync.Mutex

func getGrpcClientConn(serverAddr string) *grpc.ClientConn {
	clientsMu.Lock()
	pOnce, ok := grpcClientsOnce[serverAddr]
	if !ok {
		pOnce = &sync.Once{}
		grpcClientsOnce[serverAddr] = pOnce
	}
	clientsMu.Unlock()

	(*pOnce).Do(func() {
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
		if err != nil {
			log.Printf("grpc 连接服务端失败, %s, %v\n", serverAddr, err)
			return
		}
		grpcClientConn[serverAddr] = conn
	})

	return grpcClientConn[serverAddr]
}
