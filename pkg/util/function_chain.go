package util

type FunctionChain struct {
	chain []func()
}

func (fc *FunctionChain) Add(f func()) *FunctionChain {
	fc.chain = append(fc.chain, f)
	return fc
}

func (fc *FunctionChain) InvokeForward() *FunctionChain {
	for _, f := range fc.chain {
		f()
	}
	return fc
}

func (fc *FunctionChain) InvokeReverse() *FunctionChain {
	for i := len(fc.chain) - 1; i >= 0; i-- {
		fc.chain[i]()
	}
	return fc
}
