// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"github.com/tobiassjosten/nogfx/pkg"
	"sync"
)

// Ensure, that ModuleMock does implement pkg.Module.
// If this is not the case, regenerate this file with moq.
var _ pkg.Module = &ModuleMock{}

// ModuleMock is a mock implementation of pkg.Module.
//
// 	func TestSomethingThatUsesModule(t *testing.T) {
//
// 		// make and configure a mocked pkg.Module
// 		mockedModule := &ModuleMock{
// 			ProcessInputFunc: func(bytes []byte) [][]byte {
// 				panic("mock out the ProcessInput method")
// 			},
// 			ProcessOutputFunc: func(bytes []byte) [][]byte {
// 				panic("mock out the ProcessOutput method")
// 			},
// 		}
//
// 		// use mockedModule in code that requires pkg.Module
// 		// and then make assertions.
//
// 	}
type ModuleMock struct {
	// ProcessInputFunc mocks the ProcessInput method.
	ProcessInputFunc func(bytes []byte) [][]byte

	// ProcessOutputFunc mocks the ProcessOutput method.
	ProcessOutputFunc func(bytes []byte) [][]byte

	// calls tracks calls to the methods.
	calls struct {
		// ProcessInput holds details about calls to the ProcessInput method.
		ProcessInput []struct {
			// Bytes is the bytes argument value.
			Bytes []byte
		}
		// ProcessOutput holds details about calls to the ProcessOutput method.
		ProcessOutput []struct {
			// Bytes is the bytes argument value.
			Bytes []byte
		}
	}
	lockProcessInput  sync.RWMutex
	lockProcessOutput sync.RWMutex
}

// ProcessInput calls ProcessInputFunc.
func (mock *ModuleMock) ProcessInput(bytes []byte) [][]byte {
	if mock.ProcessInputFunc == nil {
		panic("ModuleMock.ProcessInputFunc: method is nil but Module.ProcessInput was just called")
	}
	callInfo := struct {
		Bytes []byte
	}{
		Bytes: bytes,
	}
	mock.lockProcessInput.Lock()
	mock.calls.ProcessInput = append(mock.calls.ProcessInput, callInfo)
	mock.lockProcessInput.Unlock()
	return mock.ProcessInputFunc(bytes)
}

// ProcessInputCalls gets all the calls that were made to ProcessInput.
// Check the length with:
//     len(mockedModule.ProcessInputCalls())
func (mock *ModuleMock) ProcessInputCalls() []struct {
	Bytes []byte
} {
	var calls []struct {
		Bytes []byte
	}
	mock.lockProcessInput.RLock()
	calls = mock.calls.ProcessInput
	mock.lockProcessInput.RUnlock()
	return calls
}

// ProcessOutput calls ProcessOutputFunc.
func (mock *ModuleMock) ProcessOutput(bytes []byte) [][]byte {
	if mock.ProcessOutputFunc == nil {
		panic("ModuleMock.ProcessOutputFunc: method is nil but Module.ProcessOutput was just called")
	}
	callInfo := struct {
		Bytes []byte
	}{
		Bytes: bytes,
	}
	mock.lockProcessOutput.Lock()
	mock.calls.ProcessOutput = append(mock.calls.ProcessOutput, callInfo)
	mock.lockProcessOutput.Unlock()
	return mock.ProcessOutputFunc(bytes)
}

// ProcessOutputCalls gets all the calls that were made to ProcessOutput.
// Check the length with:
//     len(mockedModule.ProcessOutputCalls())
func (mock *ModuleMock) ProcessOutputCalls() []struct {
	Bytes []byte
} {
	var calls []struct {
		Bytes []byte
	}
	mock.lockProcessOutput.RLock()
	calls = mock.calls.ProcessOutput
	mock.lockProcessOutput.RUnlock()
	return calls
}
