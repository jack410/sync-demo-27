package main

import (
	"example.com/m/config"
	"example.com/m/server"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
)

func main() {
	chChromeDie := make(chan struct{})
	chBackendDie := make(chan struct{})
	chSignal := listenToInterrupt()
	//启动gin服务
	go server.Run()
	go startBrower(chChromeDie, chBackendDie)
	//如果没有值则一直等待（阻塞），直到有信号输入
	//select可以监听多个管道，只要有一个管道有信号则进行下一步
	//如果没有信号，select就等待（阻塞）
	for {
		select {
		case <-chSignal:
			chBackendDie <- struct{}{}
		case <-chChromeDie:
			os.Exit(0)
		}
	}
}

func startBrower(chChromeDie, chBackendDie chan struct{}) {
	chromePath := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	var tmpDir string
	name, _ := ioutil.TempDir("", "lorca")
	tmpDir = name
	//删除缓存文件
	defer os.RemoveAll(tmpDir)
	fmt.Println(tmpDir)
	cmd := exec.Command(chromePath, "--app=http://localhost:"+config.GetPort()+"/static/index.html",
		fmt.Sprintf("--user-data-dir=%s", tmpDir), "--no-first-run")
	cmd.Start()
	//同时等待两件事，要么后端进程退出，要么浏览器退出
	//所以要用两个匿名函数
	go func() {
		<-chBackendDie
		cmd.Process.Kill()
	}()
	//上面<-chBackendDie触发后会触发cmd.Wait()
	go func() {
		cmd.Wait()
		chChromeDie <- struct{}{}
	}()
}

//监听中断信号
func listenToInterrupt() chan os.Signal {
	chSignal := make(chan os.Signal, 1)
	//signal.Notify订阅os.Interrupt信号,一旦有信号则往chSingal管道里写一个信号
	signal.Notify(chSignal, os.Interrupt)
	return chSignal
}
