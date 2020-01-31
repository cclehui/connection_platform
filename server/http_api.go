package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/cclehui/connection_platform/server/conn_router"
	"github.com/cclehui/connection_platform/server/protobuf_def"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
)

const (
	HTTP_SUCCESS = 0
	HTTP_FAIL    = 9999
)

//http server api
func startHttpApiServer() {
	defer func() {
		err := recover()
		log.Printf("exception HttpApiServer stoped , error:%v\n", err)
	}()

	router := gin.New()
	router.POST("http/router/send/msg", httpRouterSendMsg)

	router.Run(httpApiServerAddr)

}

//消息下行 api
func httpRouterSendMsg(c *gin.Context) {

	//yewuId string, msg string
	yewuId := c.PostForm("yewuId")
	msg := c.PostForm("msg")

	if yewuId == "" {
		c.JSON(http.StatusOK, apiFailRes("yewuId empty"))
		return
	}

	if msg == "" {
		c.JSON(http.StatusOK, apiFailRes("msg is empty"))
		return
	}

	if wsConn, err := conn_router.LoadLocalConnection(yewuId); err == nil {
		//连接在当前服务器上
		getServer().SendDownStreamMsg(wsConn, ws.OpText, []byte(msg))

		c.JSON(http.StatusOK, apiSuccessRes(fmt.Sprintf("%s_%s", yewuId, msg)))

	} else {
		//获取远程服务信息
		if grpcServerAddr, err2 := conn_router.GetServerAddr(yewuId); err2 == nil {
			if grpcServerAddr != "" {
				grpcConn := getGrpcClientConn(grpcServerAddr)
				if grpcConn != nil {
					client := protobuf_def.NewServerApiServiceClient(grpcConn)

					param := &protobuf_def.ParamSendDownStreamMsg{YewuId: yewuId, Msg: msg}

					response, err3 := client.SendDownStreamMsg(context.Background(), param)

					if err3 != nil || response.Status != GRPC_SUCCESS {
						logStr := fmt.Sprintf("调用rpc方法失败, SendDownStreamMsg, param:%v, response:%v, error:%v", param, response, err3)
						log.Printf("%s\n", logStr)
						c.JSON(http.StatusOK, apiFailRes(logStr))
					} else {

						//grpc 调用成功
						c.JSON(http.StatusOK, apiSuccessRes(response.Data))
					}

				} else {
					c.JSON(http.StatusOK, apiFailRes(fmt.Sprintf("get grpc connection fail, %s", grpcServerAddr)))
				}
			} else {
				c.JSON(http.StatusOK, apiFailRes("find connection grpcServerAddr empty"))
			}

		} else {
			c.JSON(http.StatusOK, apiFailRes(fmt.Sprintf("find connection error:%v, %v", err, err2)))
		}

	}

	return
}

func apiSuccessRes(data interface{}) gin.H {
	return gin.H{
		"status": HTTP_SUCCESS,
		"data":   data,
		"msg":    "",
	}
}

func apiFailRes(msg string) gin.H {
	return gin.H{
		"status": HTTP_FAIL,
		"data":   "",
		"msg":    msg,
	}
}
