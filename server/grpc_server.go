package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/cclehui/connection_platform/server/conn_router"
	"github.com/cclehui/connection_platform/server/protobuf_def"
	"github.com/gobwas/ws"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	GRPC_SUCCESS int32 = 0
	GRPC_FAIL    int32 = 9999
)

type grpcApiServer struct{}

//发送下行消息
func (gas *grpcApiServer) SendDownStreamMsg(ctx context.Context, param *protobuf_def.ParamSendDownStreamMsg) (*protobuf_def.ResSendDownStreamMsg, error) {

	yewuId := param.GetYewuId()
	msg := param.GetMsg()

	result := &protobuf_def.ResSendDownStreamMsg{}
	result.Status = GRPC_SUCCESS

	if wsConn, err := conn_router.LoadLocalConnection(yewuId); err == nil {
		//连接在当前服务器上
		err2 := getServer().SendDownStreamMsg(wsConn, ws.OpText, []byte(msg))
		if err2 == nil {
			return result, nil
		} else {
			result.Status = GRPC_FAIL
			result.Msg = fmt.Sprintf("grpc server SendDownStreamMsg error:%v", err2)
		}
	} else {
		result.Status = GRPC_FAIL
		result.Msg = fmt.Sprintf("grpc not find connection error:%v", err)
	}

	return result, nil
}

func startGrpcServer(serverAddr string) {
	defer func() {
		err := recover()
		log.Printf("exception GrpcApiServer stoped , error:%v\n", err)
	}()

	lis, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v, addr:%v", err, serverAddr)
	}

	var kaep = keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}

	var kasp = keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
		MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
		MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
	}

	server := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))

	protobuf_def.RegisterServerApiServiceServer(server, &grpcApiServer{})

	log.Printf("grpc 服务启动, tcp addr:%s\n", serverAddr)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
