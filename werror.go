package werror

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
)

const myName = "werror.go"

func getCallerInfo() (pkg string, function string, line string) {
	stack := string(debug.Stack())
	_, stack, _ = strings.Cut(stack, "[running]:\n") //cut the calling goroutine
	_, stack, _ = strings.Cut(stack, myName)         //cut the call made by this function
	_, stack, _ = strings.Cut(stack, myName)         //cut the call made by Wrap/New
	_, stack, _ = strings.Cut(stack, "\n")           //cut the line number and hex data at the end
	caller, stack, _ := strings.Cut(stack, "\n")     //cut the function call from the stack data
	_, stack, _ = strings.Cut(stack, ":")            //cut the stack from the line number and hex data
	line, _, _ = strings.Cut(stack, " ")             //cut the line number from the hex data
	line, _, _ = strings.Cut(line, "\n")             //if there was no extra hex data we need to cut here to recover
	pkg, function, _ = strings.Cut(caller, ".")      //split the calling package from the calling function

	return pkg, function, line
}

//wraps errors to include information about the caller and a specified message
func Wrap(msg string, err error) error {
	pkg, function, line := getCallerInfo()
	infoMsg := "error encountered in pkg: " + pkg + ", caller: " + function + ":" + line + ", msg: " + msg
	return fmt.Errorf(infoMsg+"\n >%w", err)
}

//creates a new error with information about the caller and a specified message
func New(msg string) error {
	pkg, function, line := getCallerInfo()
	infoMsg := "error raised in pkg: " + pkg + ", caller: " + function + ":" + line + ", msg: " + msg
	return errors.New(infoMsg)
}
