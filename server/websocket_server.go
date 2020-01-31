package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cclehui/connection_platform/server/conn_router"
	"github.com/cclehui/server_on_gnet/websocket"
	"github.com/gobwas/ws/wsutil"
	"github.com/panjf2000/gnet"
)

//数据处理过程
func myDataHandler(param *websocket.DataHandlerParam) {

	log.Printf("server 接收到数据, opcode:%x, %s\n", param.OpCode, string(param.Request))

	//业务id
	yewuId := "cclehui_001" //cclehui_test

	response := fmt.Sprintf("response is :%s, 当前时间:%s\n", string(param.Request), time.Now().Format("2006-01-02 15:04:05"))

	if param.WSConn.UniqId == "" {
		//未做 connection 和 业务关系路由
		conn_router.AddLocalConnection(yewuId, param.WSConn) //存储在内存中

		//业务id和 server的对应关系
		curGrpcServerAddr := fmt.Sprintf("%s:%d", getIpFromAddr(param.Server.Addr), getPortFromAddr(grpcApiServerAddr))
		err := conn_router.AddServerRoute(yewuId, curGrpcServerAddr)

		if err == nil {
			param.WSConn.UniqId = yewuId //标识连接对应关系已处理完毕
		}
	}

	//param.Writer.Write([]byte(response))

	//ws.WriteFrame(param.Writer, ws.NewTextFrame([]byte(response)))

	wsutil.WriteServerMessage(param.Writer, param.OpCode, []byte(response))

	return
}

//连接关闭事件处理
func connCloseHandler(wsConn *websocket.GnetUpgraderConn) {
	if wsConn == nil {
		return
	}

	if wsConn.UniqId == "" {
		return
	}

	conn_router.RemoveServerRoute(wsConn.UniqId)
	conn_router.RemoveLocalConnection(wsConn.UniqId)
}

var serverAddr string = "192.168.67.129:8081"
var testClientAddr string = "0.0.0.0:8080"
var httpApiServerAddr string = "0.0.0.0:8082" //http sever api
var grpcApiServerAddr string = "0.0.0.0:8083" //grpc sever api
var tcpServer *websocket.WebSocketServer

//启动server
func ServerStart() {

	//cclehui_test
	localIp := getLocalIp() //获取本机ip的方法可能不行
	port := getPortFromAddr(serverAddr)

	serverAddr = fmt.Sprintf("%s:%d", localIp, port)

	//启动测试client
	go startTestClient()

	//启动消息下行服务http api
	go startHttpApiServer()

	//grpc服务启动
	go startGrpcServer(grpcApiServerAddr)

	//cclehui_test 端口写死

	server := getServer()
	log.Fatal(gnet.Serve(server, fmt.Sprintf("tcp://%s", serverAddr), gnet.WithMulticore(true)))

}

//
func getServer() *websocket.WebSocketServer {
	if tcpServer == nil {
		tcpServer = websocket.NewServer(serverAddr)
		tcpServer.Handler = myDataHandler
		tcpServer.ConnCloseHandler = connCloseHandler //连接关闭处理函数
	}

	return tcpServer
}

func startTestClient() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("websocket client start error :%v\n", err)
		}

	}()

	wsHome := func(w http.ResponseWriter, r *http.Request) {
		//websocket.ClientTemplate.Execute(w, "ws://"+r.Host+"/echo")
		websocket.ClientTemplate.Execute(w, fmt.Sprintf("ws://%s", serverAddr))
	}

	log.Printf("http server for websocket client is listen at :%s\n", testClientAddr)

	http.HandleFunc("/", wsHome)
	log.Fatal(http.ListenAndServe(testClientAddr, nil))

}
