package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cclehui/server_on_gnet/websocket"
	"github.com/gobwas/ws/wsutil"
	"github.com/panjf2000/gnet"
)

//简单的echo server
func myDataHandler(param *websocket.DataHandlerParam) {

	log.Printf("server 接收到数据, opcode:%x, %s\n", param.OpCode, string(param.Request))

	response := fmt.Sprintf("response is :%s, 当前时间:%s\n", string(param.Request), time.Now().Format("2006-01-02 15:04:05"))

	//param.Writer.Write([]byte(response))

	//ws.WriteFrame(param.Writer, ws.NewTextFrame([]byte(response)))

	wsutil.WriteServerMessage(param.Writer, param.OpCode, []byte(response))

	return
}

//启动server
func ServerStart() {

	go startTestClient()

	//cclehui_test
	port := 8081
	tcpServer := websocket.NewServer(port)
	tcpServer.Handler = myDataHandler

	log.Fatal(gnet.Serve(tcpServer, fmt.Sprintf("tcp://:%d", port), gnet.WithMulticore(true)))

}

func startTestClient() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("websocket client start error :%v\n", err)
		}

	}()

	wsHome := func(w http.ResponseWriter, r *http.Request) {
		//websocket.ClientTemplate.Execute(w, "ws://"+r.Host+"/echo")
		websocket.ClientTemplate.Execute(w, "ws://192.168.67.129:8081")
	}

	addr := "0.0.0.0:8080"

	log.Printf("http server for websocket client is listen at :%s\n", addr)

	http.HandleFunc("/", wsHome)
	log.Fatal(http.ListenAndServe(addr, nil))

}
