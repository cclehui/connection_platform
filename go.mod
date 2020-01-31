module github.com/cclehui/connection_platform

go 1.13

require (
	github.com/cclehui/server_on_gnet v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.5.0
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/gobwas/ws v1.0.2
	github.com/golang/protobuf v1.3.2
	github.com/panjf2000/gnet v1.0.0-beta.8
	github.com/smartystreets-prototypes/go-disruptor v0.0.0-20180723194425-e0f8f9247cc2 // indirect
	google.golang.org/grpc v1.25.1
)

replace github.com/cclehui/server_on_gnet => /home/cclehui/go_workspace/src/github.com/cclehui/server_on_gnet
