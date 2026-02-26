package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/MucheXD/WSCVenueBookingBackend/cmd/batchRegister"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/logger"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/server"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/repository"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/venuePermission"
)

func main() {
	go consoleScanner()
	logger.InitLogger()
	database.InitDatabase()
	database.InitRedis()
	venuePermission.RefreshVenueAccessCache(repository.NewVenueAccessLoader())
	server.InitServer()
}

func consoleScanner() {
	scanner := bufio.NewScanner(os.Stdin)
	fileInfo, _ := os.Stdin.Stat()
	if (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		fmt.Println("未检测到交互式终端，控制台监听已禁用")
		return
	} else {
		fmt.Println("控制台指令输入已启用，输入 help 获取可用指令列表")
	}

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		cmd := strings.SplitN(input, " ", 2)[0]

		switch cmd {
		case "help":
			fmt.Println("可用指令: stop, help" +
				"\nbatch-register [csv-file]: 从文件批量注册用户" +
				"\nbatch-register-example: 输出批量用户注册输入文件示例" +
				"\nstop: 结束服务器运行" +
				"\nhelp: 显示此帮助信息")
		case "batch-register":
			inputParts := strings.SplitN(input, " ", 2)
			if len(inputParts) != 2 {
				fmt.Println("错误: batch-register 指令需要一个参数，格式为: batch-register [csv-file]")
				continue
			}
			err := batchRegister.BatchRegisterFromFile(inputParts[1])
			if err != nil {
				fmt.Printf("批量注册失败: %v\n", err)
			} else {
				fmt.Println("批量注册成功")
			}
		case "batch-register-example":
			batchRegister.ShowExample()
		case "stop":
			fmt.Println("执行程序退出")
			os.Exit(0)
		case "":
			continue
		default:
			fmt.Printf("未知指令: %s\n", input)
		}
	}
}
