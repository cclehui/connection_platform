package conn_router

import (
	"errors"
	"sync"

	"github.com/cclehui/server_on_gnet/websocket"
)

//var localConnectionMap map[string]*websocket.GnetUpgraderConn
//并发考虑可以用多个 map 来降低锁的几率， 提高性能 cclehui_todo
var localConnectionMap sync.Map = sync.Map{}

func AddLocalConnection(yewuId string, wsConn *websocket.GnetUpgraderConn) bool {
	if yewuId == "" || wsConn == nil {
		return false
	}

	localConnectionMap.Store(yewuId, wsConn)
	return true
}

func LoadLocalConnection(yewuId string) (*websocket.GnetUpgraderConn, error) {
	if value, ok := localConnectionMap.Load(yewuId); ok {
		if result, ok2 := value.(*websocket.GnetUpgraderConn); ok2 {
			return result, nil
		}
	}

	return nil, errors.New("connection not exist")
}

func RemoveLocalConnection(yewuId string) {
	localConnectionMap.Delete(yewuId)
}
