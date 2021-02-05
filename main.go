package main

import (
	"runtime"
	"syscall"

	"github.com/filebrowser/filebrowser/v2/cmd"
)

func main() {

	syscall.Umask(0) // FIX: Create file or dir mode always is 0755

	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()
}
