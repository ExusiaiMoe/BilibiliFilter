package main

import (
	"github.com/AkameMoe/BilibiliFilter/database"
	"github.com/AkameMoe/BilibiliFilter/proxy"
	"github.com/AkameMoe/BilibiliFilter/scanner"
	"github.com/AkameMoe/BilibiliFilter/utils"
	"time"
)

func main()  {
	utils.StartLoggerModule()
	database.StartDatabaseModule()
	go proxy.StartProxyModule()
	time.Sleep(time.Second*10)
	scanner.StartScannerModule()
}