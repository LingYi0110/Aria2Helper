package main

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"syscall"

	"github.com/getlantern/systray"
	"github.com/spf13/viper"
)

var cmd *exec.Cmd

func main() {
	fmt.Println("Aria2Helper  By: LingYi0110")

	// 读取配置
	viper.SetConfigFile("./config.yml")
	viper.ReadInConfig()
	command := viper.Sub("command")
	name := command.GetString("name")
	args := command.GetStringSlice("args")
	enableUpdate := viper.Sub("autoUpdateBt-tracker").GetBool("enabled")
	urls := viper.Sub("autoUpdateBt-tracker").GetStringSlice("urls")
	path := viper.Sub("autoUpdateBt-tracker").GetString("aria2ConfigPath")

	// 更新配置
	if enableUpdate {
		updateBtTracker(urls, path)
	}
	start(name, args)
	systray.Run(onReady, onExit)

}
func updateBtTracker(urls []string, configPath string) {
	var address string

	// 获取Bt-tracker的内容
	for _, vaule := range urls {
		resp, err := http.Get(vaule)
		if err != nil {
			fmt.Printf("Error: 从%s获取Bt-tracker失败", vaule)
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		address += string(body) + ","
	}

	// 写入aria2的config文件
	viper.SetConfigFile(configPath)
	viper.SetConfigType("properties") // 不知道为什么，aria2的config文件格式和java properties一致
	viper.ReadInConfig()
	viper.Set("bt-tracker", address)
	viper.WriteConfig()

}
func start(name string, args []string) {
	go func() {

		cmd = exec.Command(name)
		cmd.Args = args
		out, err := cmd.CombinedOutput()
		// err := cmd.Run()
		if err != nil {
			fmt.Printf("Aria2进程已结束: %s\n", err)
			start(name, args)
		}
		fmt.Print(out)
	}()
}
func hideWindows() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	user32 := syscall.NewLazyDLL("user32.dll")

	getConsoleWindow := kernel32.NewProc("GetConsoleWindow")
	showWindowAsync := user32.NewProc("ShowWindowAsync")
	isWindowVisible := user32.NewProc("IsWindowVisible")

	handle, r2, lastErr := getConsoleWindow.Call()
	if handle == 0 {
		fmt.Println("Error: ", handle, r2, lastErr)
	}

	status, _, _ := isWindowVisible.Call(handle)

	// uintptr to bool
	// 1          true
	// 0          false
	if status == 1 {
		showWindowAsync.Call(handle, 0)
	} else {
		showWindowAsync.Call(handle, 1)
	}

}
func onReady() {
	systray.SetIcon(Data)
	systray.SetTitle("Aria2Helper")
	systray.SetTooltip("Aria2Helper")

	mShow := systray.AddMenuItem("显示/隐藏窗口", "显示或者隐藏窗口")
	mRestart := systray.AddMenuItem("重启Aria2", "重新启动Aria2")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出", "退出")

	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				hideWindows()
			case <-mRestart.ClickedCh:
				cmd.Process.Kill()
				cmd.Wait()
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()

}

func onExit() {
}
