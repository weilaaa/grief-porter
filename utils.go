package main

import "fmt"

const (
	// pre ground color
	black  = 30
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	white  = 37

	// background color need add 10
)

func colorPrint(color int, format string, args ...interface{}) {
	content := fmt.Sprintf(format, args...)
	fmt.Printf("%c[%d;%d;%dm%s%c[0m\n", 0x1B, 0, 0, color, content, 0x1B)
}
