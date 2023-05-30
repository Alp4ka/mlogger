package misc

import (
	"fmt"
	"runtime"
)

func GetCaller() string {
	const undefined = "undefined"

	_, fileName, lineNum, ok := runtime.Caller(2)
	if !ok {
		return undefined
	}

	return fmt.Sprintf("%s:%d", fileName, lineNum)
}
