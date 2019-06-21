package output

// NewWrapper will return a new wrapper with the default logger
func NewWrapper(prefix string) *Wrapper {
	var w Wrapper
	w.o = logger
	w.prefix = prefix
	return &w
}

// NewWrapperWithOutputter will return a new wrapper with a custom outputter
func NewWrapperWithOutputter(o Outputter, prefix string) *Wrapper {
	var w Wrapper
	w.o = o
	w.prefix = prefix
	return &w
}

// Wrapper will wrap a logger
type Wrapper struct {
	o Outputter

	// Prefix string to be inserted to all format's
	prefix string
}

// Print will print to the underlying writer
func (w *Wrapper) Print(str string) {
	str = w.prefix + str
	w.o.Print(str)
}

// Printf will print a formatted message to the underlying writer
func (w *Wrapper) Printf(format string, values ...interface{}) {
	format = w.prefix + format
	w.o.Printf(format, values...)
}

// Notification will output a notification message
func (w *Wrapper) Notification(format string, values ...interface{}) {
	format = w.prefix + format
	w.o.Notification(format, values...)
}

// Success will output a success message
func (w *Wrapper) Success(format string, values ...interface{}) {
	format = w.prefix + format
	w.o.Success(format, values...)
}

// Warning will output a warning message
func (w *Wrapper) Warning(format string, values ...interface{}) {
	format = w.prefix + format
	w.o.Warning(format, values...)
}

// Error will output an error message
func (w *Wrapper) Error(format string, values ...interface{}) {
	format = w.prefix + format
	w.o.Error(format, values...)
}

// Debug will log a debug line
func (w *Wrapper) Debug(format string, values ...interface{}) {
	format = w.prefix + format
	w.o.Debug(format, values...)
}
