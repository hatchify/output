package output

func ExamplePrint() {
	Print("Hello world!")
}

func ExamplePrintf() {
	Printf("Hello world! My name is %s.", "Loggy")
}

func ExampleSuccess() {
	Success("Hello world! My name is %s.", "Loggy")
}

func ExampleWarning() {
	Warning("Hello world! My name is %s.", "Loggy")
}

func ExampleError() {
	Error("Hello world! My name is %s.", "Loggy")
}

func ExampleDebug() {
	Debug("Hello world! My name is %s.", "Loggy")
}
