package main

import (
	client2 "github.com/codingXiang/go-harbor-client/client"
	"github.com/codingXiang/go-harbor-client/module/user"
	"github.com/codingXiang/configer"
	"github.com/codingXiang/go-logger"
)

func main() {
	//初始化 configer，設定預設讀取環境變數
	config := configer.NewConfigerCore("yaml", "harbor", "./config")
	config.SetAutomaticEnv("")

	logger.Log = logger.NewLogger(logger.Logger{Format: "text", Level: "debug"})

	client := client2.NewClient(config)

	userSvc := user.NewUserService(client)
	if users, _, err := userSvc.Current(); err == nil {
		//for _, user := range users {
		logger.Log.Debug(users)
		//}
	} else {
		logger.Log.Fatal(err)
	}
}
