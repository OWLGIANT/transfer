package main

import (
	"github.com/sirupsen/logrus"
	"time"
	"transfer/utils"
)

func main() {
	utils.LogInit()
	logrus.Info("bpy_transfer:服务启动", time.Now().Format("2006-01-02 15:04:05"))
	ticker := time.Tick(time.Hour)
	for range ticker {
		utils.CompressAndTransferLoop()
	}
}
