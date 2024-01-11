package misc

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

// GetCallerWithLevel get shortened filename and line number of a caller function with a specific shift level.
// returns "undefined" in case of runtime.Caller function returns not ok code.
// Important: it calls runtime.Caller(level+1) beneath.
//
// Example:
//
//	1 func main() {
//	2     fmt.Println(GetCallerWithLevel(0))
//	3 }
//
// Output: tmp/sandbox/prog.go:2
func GetCallerWithLevel(level int) string {
	const (
		undefined = "undefined"
	)

	_, fileName, lineNum, ok := runtime.Caller(level + 1)
	if !ok {
		return undefined
	}

	return fmt.Sprintf("%s:%d", shortenFileName(fileName), lineNum)
}

// shortenFileName shortens file name.
func shortenFileName(filename string) string {
	wd, _ := os.Getwd()
	fullPath := path.Clean(filename)
	pathParts := strings.Split(wd, string(os.PathSeparator))
	last := pathParts[len(pathParts)-1]
	index := len(wd) - len(last)

	return fullPath[index:]
}

// Coalesce returns first not-nil element from elems variadic variable.
func Coalesce[T comparable](elems ...T) T {
	var zero T

	for _, elem := range elems {
		if elem != zero {
			return elem
		}
	}

	return [1]T{}[0]
}
