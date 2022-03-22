package controller

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

func AddressesController(c *gin.Context) {
	//获取当前电脑的所有ip地址
	addrs, _ := net.InterfaceAddrs()
	var result []string
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		// 断言address里的地址是ip地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				result = append(result, ipnet.IP.String())
			}
		}
	}
	//转为json写入address接口的http respond
	c.JSON(http.StatusOK, gin.H{"addresses": result})
}
