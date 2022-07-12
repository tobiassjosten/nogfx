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
// 			SetCharacterFunc: func(character pkg.Character)  {
// 				panic("mock out the SetCharacter method")
// 			},
// 			SetRoomFunc: func(room *navigation.Room)  {
// 				panic("mock out the SetRoom method")
// 			},
// 			SetTargetFunc: func(target *pkg.Target)  {
// 				panic("mock out the SetTarget method")
// 			},
// 			UnmaskInputFunc: func()  {
// 				panic("mock out the UnmaskInput method")
// 			},
// 		}
//
// 		// use mockedUI in code that requires pkg.UI
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

	// SetCharacterFunc mocks the SetCharacter method.
	SetCharacterFunc func(character pkg.Character)

	// SetRoomFunc mocks the SetRoom method.
	SetRoomFunc func(room *navigation.Room)

	// SetTargetFunc mocks the SetTarget method.
	SetTargetFunc func(target *pkg.Target)

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
		// SetCharacter holds details about calls to the SetCharacter method.
		SetCharacter []struct {
			// Character is the character argument value.
			Character pkg.Character
		}
		// SetRoom holds details about calls to the SetRoom method.
		SetRoom []struct {
			// Room is the room argument value.
			Room *navigation.Room
		}
		// SetTarget holds details about calls to the SetTarget method.
		SetTarget []struct {
			// Target is the target argument value.
			Target *pkg.Target
		}
		// UnmaskInput holds details about calls to the UnmaskInput method.
		UnmaskInput []struct {
		}
	}
	lockInputs       sync.RWMutex
	lockMaskInput    sync.RWMutex
	lockOutputs      sync.RWMutex
	lockPrint        sync.RWMutex
	lockRun          sync.RWMutex
	lockSetCharacter sync.RWMutex
	lockSetRoom      sync.RWMutex
	lockSetTarget    sync.RWMutex
	lockUnmaskInput  sync.RWMutex
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

// SetCharacter calls SetCharacterFunc.
func (mock *UIMock) SetCharacter(character pkg.Character) {
	if mock.SetCharacterFunc == nil {
		panic("UIMock.SetCharacterFunc: method is nil but UI.SetCharacter was just called")
	}
	callInfo := struct {
		Character pkg.Character
	}{
		Character: character,
	}
	mock.lockSetCharacter.Lock()
	mock.calls.SetCharacter = append(mock.calls.SetCharacter, callInfo)
	mock.lockSetCharacter.Unlock()
	mock.SetCharacterFunc(character)
}

// SetCharacterCalls gets all the calls that were made to SetCharacter.
// Check the length with:
//     len(mockedUI.SetCharacterCalls())
func (mock *UIMock) SetCharacterCalls() []struct {
	Character pkg.Character
} {
	var calls []struct {
		Character pkg.Character
	}
	mock.lockSetCharacter.RLock()
	calls = mock.calls.SetCharacter
	mock.lockSetCharacter.RUnlock()
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

// SetTarget calls SetTargetFunc.
func (mock *UIMock) SetTarget(target *pkg.Target) {
	if mock.SetTargetFunc == nil {
		panic("UIMock.SetTargetFunc: method is nil but UI.SetTarget was just called")
	}
	callInfo := struct {
		Target *pkg.Target
	}{
		Target: target,
	}
	mock.lockSetTarget.Lock()
	mock.calls.SetTarget = append(mock.calls.SetTarget, callInfo)
	mock.lockSetTarget.Unlock()
	mock.SetTargetFunc(target)
}

// SetTargetCalls gets all the calls that were made to SetTarget.
// Check the length with:
//     len(mockedUI.SetTargetCalls())
func (mock *UIMock) SetTargetCalls() []struct {
	Target *pkg.Target
} {
	var calls []struct {
		Target *pkg.Target
	}
	mock.lockSetTarget.RLock()
	calls = mock.calls.SetTarget
	mock.lockSetTarget.RUnlock()
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
