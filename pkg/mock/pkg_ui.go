// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/navigation"
	"sync"
)

// Ensure, that UIMock does implement pkg.UI.
// If this is not the case, regenerate this file with moq.
var _ pkg.UI = &UIMock{}

// UIMock is a mock implementation of pkg.UI.
//
// 	func TestSomethingThatUsesUI(t *testing.T) {
//
// 		// make and configure a mocked pkg.UI
// 		mockedUI := &UIMock{
// 			AddVitalFunc: func(s string, ifaceVal interface{})  {
// 				panic("mock out the AddVital method")
// 			},
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
// 			SetRoomFunc: func(room *navigation.Room)  {
// 				panic("mock out the SetRoom method")
// 			},
// 			UnmaskInputFunc: func()  {
// 				panic("mock out the UnmaskInput method")
// 			},
// 			UpdateVitalFunc: func(s string, n1 int, n2 int)  {
// 				panic("mock out the UpdateVital method")
// 			},
// 		}
//
// 		// use mockedUI in code that requires pkg.UI
// 		// and then make assertions.
//
// 	}
type UIMock struct {
	// AddVitalFunc mocks the AddVital method.
	AddVitalFunc func(s string, ifaceVal interface{})

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

	// SetRoomFunc mocks the SetRoom method.
	SetRoomFunc func(room *navigation.Room)

	// UnmaskInputFunc mocks the UnmaskInput method.
	UnmaskInputFunc func()

	// UpdateVitalFunc mocks the UpdateVital method.
	UpdateVitalFunc func(s string, n1 int, n2 int)

	// calls tracks calls to the methods.
	calls struct {
		// AddVital holds details about calls to the AddVital method.
		AddVital []struct {
			// S is the s argument value.
			S string
			// IfaceVal is the ifaceVal argument value.
			IfaceVal interface{}
		}
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
		// SetRoom holds details about calls to the SetRoom method.
		SetRoom []struct {
			// Room is the room argument value.
			Room *navigation.Room
		}
		// UnmaskInput holds details about calls to the UnmaskInput method.
		UnmaskInput []struct {
		}
		// UpdateVital holds details about calls to the UpdateVital method.
		UpdateVital []struct {
			// S is the s argument value.
			S string
			// N1 is the n1 argument value.
			N1 int
			// N2 is the n2 argument value.
			N2 int
		}
	}
	lockAddVital    sync.RWMutex
	lockInputs      sync.RWMutex
	lockMaskInput   sync.RWMutex
	lockOutputs     sync.RWMutex
	lockPrint       sync.RWMutex
	lockRun         sync.RWMutex
	lockSetRoom     sync.RWMutex
	lockUnmaskInput sync.RWMutex
	lockUpdateVital sync.RWMutex
}

// AddVital calls AddVitalFunc.
func (mock *UIMock) AddVital(s string, ifaceVal interface{}) {
	if mock.AddVitalFunc == nil {
		panic("UIMock.AddVitalFunc: method is nil but UI.AddVital was just called")
	}
	callInfo := struct {
		S        string
		IfaceVal interface{}
	}{
		S:        s,
		IfaceVal: ifaceVal,
	}
	mock.lockAddVital.Lock()
	mock.calls.AddVital = append(mock.calls.AddVital, callInfo)
	mock.lockAddVital.Unlock()
	mock.AddVitalFunc(s, ifaceVal)
}

// AddVitalCalls gets all the calls that were made to AddVital.
// Check the length with:
//     len(mockedUI.AddVitalCalls())
func (mock *UIMock) AddVitalCalls() []struct {
	S        string
	IfaceVal interface{}
} {
	var calls []struct {
		S        string
		IfaceVal interface{}
	}
	mock.lockAddVital.RLock()
	calls = mock.calls.AddVital
	mock.lockAddVital.RUnlock()
	return calls
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

// SetRoom calls SetRoomFunc.
func (mock *UIMock) SetRoom(room *navigation.Room) {
	if mock.SetRoomFunc == nil {
		panic("UIMock.SetRoomFunc: method is nil but UI.SetRoom was just called")
	}
	callInfo := struct {
		Room *navigation.Room
	}{
		Room: room,
	}
	mock.lockSetRoom.Lock()
	mock.calls.SetRoom = append(mock.calls.SetRoom, callInfo)
	mock.lockSetRoom.Unlock()
	mock.SetRoomFunc(room)
}

// SetRoomCalls gets all the calls that were made to SetRoom.
// Check the length with:
//     len(mockedUI.SetRoomCalls())
func (mock *UIMock) SetRoomCalls() []struct {
	Room *navigation.Room
} {
	var calls []struct {
		Room *navigation.Room
	}
	mock.lockSetRoom.RLock()
	calls = mock.calls.SetRoom
	mock.lockSetRoom.RUnlock()
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

// UpdateVital calls UpdateVitalFunc.
func (mock *UIMock) UpdateVital(s string, n1 int, n2 int) {
	if mock.UpdateVitalFunc == nil {
		panic("UIMock.UpdateVitalFunc: method is nil but UI.UpdateVital was just called")
	}
	callInfo := struct {
		S  string
		N1 int
		N2 int
	}{
		S:  s,
		N1: n1,
		N2: n2,
	}
	mock.lockUpdateVital.Lock()
	mock.calls.UpdateVital = append(mock.calls.UpdateVital, callInfo)
	mock.lockUpdateVital.Unlock()
	mock.UpdateVitalFunc(s, n1, n2)
}

// UpdateVitalCalls gets all the calls that were made to UpdateVital.
// Check the length with:
//     len(mockedUI.UpdateVitalCalls())
func (mock *UIMock) UpdateVitalCalls() []struct {
	S  string
	N1 int
	N2 int
} {
	var calls []struct {
		S  string
		N1 int
		N2 int
	}
	mock.lockUpdateVital.RLock()
	calls = mock.calls.UpdateVital
	mock.lockUpdateVital.RUnlock()
	return calls
}
