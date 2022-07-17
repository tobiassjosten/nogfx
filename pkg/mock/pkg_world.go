// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"github.com/tobiassjosten/nogfx/pkg"
	"sync"
)

// Ensure, that WorldMock does implement pkg.World.
// If this is not the case, regenerate this file with moq.
var _ pkg.World = &WorldMock{}

// WorldMock is a mock implementation of pkg.World.
//
// 	func TestSomethingThatUsesWorld(t *testing.T) {
//
// 		// make and configure a mocked pkg.World
// 		mockedWorld := &WorldMock{
// 			PrintFunc: func(bytes []byte)  {
// 				panic("mock out the Print method")
// 			},
// 			ProcessCommandFunc: func(bytes []byte)  {
// 				panic("mock out the ProcessCommand method")
// 			},
// 			ProcessInputFunc: func(input pkg.Input) pkg.Input {
// 				panic("mock out the ProcessInput method")
// 			},
// 			ProcessOutputFunc: func(output pkg.Output) pkg.Output {
// 				panic("mock out the ProcessOutput method")
// 			},
// 			SendFunc: func(bytes []byte)  {
// 				panic("mock out the Send method")
// 			},
// 		}
//
// 		// use mockedWorld in code that requires pkg.World
// 		// and then make assertions.
//
// 	}
type WorldMock struct {
	// PrintFunc mocks the Print method.
	PrintFunc func(bytes []byte)

	// ProcessCommandFunc mocks the ProcessCommand method.
	ProcessCommandFunc func(bytes []byte)

	// ProcessInputFunc mocks the ProcessInput method.
	ProcessInputFunc func(input pkg.Input) pkg.Input

	// ProcessOutputFunc mocks the ProcessOutput method.
	ProcessOutputFunc func(output pkg.Output) pkg.Output

	// SendFunc mocks the Send method.
	SendFunc func(bytes []byte)

	// calls tracks calls to the methods.
	calls struct {
		// Print holds details about calls to the Print method.
		Print []struct {
			// Bytes is the bytes argument value.
			Bytes []byte
		}
		// ProcessCommand holds details about calls to the ProcessCommand method.
		ProcessCommand []struct {
			// Bytes is the bytes argument value.
			Bytes []byte
		}
		// ProcessInput holds details about calls to the ProcessInput method.
		ProcessInput []struct {
			// Input is the input argument value.
			Input pkg.Input
		}
		// ProcessOutput holds details about calls to the ProcessOutput method.
		ProcessOutput []struct {
			// Output is the output argument value.
			Output pkg.Output
		}
		// Send holds details about calls to the Send method.
		Send []struct {
			// Bytes is the bytes argument value.
			Bytes []byte
		}
	}
	lockPrint          sync.RWMutex
	lockProcessCommand sync.RWMutex
	lockProcessInput   sync.RWMutex
	lockProcessOutput  sync.RWMutex
	lockSend           sync.RWMutex
}

// Print calls PrintFunc.
func (mock *WorldMock) Print(bytes []byte) {
	if mock.PrintFunc == nil {
		panic("WorldMock.PrintFunc: method is nil but World.Print was just called")
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
//     len(mockedWorld.PrintCalls())
func (mock *WorldMock) PrintCalls() []struct {
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

// ProcessCommand calls ProcessCommandFunc.
func (mock *WorldMock) ProcessCommand(bytes []byte) {
	if mock.ProcessCommandFunc == nil {
		panic("WorldMock.ProcessCommandFunc: method is nil but World.ProcessCommand was just called")
	}
	callInfo := struct {
		Bytes []byte
	}{
		Bytes: bytes,
	}
	mock.lockProcessCommand.Lock()
	mock.calls.ProcessCommand = append(mock.calls.ProcessCommand, callInfo)
	mock.lockProcessCommand.Unlock()
	mock.ProcessCommandFunc(bytes)
}

// ProcessCommandCalls gets all the calls that were made to ProcessCommand.
// Check the length with:
//     len(mockedWorld.ProcessCommandCalls())
func (mock *WorldMock) ProcessCommandCalls() []struct {
	Bytes []byte
} {
	var calls []struct {
		Bytes []byte
	}
	mock.lockProcessCommand.RLock()
	calls = mock.calls.ProcessCommand
	mock.lockProcessCommand.RUnlock()
	return calls
}

// ProcessInput calls ProcessInputFunc.
func (mock *WorldMock) ProcessInput(input pkg.Input) pkg.Input {
	if mock.ProcessInputFunc == nil {
		panic("WorldMock.ProcessInputFunc: method is nil but World.ProcessInput was just called")
	}
	callInfo := struct {
		Input pkg.Input
	}{
		Input: input,
	}
	mock.lockProcessInput.Lock()
	mock.calls.ProcessInput = append(mock.calls.ProcessInput, callInfo)
	mock.lockProcessInput.Unlock()
	return mock.ProcessInputFunc(input)
}

// ProcessInputCalls gets all the calls that were made to ProcessInput.
// Check the length with:
//     len(mockedWorld.ProcessInputCalls())
func (mock *WorldMock) ProcessInputCalls() []struct {
	Input pkg.Input
} {
	var calls []struct {
		Input pkg.Input
	}
	mock.lockProcessInput.RLock()
	calls = mock.calls.ProcessInput
	mock.lockProcessInput.RUnlock()
	return calls
}

// ProcessOutput calls ProcessOutputFunc.
func (mock *WorldMock) ProcessOutput(output pkg.Output) pkg.Output {
	if mock.ProcessOutputFunc == nil {
		panic("WorldMock.ProcessOutputFunc: method is nil but World.ProcessOutput was just called")
	}
	callInfo := struct {
		Output pkg.Output
	}{
		Output: output,
	}
	mock.lockProcessOutput.Lock()
	mock.calls.ProcessOutput = append(mock.calls.ProcessOutput, callInfo)
	mock.lockProcessOutput.Unlock()
	return mock.ProcessOutputFunc(output)
}

// ProcessOutputCalls gets all the calls that were made to ProcessOutput.
// Check the length with:
//     len(mockedWorld.ProcessOutputCalls())
func (mock *WorldMock) ProcessOutputCalls() []struct {
	Output pkg.Output
} {
	var calls []struct {
		Output pkg.Output
	}
	mock.lockProcessOutput.RLock()
	calls = mock.calls.ProcessOutput
	mock.lockProcessOutput.RUnlock()
	return calls
}

// Send calls SendFunc.
func (mock *WorldMock) Send(bytes []byte) {
	if mock.SendFunc == nil {
		panic("WorldMock.SendFunc: method is nil but World.Send was just called")
	}
	callInfo := struct {
		Bytes []byte
	}{
		Bytes: bytes,
	}
	mock.lockSend.Lock()
	mock.calls.Send = append(mock.calls.Send, callInfo)
	mock.lockSend.Unlock()
	mock.SendFunc(bytes)
}

// SendCalls gets all the calls that were made to Send.
// Check the length with:
//     len(mockedWorld.SendCalls())
func (mock *WorldMock) SendCalls() []struct {
	Bytes []byte
} {
	var calls []struct {
		Bytes []byte
	}
	mock.lockSend.RLock()
	calls = mock.calls.Send
	mock.lockSend.RUnlock()
	return calls
}
