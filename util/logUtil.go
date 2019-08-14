package util

import (
	"autoCourse/config"
	"log"
	"os"
)

var lerr, linfo, laccept *log.Logger

func init() {
	p, err := os.Getwd()
	ErrExit(err)

	p += config.LogPath
	pe := p + config.LogErrFileName
	pi := p + config.LogInfoFileName
	pa := p + config.LogAcceptFileName

	err = os.MkdirAll(p, config.FilePermission)
	ErrExit(err)

	fe := openFile(&pe)
	fi := openFile(&pi)
	fa := openFile(&pa)

	lerr = log.New(fe, "", log.LstdFlags|log.Lshortfile)
	linfo = log.New(fi, "", log.LstdFlags)
	laccept = log.New(fa, "", log.LstdFlags)
}

// 打开文件，不存在，则创建
func openFile(path *string) (f *os.File) {
	_, err := os.Stat(*path)
	if err != nil {
		f, err = os.Create(*path)
		ErrExit(err)
	} else {
		f, err = os.OpenFile(*path, os.O_RDWR|os.O_APPEND, config.FilePermission)
		ErrExit(err)
	}
	return
}

// 打印错误，存在错误返回true
func LogErr(err error) (isErr bool) {
	if err != nil {
		isErr = true
		lerr.Println(err)
	}
	return
}

// 打印错误日志
func LogErrStr(str string) {
	lerr.Println(str)
}

// 打印错误，存在错误返回true
func LogErrAndStr(str string, err error) (isErr bool) {
	if err != nil {
		isErr = true
		lerr.Println(str+"  : ", err)
	}
	return
}

// 打印信息
func LogInfoF(format string, v ...interface{}) {
	linfo.Printf(format, v...)
}

// 打印连接信息
func LogAcceptF(format string, v ...interface{}) {
	laccept.Printf(format, v...)
}

// 打印信息
func LogInfo(infos ...interface{}) {
	linfo.Println(infos...)
}

// 打印连接信息
func LogAccept(accepts ...interface{}) {
	laccept.Println(accepts...)
}
