package util

import (
	"os"
	"reflect"
	"runtime"
)

func GetFunctionName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

func CurrentFunction() string {
	counter, _, _, success := runtime.Caller(1)

	if !success {
		println("functionName: runtime.Caller: failed")
		os.Exit(1)
	}

	return runtime.FuncForPC(counter).Name()
}
