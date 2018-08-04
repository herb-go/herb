package util

import (
	"fmt"
	"log"
	"runtime/debug"
	"strings"
)

func Recover(args ...interface{}) {
	if r := recover(); r != nil {
		err := r.(error)
		if _, ok := LoggerIgnoredErrors[err]; ok == false {
			lines := strings.Split(string(debug.Stack()), "\n")
			length := len(lines)
			maxLength := LoggerMaxLength*2 + 7
			if length > maxLength {
				length = maxLength
			}
			var output = make([]string, length-6)
			output[0] = fmt.Sprintf("Panic: %s", err.Error())
			output[0] += "\n" + lines[0]
			copy(output[1:], lines[7:])
			log.Println(strings.Join(output, "\n"))

		}
	}

}
