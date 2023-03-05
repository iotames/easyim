package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/iotames/miniutils"
)

const RUN_LOCK_FILE = "run.lock"

// 停止程序运行
func stopApp(port int) error {
	var err error
	// pid, err := os.ReadFile(RUN_LOCK_FILE)
	pid := miniutils.GetPidByPort(sconf.Port)
	if pid > -1 {
		err = miniutils.KillPid(fmt.Sprintf("%d", pid))
		if err != nil {
			return fmt.Errorf("kill app fail(%v)", err)
		}
		return os.Remove(RUN_LOCK_FILE)
	}
	return fmt.Errorf("app is not running!")
}

// 以守护进程的方式后台运行
func startDaemon() error {
	args := os.Args
	var newArgs []string
	for i, arg := range args {
		if arg == "-d" || arg == "--d" || i == 0 {
			continue
		}
		newArgs = append(newArgs, arg)
	}
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		winArgs := append([]string{"/c", args[0]}, newArgs...)
		cmd = exec.Command("cmd", winArgs...)
	} else {
		cmd = exec.Command(args[0], newArgs...)
	}
	cmd.Env = os.Environ()
	err := cmd.Start()
	if err != nil {
		return err
	}
	return os.WriteFile(RUN_LOCK_FILE, []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0755)
}
