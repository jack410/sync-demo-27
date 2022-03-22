package server

import (
	"embed"
	"example.com/m/config"
	c "example.com/m/controller"
	"example.com/m/server/ws"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

//下面这句话的意思是打包go的时候把后面这个目录打包进去
//go:embed frontend/dist/*
var FS embed.FS

func Run() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	//把打包好的静态文件变成一个结构化的目录
	staticFiles, _ := fs.Sub(FS, "frontend/dist")
	router.POST("/api/v1/files", c.FilesController)
	router.GET("/api/v1/qrcodes", c.QrcodesController)
	router.GET("/uploads/:path", c.UploadsController)
	router.POST("/api/v1/texts", c.TextsController)
	router.GET("/api/v1/addresses", c.AddressesController)
	router.StaticFS("/static", http.FS(staticFiles))
	hub := ws.NewHub()
	go hub.Run()
	router.GET("/ws", func(c *gin.Context) {
		ws.HttpController(c, hub)
	})
	//NoRoute表示用户访问路径没匹配到程序定义的路由
	router.NoRoute(func(c *gin.Context) {
		//获取用户访问的路径
		path := c.Request.URL.Path
		//判断路径是否以static开头
		if strings.HasPrefix(path, "/static/") {
			reader, err := staticFiles.Open("index.html")
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()
			stat, err := reader.Stat()
			if err != nil {
				log.Fatal(err)
			}
			c.DataFromReader(http.StatusOK, stat.Size(), "text/html", reader, nil)
			//如果不是static开头则返回404
		} else {
			c.Status(http.StatusNotFound)
		}
	})
	router.Run(":" + config.GetPort())
}
