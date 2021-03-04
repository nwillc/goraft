package util

import (
	"log"
)

// FunctionStack represents a stack of functions that can be called in sequence.
// Functions added to the stack will have any panics they generate caught and logged.
type FunctionStack struct {
	fn    func()
	count int
}

// Append a function to the FunctionStack.
func (fs *FunctionStack) Append(f func()) *FunctionStack {
	if fs == nil {
		//goland:noinspection GoAssignmentToReceiver
		fs = &FunctionStack{}
	}
	var priorFn = fs.fn
	fs.fn = func() {
		defer func() {
			if p := recover(); p != nil {
				log.Println("FunctionStack suppressing a panic:", p)
			}
		}()
		f()
		if priorFn != nil {
			priorFn()
		}
	}
	fs.count++
	return fs
}

// Invoke invokes the functions in the FunctionStack. The functions are invoked in the reverse order they were added and
// any panic in a function is suppressed.
func (fs *FunctionStack) Invoke() *FunctionStack {
	if fs == nil || fs.count == 0 {
		return fs
	}
	fs.fn()
	fs.fn = nil
	fs.count = 0
	return fs
}

// InvokeIf invokes the FunctionStack ifTrue is true.
func (fs *FunctionStack) InvokeIf(ifTrue bool) *FunctionStack {
	if !ifTrue {
		return fs
	}
	return fs.Invoke()
}

// Depth of the stack.
func (fs *FunctionStack) Depth() int {
	if fs == nil {
		return 0
	}

	return fs.count
}
