package util

import log "github.com/sirupsen/logrus"

// FunctionChain represents a chain of functions that can be called in sequence.
// Functions added to the chain will have any panics they generate caught and logged.
type FunctionChain struct {
	chain []func()
}

// Add a function to the FunctionChain.
func (fc *FunctionChain) Add(f func()) *FunctionChain {
	fc.chain = append(fc.chain, func() {
		defer func() {
			if err := recover(); err != nil {
				log.Warnln("panic occurred in function chain:", err)
			}
		}()
		f()
	})
	return fc
}

// InvokeForward invokes the functions in the chain in the order they were added.
func (fc *FunctionChain) InvokeForward() *FunctionChain {
	for _, f := range fc.chain {
		f()
	}
	return fc
}

// InvokeReverse invokes the functions in the chain in the reverse order they were added.
func (fc *FunctionChain) InvokeReverse() *FunctionChain {
	for i := len(fc.chain) - 1; i >= 0; i-- {
		fc.chain[i]()
	}
	return fc
}
