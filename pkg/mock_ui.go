// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package pkg

import (
	"context"
	"sync"
)

// Ensure, that UIMock does implement UI.
// If this is not the case, regenerate this file with moq.
var _ UI = &UIMock{}

// UIMock is a mock implementation of UI.
//
// 	func TestSomethingThatUsesUI(t *testing.T) {
//
// 		// make and configure a mocked UI
// 		mockedUI := &UIMock{
// 			InputsFunc: func() <-chan []byte {
// 				panic("mock out the Inputs method")
// 			},
// 			MaskInputFunc: func()  {
// 				panic("mock out the MaskInput method")
// 			},
// 			OutputsFunc: func() chan<- []byte {
// 				panic("mock out the Outputs method")
// 			},
// 			PrintFunc: func(bytes []byte)  {
// 				panic("mock out the Print method")
// 			},
// 			RunFunc: func(contextMoqParam context.Context) error {
// 				panic("mock out the Run method")
// 			},
// 			UnmaskInputFunc: func()  {
// 				panic("mock out the UnmaskInput method")
// 			},
// 		}
//
// 		// use mockedUI in code that requires UI
// 		// and then make assertions.
//
// 	}
type UIMock struct {
	// InputsFunc mocks the Inputs method.
	InputsFunc func() <-chan []byte

	// MaskInputFunc mocks the MaskInput method.
	MaskInputFunc func()

	// OutputsFunc mocks the Outputs method.
	OutputsFunc func() chan<- []byte

	// PrintFunc mocks the Print method.
	PrintFunc func(bytes []byte)

	// RunFunc mocks the Run method.
	RunFunc func(contextMoqParam context.Context) error

	// UnmaskInputFunc mocks the UnmaskInput method.
	UnmaskInputFunc func()

	// calls tracks calls to the methods.
	calls struct {
		// Inputs holds details about calls to the Inputs method.
		Inputs []struct {
		}
		// MaskInput holds details about calls to the MaskInput method.
		MaskInput []struct {
		}
		// Outputs holds details about calls to the Outputs method.
		Outputs []struct {
		}
		// Print holds details about calls to the Print method.
		Print []struct {
			// Bytes is the bytes argument value.
			Bytes []byte
		}
		// Run holds details about calls to the Run method.
		Run []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
		}
		// UnmaskInput holds details about calls to the UnmaskInput method.
		UnmaskInput []struct {
		}
	}
	lockInputs      sync.RWMutex
	lockMaskInput   sync.RWMutex
	lockOutputs     sync.RWMutex
	lockPrint       sync.RWMutex
	lockRun         sync.RWMutex
	lockUnmaskInput sync.RWMutex
}

// Inputs calls InputsFunc.
func (mock *UIMock) Inputs() <-chan []byte {
	if mock.InputsFunc == nil {
		panic("UIMock.InputsFunc: method is nil but UI.Inputs was just called")
	}
	callInfo := struct {
	}{}
	mock.lockInputs.Lock()
	mock.calls.Inputs = append(mock.calls.Inputs, callInfo)
	mock.lockInputs.Unlock()
	return mock.InputsFunc()
}

// InputsCalls gets all the calls that were made to Inputs.
// Check the length with:
//     len(mockedUI.InputsCalls())
func (mock *UIMock) InputsCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockInputs.RLock()
	calls = mock.calls.Inputs
	mock.lockInputs.RUnlock()
	return calls
}

// MaskInput calls MaskInputFunc.
func (mock *UIMock) MaskInput() {
	if mock.MaskInputFunc == nil {
		panic("UIMock.MaskInputFunc: method is nil but UI.MaskInput was just called")
	}
	callInfo := struct {
	}{}
	mock.lockMaskInput.Lock()
	mock.calls.MaskInput = append(mock.calls.MaskInput, callInfo)
	mock.lockMaskInput.Unlock()
	mock.MaskInputFunc()
}

// MaskInputCalls gets all the calls that were made to MaskInput.
// Check the length with:
//     len(mockedUI.MaskInputCalls())
func (mock *UIMock) MaskInputCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockMaskInput.RLock()
	calls = mock.calls.MaskInput
	mock.lockMaskInput.RUnlock()
	return calls
}

// Outputs calls OutputsFunc.
func (mock *UIMock) Outputs() chan<- []byte {
	if mock.OutputsFunc == nil {
		panic("UIMock.OutputsFunc: method is nil but UI.Outputs was just called")
	}
	callInfo := struct {
	}{}
	mock.lockOutputs.Lock()
	mock.calls.Outputs = append(mock.calls.Outputs, callInfo)
	mock.lockOutputs.Unlock()
	return mock.OutputsFunc()
}

// OutputsCalls gets all the calls that were made to Outputs.
// Check the length with:
//     len(mockedUI.OutputsCalls())
func (mock *UIMock) OutputsCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockOutputs.RLock()
	calls = mock.calls.Outputs
	mock.lockOutputs.RUnlock()
	return calls
}

// Print calls PrintFunc.
func (mock *UIMock) Print(bytes []byte) {
	if mock.PrintFunc == nil {
		panic("UIMock.PrintFunc: method is nil but UI.Print was just called")
	}
	callInfo := struct {
		Bytes []byte
	}{
		Bytes: bytes,
	}
	mock.lockPrint.Lock()
	mock.calls.Print = append(mock.calls.Print, callInfo)
	mock.lockPrint.Unlock()
	mock.PrintFunc(bytes)
}

// PrintCalls gets all the calls that were made to Print.
// Check the length with:
//     len(mockedUI.PrintCalls())
func (mock *UIMock) PrintCalls() []struct {
	Bytes []byte
} {
	var calls []struct {
		Bytes []byte
	}
	mock.lockPrint.RLock()
	calls = mock.calls.Print
	mock.lockPrint.RUnlock()
	return calls
}

// Run calls RunFunc.
func (mock *UIMock) Run(contextMoqParam context.Context) error {
	if mock.RunFunc == nil {
		panic("UIMock.RunFunc: method is nil but UI.Run was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
	}{
		ContextMoqParam: contextMoqParam,
	}
	mock.lockRun.Lock()
	mock.calls.Run = append(mock.calls.Run, callInfo)
	mock.lockRun.Unlock()
	return mock.RunFunc(contextMoqParam)
}

// RunCalls gets all the calls that were made to Run.
// Check the length with:
//     len(mockedUI.RunCalls())
func (mock *UIMock) RunCalls() []struct {
	ContextMoqParam context.Context
} {
	var calls []struct {
		ContextMoqParam context.Context
	}
	mock.lockRun.RLock()
	calls = mock.calls.Run
	mock.lockRun.RUnlock()
	return calls
}

// UnmaskInput calls UnmaskInputFunc.
func (mock *UIMock) UnmaskInput() {
	if mock.UnmaskInputFunc == nil {
		panic("UIMock.UnmaskInputFunc: method is nil but UI.UnmaskInput was just called")
	}
	callInfo := struct {
	}{}
	mock.lockUnmaskInput.Lock()
	mock.calls.UnmaskInput = append(mock.calls.UnmaskInput, callInfo)
	mock.lockUnmaskInput.Unlock()
	mock.UnmaskInputFunc()
}

// UnmaskInputCalls gets all the calls that were made to UnmaskInput.
// Check the length with:
//     len(mockedUI.UnmaskInputCalls())
func (mock *UIMock) UnmaskInputCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockUnmaskInput.RLock()
	calls = mock.calls.UnmaskInput
	mock.lockUnmaskInput.RUnlock()
	return calls
}
