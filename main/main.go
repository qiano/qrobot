package main

import (
	"fmt"
	// "github.com/qshuai162/btcrobot/src/config"
	// "io/ioutil"
	"math/rand"
	// "os"
	// "path/filepath"
	"runtime"
	// "strconv"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 设置随机数种子
	rand.Seed(time.Now().Unix())
}

func main() {
	printBanner()

	//go webui.RunServer()

	RunRobot()
}

func printBanner() {
	version := "V0.1"
	fmt.Println("[ ---------------------------------------------------------->>> ")
	fmt.Println(" BTC/LTC自动化算法交易引擎", version)
	fmt.Println(" *@警告：API key和密码存放在conf/secret.json文件内，共享给他人前请务必删除，注意账号安全！！")
	fmt.Println(" <<<----------------------------------------------------------] ")
}
