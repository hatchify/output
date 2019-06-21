# Output [![GoDoc](https://godoc.org/github.com/hatchify/output?status.svg)](https://godoc.org/github.com/hatchify/output) ![Status](https://img.shields.io/badge/status-beta-yellow.svg)
Output is an output logger helper library a few simple features:

* Thread-safe
* Colored dot prefix for status (success, warning, error)
* Debug output (with filename and line number)

![screenshot](https://github.com/hatchify/output/blob/master/screenshot.png?raw=true "Screenshot of output")

## Usage
The primary usage of output is utilizing the package-level logger. See below for examples of the available methods:

### Print
```go 
func ExamplePrint() {
	Print("Hello world!")
}
```

### Printf
```go 
func ExamplePrintf() {
	Printf("Hello world! My name is %s.", "Loggy")
}
```

### Success
```go 
func ExampleSuccess() {
	Success("Hello world! My name is %s.", "Loggy")
}
```

### Warning
```go 
func ExampleWarning() {
	Warning("Hello world! My name is %s.", "Loggy")
}
```

### Error
```go 
func ExampleError() {
	Error("Hello world! My name is %s.", "Loggy")
}
```

### Debug
```go 
func ExampleDebug() {
	Debug("Hello world! My name is %s.", "Loggy")
}
```











