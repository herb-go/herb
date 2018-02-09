package query

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"
)

var Debug = false

var Logger func(timestamp int64, cmd string, args []interface{})

func init() {
	Logger = func(timestamp int64, cmd string, args []interface{}) {
		spent := time.Duration((time.Now().UnixNano() - timestamp)) * time.Nanosecond
		fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " SQL query debug:")
		fmt.Println("Query:")
		lines := strings.Split(cmd, "\n")
		for k := range lines {
			fmt.Println("\t" + lines[k])
		}
		fmt.Println("Args:")
		argsString := make([]string, len(args))
		for k := range args {
			argsString[k] = fmt.Sprint(args[k])
		}
		fmt.Println("\t[" + strings.Join(argsString, " , ") + "]")
		stacks := string(debug.Stack())
		stacklines := strings.Split(stacks, "\n")
		if len(stacklines) > 9 {
			fmt.Println("Stack:")
			fmt.Println("\t" + stacklines[7])
			fmt.Println("\t" + stacklines[8])
		}
		fmt.Println("Time spent:")
		if spent > time.Hour {
			fmt.Printf("\t%f %s \n", spent.Hours(), "hours")
		} else if spent > time.Minute {
			fmt.Printf("\t%f %s \n", spent.Minutes(), "minutes")
		} else if spent > time.Second {
			fmt.Printf("\t%f %s \n", spent.Seconds(), "seconds")
		} else if spent > time.Millisecond {
			fmt.Printf("\t%f %s \n", spent.Seconds()*1000, "milliseconds")
		} else if spent > time.Microsecond {
			fmt.Printf("\t%f %s \n", spent.Seconds()*1000*1000, "microseconds")
		} else {
			fmt.Printf("\t%d %s \n", spent.Nanoseconds(), "nanoseconds")
		}
		fmt.Println()
	}
}
