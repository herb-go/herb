package querybuilder

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"
)

//Debug querybuild debug mode.
//If enabled,all sql commands and args will bt sent to logger
var Debug = false

type timeDuation struct {
	Label    string
	Duration time.Duration
}

var timeDurationList = []timeDuation{

	timeDuation{
		Label:    "minute",
		Duration: time.Minute,
	},
	timeDuation{
		Label:    "seconds",
		Duration: time.Second,
	},
	timeDuation{
		Label:    "milliseconds",
		Duration: time.Millisecond,
	},
	timeDuation{
		Label:    "microseconds",
		Duration: time.Microsecond,
	},
	timeDuation{
		Label:    "nanoseconds",
		Duration: time.Nanosecond,
	},
}

//Logger query logger
var Logger func(timestamp int64, cmd string, args []interface{})

//DefaultLogger default logger which print qurey  command.args and spent time to std.output
var DefaultLogger = func(timestamp int64, cmd string, args []interface{}) {
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
	for _, v := range timeDurationList {
		if spent > 10*v.Duration {
			fmt.Printf("\t%d %s \n", spent/v.Duration, v.Label)
			break
		}
	}
	fmt.Println()
}

func init() {
	Logger = DefaultLogger
}
