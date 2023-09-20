package misc

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

func GetCallerWithLevel(level int) string {
	const (
		undefined = "undefined"
	)

	_, fileName, lineNum, ok := runtime.Caller(level)
	if !ok {
		return undefined
	}

	return fmt.Sprintf("%s:%d", shortenFileName(fileName), lineNum)
}

func shortenFileName(filename string) string {
	wd, _ := os.Getwd()
	fullPath := path.Clean(filename)
	pathParts := strings.Split(wd, string(os.PathSeparator))
	last := pathParts[len(pathParts)-1]
	index := len(wd) - len(last)

	return fullPath[index:]
}
