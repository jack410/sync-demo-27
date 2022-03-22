package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func getUploadDir() (uploads string) {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(exe)
	uploads = filepath.Join(dir, "uploads")
	return
}

func UploadsController(c *gin.Context) {
	//获取到路由里的path，即:path这里的path的值
	if path := c.Param("path"); path != "" {
		target := filepath.Join(getUploadDir(), path)
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+path)
		c.Header("Content-Type", "application/octet-stream")
		//给前端发送target，即文件下载路径
		c.File(target)
	} else {
		c.Status(http.StatusNotFound)
	}
}
