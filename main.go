package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"syscall"

	"github.com/getlantern/systray"
	"gopkg.in/yaml.v3"
)

var cmd *exec.Cmd

type config struct {
	command struct {
	}
}

func main() {
	fmt.Println("Aria2Helper  By: LingYi0110")
	file, _ := ioutil.ReadFile("config.yaml")
	var data [7]Users
	yaml.Unmarshal(file, &data)
	start()
	systray.Run(onReady, onExit)

}

func start() {
	go func() {
		cmd = exec.Command("./aria2c", "--conf-path=aria2.conf")
		out, err := cmd.CombinedOutput()
		// err := cmd.Run()
		if err != nil {
			fmt.Printf("Aria2进程已结束: %s\n", err)
			start()
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
	// clean up here
}
