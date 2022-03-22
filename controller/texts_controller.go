package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func TextsController(c *gin.Context) {
	var json struct {
		Raw string `json:"raw"`
	}
	//将获取到的数据传给json变量
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		//获取当前目录
		exe, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		dir := filepath.Dir(exe)
		if err != nil {
			log.Fatal(err)
		}
		//随机生成一个字符串并复制给filename，用来做上传后的文件名
		filename := uuid.New().String()
		//拼接uploads的绝对路径
		uploads := filepath.Join(dir, "uploads")
		//创建uploads目录
		err = os.MkdirAll(uploads, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		//合成上传后文件的路径
		fullpath := path.Join("uploads", filename+".txt")
		//把json.Raw写到文件里
		err = ioutil.WriteFile(filepath.Join(dir, fullpath), []byte(json.Raw), 0644)
		if err != nil {
			log.Fatal(err)
		}
		//返回文件的路径到texts接口的http respond，比如/uploads/c07b266c-53ce-435d-91d0-bb4cbbb00ecb.txt
		c.JSON(http.StatusOK, gin.H{"url": "/" + fullpath})
	}
}
