package main

import (
	"flag"
	"os"
	"runtime"
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/xtls/xray-core/main/commands/base"
	_ "github.com/xtls/xray-core/main/distro/all"
)

// setProcessName sets the process name for the current process.
// This is a Linux-specific implementation using prctl.
func setProcessName(name string) {
	if runtime.GOOS != "linux" {
		return // Or handle for other OSes if needed
	}
	// The process name is limited to 16 bytes in Linux, including the null terminator.
	bytes := make([]byte, 16)
	copy(bytes, name)
	err := unix.Prctl(unix.PR_SET_NAME, uintptr(unsafe.Pointer(&bytes[0])), 0, 0, 0)
	if err != nil {
		// Silently ignore the error, as it's not critical for the application to run.
	}
}

func main() {
	setProcessName("caddy-debug") // You can change the process name here

	os.Args = getArgsV4Compatible()

	base.RootCommand.Long = "Xray is a platform for building proxies."
	base.RootCommand.Commands = append(
		[]*base.Command{
			cmdRun,
			cmdVersion,
		},
		base.RootCommand.Commands...,
	)
	base.Execute()
}

func getArgsV4Compatible() []string {
	if len(os.Args) == 1 {
		return []string{os.Args[0], "run"}
	}
	if os.Args[1][0] != '-' {
		return os.Args
	}
	version := false
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.BoolVar(&version, "version", false, "")
	// parse silently, no usage, no error output
	fs.Usage = func() {}
	fs.SetOutput(&null{})
	err := fs.Parse(os.Args[1:])
	if err == flag.ErrHelp {
		// fmt.Println("DEPRECATED: -h, WILL BE REMOVED IN V5.")
		// fmt.Println("PLEASE USE: xray help")
		// fmt.Println()
		return []string{os.Args[0], "help"}
	}
	if version {
		// fmt.Println("DEPRECATED: -version, WILL BE REMOVED IN V5.")
		// fmt.Println("PLEASE USE: xray version")
		// fmt.Println()
		return []string{os.Args[0], "version"}
	}
	// fmt.Println("COMPATIBLE MODE, DEPRECATED.")
	// fmt.Println("PLEASE USE: xray run [arguments] INSTEAD.")
	// fmt.Println()
	return append([]string{os.Args[0], "run"}, os.Args[1:]...)
}

type null struct{}

func (n *null) Write(p []byte) (int, error) {
	return len(p), nil
}
