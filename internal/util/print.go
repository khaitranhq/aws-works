package util

import "fmt"

func ErrorPrint(txt string) {
	fmt.Printf("\033[1;31m%s\033[0m\n", txt)
}
