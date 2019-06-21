package output

import (
	"os"

	"github.com/fatih/color"
)

const dot = "● "

var (
	green  = color.New(color.FgGreen)
	yellow = color.New(color.FgYellow)
	red    = color.New(color.FgRed)

	greenDot  = green.Sprint(dot)
	yellowDot = yellow.Sprint(dot)
	redDot    = red.Sprint(dot)
)

// Default output
var logger = New(NopCloser(os.Stdout), nil)

// Print will print to the underlying writer
func Print(str string) {
	logger.Print(str)
}

// Printf will print a formatted message to the underlying writer
func Printf(format string, values ...interface{}) {
	logger.Printf(format, values...)
}

// Success will output a success message
func Success(format string, values ...interface{}) {
	logger.Success(format, values...)
}

// Warning will output a warning message
func Warning(format string, values ...interface{}) {
	logger.Warning(format, values...)
}

// Error will output an error message
func Error(format string, values ...interface{}) {
	logger.Error(format, values...)
}

// Debug will log a debug line
func Debug(format string, values ...interface{}) {
	logger.Debug(format, values...)
}
