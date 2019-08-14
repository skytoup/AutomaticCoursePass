package util

import (
	"fmt"
	"os"
)

// 存在err则打印，且返回true
func ErrPrintln(err error) (isErr bool) {
	if err != nil {
		fmt.Println(err)
		isErr = true
	}
	return
}

// 存在err，则打印后panic
func ErrPanic(err error) {
	if ErrPrintln(err) {
		panic(err)
	}
}

// 存在err，则打印后退出程序
func ErrExit(err error) {
	if ErrPrintln(err) {
		os.Exit(-1)
	}
}
