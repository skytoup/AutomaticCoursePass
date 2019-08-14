package main

import (
	"autoCourse/config"
	_ "autoCourse/db"
	"autoCourse/util"
	"autoCourse/viewController"
	"net/http"
	"time"
)

func main() {
	handle()
	listen()
}

// 设置路由
func handle() {
	http.Handle("/", &viewController.MainVC{})
}

// 监听端口，设置服务器参数
func listen() {
	server := http.Server{Addr: config.Address + config.Port, Handler: nil}
	server.ListenAndServe()
	var err error
	server.ReadTimeout, err = time.ParseDuration(config.ReadTimeout)
	util.ErrExit(err)
	server.WriteTimeout, err = time.ParseDuration(config.WriteTimeout)
	util.ErrExit(err)
	server.MaxHeaderBytes = config.MaxHeaderBytes
	util.ErrExit(err)
}
