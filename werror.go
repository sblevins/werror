package werror

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// getCallerInfo returns the package, function, and line of the first external caller
func getCallerInfo() (pkg string, function string, line int) {
	// skip=0 is getCallerInfo
	// skip=1 is Wrap or New
	// skip=2 is the actual caller
	pc := make([]uintptr, 10)
	n := runtime.Callers(3, pc) // skip 3 to get the first external caller
	if n == 0 {
		return "unknown", "unknown", 0
	}
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	funcName := frame.Function // e.g., "main.doSomething"
	split := strings.Split(funcName, ".")
	if len(split) >= 2 {
		pkg = strings.Join(split[:len(split)-1], ".")
		function = split[len(split)-1]
	} else {
		pkg = "unknown"
		function = funcName
	}

	return pkg, function, frame.Line
}

func Wrap(msg string, err error) error {
	pkg, function, line := getCallerInfo()
	infoMsg := fmt.Sprintf("error encountered in pkg: %s, caller: %s:%d, msg: %s", pkg, function, line, msg)
	return fmt.Errorf("%s\n > %w", infoMsg, err)
}

func New(msg string) error {
	pkg, function, line := getCallerInfo()
	infoMsg := fmt.Sprintf("error raised in pkg: %s, caller: %s:%d, msg: %s", pkg, function, line, msg)
	return errors.New(infoMsg)
}

