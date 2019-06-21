package output

import (
	"log"
	"os"
)

var testLogger *Logger

func ExampleNew() {
	var (
		f   *os.File
		err error
	)

	if f, err = os.Create("output.log"); err != nil {
		log.Fatalf("Error creating file: %v", err)
	}

	testLogger = New(f, nil)
}

func ExampleLogger_Print() {
	testLogger.Print("Hello world!")
}

func ExampleLogger_Printf() {
	testLogger.Printf("Hello world! My name is %s.", "Loggy")
}

func ExampleLogger_Notification() {
	testLogger.Notification("Hello world! My name is %s.", "Loggy")
}

func ExampleLogger_Success() {
	testLogger.Success("Hello world! My name is %s.", "Loggy")
}

func ExampleLogger_Warning() {
	testLogger.Warning("Hello world! My name is %s.", "Loggy")
}

func ExampleLogger_Error() {
	testLogger.Error("Hello world! My name is %s.", "Loggy")
}

func ExampleLogger_Debug() {
	testLogger.Debug("Hello world! My name is %s.", "Loggy")
}
