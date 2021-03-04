package util

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type FunctionStackTestSuite struct {
	suite.Suite
}

func (suite *FunctionStackTestSuite) SetupTest() {

}

func TestFunctionStack(t *testing.T) {
	suite.Run(t, new(FunctionStackTestSuite))
}

func (suite *FunctionStackTestSuite) Test_AppendToNil() {
	var fs *FunctionStack

	fs = fs.Append(func() {})
	suite.NotNil(fs)
	suite.Equal(1, fs.Depth())
}

func (suite *FunctionStackTestSuite) Test_AppendTwo() {
	var fs *FunctionStack

	fs = fs.
		Append(func() {}).
		Append(func() {})
	suite.Equal(2, fs.Depth())
}

func (suite *FunctionStackTestSuite) Test_InvokeNil() {
	var fs *FunctionStack

	fs = fs.
		Invoke()
	suite.Nil(fs)
}

func (suite *FunctionStackTestSuite) Test_InvokeIfNil() {
	var fs *FunctionStack

	fs = fs.
		InvokeIf(true)
	suite.Nil(fs)
}

func (suite *FunctionStackTestSuite) Test_InvokeOne() {
	var fs *FunctionStack

	var counter = 0
	fs = fs.
		Append(func() { counter++ }).
		Invoke()
	suite.Equal(1, counter)
	suite.Equal(0, fs.Depth())
}

func (suite *FunctionStackTestSuite) Test_InvokeTwoWithOnePanic() {
	var fs *FunctionStack

	var counter = 0
	fs = fs.
		Append(func() { panic("The sky is falling!") }).
		Append(func() { counter++ }).
		Invoke()
	suite.Equal(1, counter)
	suite.Equal(0, fs.Depth())
}

func (suite *FunctionStackTestSuite) Test_InvokeOneIfFalse() {
	var fs *FunctionStack

	var counter = 0
	fs = fs.
		Append(func() { counter++ }).
		InvokeIf(false)
	suite.Equal(0, counter)
	suite.Equal(1, fs.Depth())
}

func (suite *FunctionStackTestSuite) Test_InvokeMultipleOrder() {
	var fs *FunctionStack

	var counter = 0
	_ = fs.
		Append(func() { counter++ }).
		Append(func() { counter = 1 }).
		Invoke()

	suite.Equal(2, counter)
}
