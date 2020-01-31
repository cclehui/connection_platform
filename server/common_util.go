package server

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
)

//获取本机ip 地址
func getLocalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(fmt.Sprintf("getLocalIp error, %v", err))
	}

	result := "0.0.0.0"

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				result = ipnet.IP.String()
				break
				//fmt.Println(result)
			}
		}
	}

	return result
}

func getPortFromAddr(addr string) int {

	r, _ := regexp.Compile(".*?:(\\d+)")

	matches := r.FindStringSubmatch(addr)
	//fmt.Println(matches)

	result := 0

	if len(matches) == 2 {
		result, _ = strconv.Atoi(matches[1])
	}

	return result

}
