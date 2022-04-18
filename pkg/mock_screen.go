// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package pkg

import (
	"github.com/gdamore/tcell/v2"
	"sync"
)

// Ensure, that ScreenMock does implement tcell.Screen.
// If this is not the case, regenerate this file with moq.
var _ tcell.Screen = &ScreenMock{}

// ScreenMock is a mock implementation of tcell.Screen.
//
// 	func TestSomethingThatUsesScreen(t *testing.T) {
//
// 		// make and configure a mocked tcell.Screen
// 		mockedScreen := &ScreenMock{
// 			BeepFunc: func() error {
// 				panic("mock out the Beep method")
// 			},
// 			CanDisplayFunc: func(r rune, checkFallbacks bool) bool {
// 				panic("mock out the CanDisplay method")
// 			},
// 			ChannelEventsFunc: func(ch chan<- tcell.Event, quit <-chan struct{})  {
// 				panic("mock out the ChannelEvents method")
// 			},
// 			CharacterSetFunc: func() string {
// 				panic("mock out the CharacterSet method")
// 			},
// 			ClearFunc: func()  {
// 				panic("mock out the Clear method")
// 			},
// 			ColorsFunc: func() int {
// 				panic("mock out the Colors method")
// 			},
// 			DisableMouseFunc: func()  {
// 				panic("mock out the DisableMouse method")
// 			},
// 			DisablePasteFunc: func()  {
// 				panic("mock out the DisablePaste method")
// 			},
// 			EnableMouseFunc: func(mouseFlagss ...tcell.MouseFlags)  {
// 				panic("mock out the EnableMouse method")
// 			},
// 			EnablePasteFunc: func()  {
// 				panic("mock out the EnablePaste method")
// 			},
// 			FillFunc: func(n rune, style tcell.Style)  {
// 				panic("mock out the Fill method")
// 			},
// 			FiniFunc: func()  {
// 				panic("mock out the Fini method")
// 			},
// 			GetContentFunc: func(x int, y int) (rune, []rune, tcell.Style, int) {
// 				panic("mock out the GetContent method")
// 			},
// 			HasKeyFunc: func(key tcell.Key) bool {
// 				panic("mock out the HasKey method")
// 			},
// 			HasMouseFunc: func() bool {
// 				panic("mock out the HasMouse method")
// 			},
// 			HasPendingEventFunc: func() bool {
// 				panic("mock out the HasPendingEvent method")
// 			},
// 			HideCursorFunc: func()  {
// 				panic("mock out the HideCursor method")
// 			},
// 			InitFunc: func() error {
// 				panic("mock out the Init method")
// 			},
// 			PollEventFunc: func() tcell.Event {
// 				panic("mock out the PollEvent method")
// 			},
// 			PostEventFunc: func(ev tcell.Event) error {
// 				panic("mock out the PostEvent method")
// 			},
// 			PostEventWaitFunc: func(ev tcell.Event)  {
// 				panic("mock out the PostEventWait method")
// 			},
// 			RegisterRuneFallbackFunc: func(r rune, subst string)  {
// 				panic("mock out the RegisterRuneFallback method")
// 			},
// 			ResizeFunc: func(n1 int, n2 int, n3 int, n4 int)  {
// 				panic("mock out the Resize method")
// 			},
// 			ResumeFunc: func() error {
// 				panic("mock out the Resume method")
// 			},
// 			SetCellFunc: func(x int, y int, style tcell.Style, ch ...rune)  {
// 				panic("mock out the SetCell method")
// 			},
// 			SetContentFunc: func(x int, y int, mainc rune, combc []rune, style tcell.Style)  {
// 				panic("mock out the SetContent method")
// 			},
// 			SetCursorStyleFunc: func(cursorStyle tcell.CursorStyle)  {
// 				panic("mock out the SetCursorStyle method")
// 			},
// 			SetStyleFunc: func(style tcell.Style)  {
// 				panic("mock out the SetStyle method")
// 			},
// 			ShowFunc: func()  {
// 				panic("mock out the Show method")
// 			},
// 			ShowCursorFunc: func(x int, y int)  {
// 				panic("mock out the ShowCursor method")
// 			},
// 			SizeFunc: func() (int, int) {
// 				panic("mock out the Size method")
// 			},
// 			SuspendFunc: func() error {
// 				panic("mock out the Suspend method")
// 			},
// 			SyncFunc: func()  {
// 				panic("mock out the Sync method")
// 			},
// 			UnregisterRuneFallbackFunc: func(r rune)  {
// 				panic("mock out the UnregisterRuneFallback method")
// 			},
// 		}
//
// 		// use mockedScreen in code that requires tcell.Screen
// 		// and then make assertions.
//
// 	}
type ScreenMock struct {
	// BeepFunc mocks the Beep method.
	BeepFunc func() error

	// CanDisplayFunc mocks the CanDisplay method.
	CanDisplayFunc func(r rune, checkFallbacks bool) bool

	// ChannelEventsFunc mocks the ChannelEvents method.
	ChannelEventsFunc func(ch chan<- tcell.Event, quit <-chan struct{})

	// CharacterSetFunc mocks the CharacterSet method.
	CharacterSetFunc func() string

	// ClearFunc mocks the Clear method.
	ClearFunc func()

	// ColorsFunc mocks the Colors method.
	ColorsFunc func() int

	// DisableMouseFunc mocks the DisableMouse method.
	DisableMouseFunc func()

	// DisablePasteFunc mocks the DisablePaste method.
	DisablePasteFunc func()

	// EnableMouseFunc mocks the EnableMouse method.
	EnableMouseFunc func(mouseFlagss ...tcell.MouseFlags)

	// EnablePasteFunc mocks the EnablePaste method.
	EnablePasteFunc func()

	// FillFunc mocks the Fill method.
	FillFunc func(n rune, style tcell.Style)

	// FiniFunc mocks the Fini method.
	FiniFunc func()

	// GetContentFunc mocks the GetContent method.
	GetContentFunc func(x int, y int) (rune, []rune, tcell.Style, int)

	// HasKeyFunc mocks the HasKey method.
	HasKeyFunc func(key tcell.Key) bool

	// HasMouseFunc mocks the HasMouse method.
	HasMouseFunc func() bool

	// HasPendingEventFunc mocks the HasPendingEvent method.
	HasPendingEventFunc func() bool

	// HideCursorFunc mocks the HideCursor method.
	HideCursorFunc func()

	// InitFunc mocks the Init method.
	InitFunc func() error

	// PollEventFunc mocks the PollEvent method.
	PollEventFunc func() tcell.Event

	// PostEventFunc mocks the PostEvent method.
	PostEventFunc func(ev tcell.Event) error

	// PostEventWaitFunc mocks the PostEventWait method.
	PostEventWaitFunc func(ev tcell.Event)

	// RegisterRuneFallbackFunc mocks the RegisterRuneFallback method.
	RegisterRuneFallbackFunc func(r rune, subst string)

	// ResizeFunc mocks the Resize method.
	ResizeFunc func(n1 int, n2 int, n3 int, n4 int)

	// ResumeFunc mocks the Resume method.
	ResumeFunc func() error

	// SetCellFunc mocks the SetCell method.
	SetCellFunc func(x int, y int, style tcell.Style, ch ...rune)

	// SetContentFunc mocks the SetContent method.
	SetContentFunc func(x int, y int, mainc rune, combc []rune, style tcell.Style)

	// SetCursorStyleFunc mocks the SetCursorStyle method.
	SetCursorStyleFunc func(cursorStyle tcell.CursorStyle)

	// SetStyleFunc mocks the SetStyle method.
	SetStyleFunc func(style tcell.Style)

	// ShowFunc mocks the Show method.
	ShowFunc func()

	// ShowCursorFunc mocks the ShowCursor method.
	ShowCursorFunc func(x int, y int)

	// SizeFunc mocks the Size method.
	SizeFunc func() (int, int)

	// SuspendFunc mocks the Suspend method.
	SuspendFunc func() error

	// SyncFunc mocks the Sync method.
	SyncFunc func()

	// UnregisterRuneFallbackFunc mocks the UnregisterRuneFallback method.
	UnregisterRuneFallbackFunc func(r rune)

	// calls tracks calls to the methods.
	calls struct {
		// Beep holds details about calls to the Beep method.
		Beep []struct {
		}
		// CanDisplay holds details about calls to the CanDisplay method.
		CanDisplay []struct {
			// R is the r argument value.
			R rune
			// CheckFallbacks is the checkFallbacks argument value.
			CheckFallbacks bool
		}
		// ChannelEvents holds details about calls to the ChannelEvents method.
		ChannelEvents []struct {
			// Ch is the ch argument value.
			Ch chan<- tcell.Event
			// Quit is the quit argument value.
			Quit <-chan struct{}
		}
		// CharacterSet holds details about calls to the CharacterSet method.
		CharacterSet []struct {
		}
		// Clear holds details about calls to the Clear method.
		Clear []struct {
		}
		// Colors holds details about calls to the Colors method.
		Colors []struct {
		}
		// DisableMouse holds details about calls to the DisableMouse method.
		DisableMouse []struct {
		}
		// DisablePaste holds details about calls to the DisablePaste method.
		DisablePaste []struct {
		}
		// EnableMouse holds details about calls to the EnableMouse method.
		EnableMouse []struct {
			// MouseFlagss is the mouseFlagss argument value.
			MouseFlagss []tcell.MouseFlags
		}
		// EnablePaste holds details about calls to the EnablePaste method.
		EnablePaste []struct {
		}
		// Fill holds details about calls to the Fill method.
		Fill []struct {
			// N is the n argument value.
			N rune
			// Style is the style argument value.
			Style tcell.Style
		}
		// Fini holds details about calls to the Fini method.
		Fini []struct {
		}
		// GetContent holds details about calls to the GetContent method.
		GetContent []struct {
			// X is the x argument value.
			X int
			// Y is the y argument value.
			Y int
		}
		// HasKey holds details about calls to the HasKey method.
		HasKey []struct {
			// Key is the key argument value.
			Key tcell.Key
		}
		// HasMouse holds details about calls to the HasMouse method.
		HasMouse []struct {
		}
		// HasPendingEvent holds details about calls to the HasPendingEvent method.
		HasPendingEvent []struct {
		}
		// HideCursor holds details about calls to the HideCursor method.
		HideCursor []struct {
		}
		// Init holds details about calls to the Init method.
		Init []struct {
		}
		// PollEvent holds details about calls to the PollEvent method.
		PollEvent []struct {
		}
		// PostEvent holds details about calls to the PostEvent method.
		PostEvent []struct {
			// Ev is the ev argument value.
			Ev tcell.Event
		}
		// PostEventWait holds details about calls to the PostEventWait method.
		PostEventWait []struct {
			// Ev is the ev argument value.
			Ev tcell.Event
		}
		// RegisterRuneFallback holds details about calls to the RegisterRuneFallback method.
		RegisterRuneFallback []struct {
			// R is the r argument value.
			R rune
			// Subst is the subst argument value.
			Subst string
		}
		// Resize holds details about calls to the Resize method.
		Resize []struct {
			// N1 is the n1 argument value.
			N1 int
			// N2 is the n2 argument value.
			N2 int
			// N3 is the n3 argument value.
			N3 int
			// N4 is the n4 argument value.
			N4 int
		}
		// Resume holds details about calls to the Resume method.
		Resume []struct {
		}
		// SetCell holds details about calls to the SetCell method.
		SetCell []struct {
			// X is the x argument value.
			X int
			// Y is the y argument value.
			Y int
			// Style is the style argument value.
			Style tcell.Style
			// Ch is the ch argument value.
			Ch []rune
		}
		// SetContent holds details about calls to the SetContent method.
		SetContent []struct {
			// X is the x argument value.
			X int
			// Y is the y argument value.
			Y int
			// Mainc is the mainc argument value.
			Mainc rune
			// Combc is the combc argument value.
			Combc []rune
			// Style is the style argument value.
			Style tcell.Style
		}
		// SetCursorStyle holds details about calls to the SetCursorStyle method.
		SetCursorStyle []struct {
			// CursorStyle is the cursorStyle argument value.
			CursorStyle tcell.CursorStyle
		}
		// SetStyle holds details about calls to the SetStyle method.
		SetStyle []struct {
			// Style is the style argument value.
			Style tcell.Style
		}
		// Show holds details about calls to the Show method.
		Show []struct {
		}
		// ShowCursor holds details about calls to the ShowCursor method.
		ShowCursor []struct {
			// X is the x argument value.
			X int
			// Y is the y argument value.
			Y int
		}
		// Size holds details about calls to the Size method.
		Size []struct {
		}
		// Suspend holds details about calls to the Suspend method.
		Suspend []struct {
		}
		// Sync holds details about calls to the Sync method.
		Sync []struct {
		}
		// UnregisterRuneFallback holds details about calls to the UnregisterRuneFallback method.
		UnregisterRuneFallback []struct {
			// R is the r argument value.
			R rune
		}
	}
	lockBeep                   sync.RWMutex
	lockCanDisplay             sync.RWMutex
	lockChannelEvents          sync.RWMutex
	lockCharacterSet           sync.RWMutex
	lockClear                  sync.RWMutex
	lockColors                 sync.RWMutex
	lockDisableMouse           sync.RWMutex
	lockDisablePaste           sync.RWMutex
	lockEnableMouse            sync.RWMutex
	lockEnablePaste            sync.RWMutex
	lockFill                   sync.RWMutex
	lockFini                   sync.RWMutex
	lockGetContent             sync.RWMutex
	lockHasKey                 sync.RWMutex
	lockHasMouse               sync.RWMutex
	lockHasPendingEvent        sync.RWMutex
	lockHideCursor             sync.RWMutex
	lockInit                   sync.RWMutex
	lockPollEvent              sync.RWMutex
	lockPostEvent              sync.RWMutex
	lockPostEventWait          sync.RWMutex
	lockRegisterRuneFallback   sync.RWMutex
	lockResize                 sync.RWMutex
	lockResume                 sync.RWMutex
	lockSetCell                sync.RWMutex
	lockSetContent             sync.RWMutex
	lockSetCursorStyle         sync.RWMutex
	lockSetStyle               sync.RWMutex
	lockShow                   sync.RWMutex
	lockShowCursor             sync.RWMutex
	lockSize                   sync.RWMutex
	lockSuspend                sync.RWMutex
	lockSync                   sync.RWMutex
	lockUnregisterRuneFallback sync.RWMutex
}

// Beep calls BeepFunc.
func (mock *ScreenMock) Beep() error {
	if mock.BeepFunc == nil {
		panic("ScreenMock.BeepFunc: method is nil but Screen.Beep was just called")
	}
	callInfo := struct {
	}{}
	mock.lockBeep.Lock()
	mock.calls.Beep = append(mock.calls.Beep, callInfo)
	mock.lockBeep.Unlock()
	return mock.BeepFunc()
}

// BeepCalls gets all the calls that were made to Beep.
// Check the length with:
//     len(mockedScreen.BeepCalls())
func (mock *ScreenMock) BeepCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockBeep.RLock()
	calls = mock.calls.Beep
	mock.lockBeep.RUnlock()
	return calls
}

// CanDisplay calls CanDisplayFunc.
func (mock *ScreenMock) CanDisplay(r rune, checkFallbacks bool) bool {
	if mock.CanDisplayFunc == nil {
		panic("ScreenMock.CanDisplayFunc: method is nil but Screen.CanDisplay was just called")
	}
	callInfo := struct {
		R              rune
		CheckFallbacks bool
	}{
		R:              r,
		CheckFallbacks: checkFallbacks,
	}
	mock.lockCanDisplay.Lock()
	mock.calls.CanDisplay = append(mock.calls.CanDisplay, callInfo)
	mock.lockCanDisplay.Unlock()
	return mock.CanDisplayFunc(r, checkFallbacks)
}

// CanDisplayCalls gets all the calls that were made to CanDisplay.
// Check the length with:
//     len(mockedScreen.CanDisplayCalls())
func (mock *ScreenMock) CanDisplayCalls() []struct {
	R              rune
	CheckFallbacks bool
} {
	var calls []struct {
		R              rune
		CheckFallbacks bool
	}
	mock.lockCanDisplay.RLock()
	calls = mock.calls.CanDisplay
	mock.lockCanDisplay.RUnlock()
	return calls
}

// ChannelEvents calls ChannelEventsFunc.
func (mock *ScreenMock) ChannelEvents(ch chan<- tcell.Event, quit <-chan struct{}) {
	if mock.ChannelEventsFunc == nil {
		panic("ScreenMock.ChannelEventsFunc: method is nil but Screen.ChannelEvents was just called")
	}
	callInfo := struct {
		Ch   chan<- tcell.Event
		Quit <-chan struct{}
	}{
		Ch:   ch,
		Quit: quit,
	}
	mock.lockChannelEvents.Lock()
	mock.calls.ChannelEvents = append(mock.calls.ChannelEvents, callInfo)
	mock.lockChannelEvents.Unlock()
	mock.ChannelEventsFunc(ch, quit)
}

// ChannelEventsCalls gets all the calls that were made to ChannelEvents.
// Check the length with:
//     len(mockedScreen.ChannelEventsCalls())
func (mock *ScreenMock) ChannelEventsCalls() []struct {
	Ch   chan<- tcell.Event
	Quit <-chan struct{}
} {
	var calls []struct {
		Ch   chan<- tcell.Event
		Quit <-chan struct{}
	}
	mock.lockChannelEvents.RLock()
	calls = mock.calls.ChannelEvents
	mock.lockChannelEvents.RUnlock()
	return calls
}

// CharacterSet calls CharacterSetFunc.
func (mock *ScreenMock) CharacterSet() string {
	if mock.CharacterSetFunc == nil {
		panic("ScreenMock.CharacterSetFunc: method is nil but Screen.CharacterSet was just called")
	}
	callInfo := struct {
	}{}
	mock.lockCharacterSet.Lock()
	mock.calls.CharacterSet = append(mock.calls.CharacterSet, callInfo)
	mock.lockCharacterSet.Unlock()
	return mock.CharacterSetFunc()
}

// CharacterSetCalls gets all the calls that were made to CharacterSet.
// Check the length with:
//     len(mockedScreen.CharacterSetCalls())
func (mock *ScreenMock) CharacterSetCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockCharacterSet.RLock()
	calls = mock.calls.CharacterSet
	mock.lockCharacterSet.RUnlock()
	return calls
}

// Clear calls ClearFunc.
func (mock *ScreenMock) Clear() {
	if mock.ClearFunc == nil {
		panic("ScreenMock.ClearFunc: method is nil but Screen.Clear was just called")
	}
	callInfo := struct {
	}{}
	mock.lockClear.Lock()
	mock.calls.Clear = append(mock.calls.Clear, callInfo)
	mock.lockClear.Unlock()
	mock.ClearFunc()
}

// ClearCalls gets all the calls that were made to Clear.
// Check the length with:
//     len(mockedScreen.ClearCalls())
func (mock *ScreenMock) ClearCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockClear.RLock()
	calls = mock.calls.Clear
	mock.lockClear.RUnlock()
	return calls
}

// Colors calls ColorsFunc.
func (mock *ScreenMock) Colors() int {
	if mock.ColorsFunc == nil {
		panic("ScreenMock.ColorsFunc: method is nil but Screen.Colors was just called")
	}
	callInfo := struct {
	}{}
	mock.lockColors.Lock()
	mock.calls.Colors = append(mock.calls.Colors, callInfo)
	mock.lockColors.Unlock()
	return mock.ColorsFunc()
}

// ColorsCalls gets all the calls that were made to Colors.
// Check the length with:
//     len(mockedScreen.ColorsCalls())
func (mock *ScreenMock) ColorsCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockColors.RLock()
	calls = mock.calls.Colors
	mock.lockColors.RUnlock()
	return calls
}

// DisableMouse calls DisableMouseFunc.
func (mock *ScreenMock) DisableMouse() {
	if mock.DisableMouseFunc == nil {
		panic("ScreenMock.DisableMouseFunc: method is nil but Screen.DisableMouse was just called")
	}
	callInfo := struct {
	}{}
	mock.lockDisableMouse.Lock()
	mock.calls.DisableMouse = append(mock.calls.DisableMouse, callInfo)
	mock.lockDisableMouse.Unlock()
	mock.DisableMouseFunc()
}

// DisableMouseCalls gets all the calls that were made to DisableMouse.
// Check the length with:
//     len(mockedScreen.DisableMouseCalls())
func (mock *ScreenMock) DisableMouseCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockDisableMouse.RLock()
	calls = mock.calls.DisableMouse
	mock.lockDisableMouse.RUnlock()
	return calls
}

// DisablePaste calls DisablePasteFunc.
func (mock *ScreenMock) DisablePaste() {
	if mock.DisablePasteFunc == nil {
		panic("ScreenMock.DisablePasteFunc: method is nil but Screen.DisablePaste was just called")
	}
	callInfo := struct {
	}{}
	mock.lockDisablePaste.Lock()
	mock.calls.DisablePaste = append(mock.calls.DisablePaste, callInfo)
	mock.lockDisablePaste.Unlock()
	mock.DisablePasteFunc()
}

// DisablePasteCalls gets all the calls that were made to DisablePaste.
// Check the length with:
//     len(mockedScreen.DisablePasteCalls())
func (mock *ScreenMock) DisablePasteCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockDisablePaste.RLock()
	calls = mock.calls.DisablePaste
	mock.lockDisablePaste.RUnlock()
	return calls
}

// EnableMouse calls EnableMouseFunc.
func (mock *ScreenMock) EnableMouse(mouseFlagss ...tcell.MouseFlags) {
	if mock.EnableMouseFunc == nil {
		panic("ScreenMock.EnableMouseFunc: method is nil but Screen.EnableMouse was just called")
	}
	callInfo := struct {
		MouseFlagss []tcell.MouseFlags
	}{
		MouseFlagss: mouseFlagss,
	}
	mock.lockEnableMouse.Lock()
	mock.calls.EnableMouse = append(mock.calls.EnableMouse, callInfo)
	mock.lockEnableMouse.Unlock()
	mock.EnableMouseFunc(mouseFlagss...)
}

// EnableMouseCalls gets all the calls that were made to EnableMouse.
// Check the length with:
//     len(mockedScreen.EnableMouseCalls())
func (mock *ScreenMock) EnableMouseCalls() []struct {
	MouseFlagss []tcell.MouseFlags
} {
	var calls []struct {
		MouseFlagss []tcell.MouseFlags
	}
	mock.lockEnableMouse.RLock()
	calls = mock.calls.EnableMouse
	mock.lockEnableMouse.RUnlock()
	return calls
}

// EnablePaste calls EnablePasteFunc.
func (mock *ScreenMock) EnablePaste() {
	if mock.EnablePasteFunc == nil {
		panic("ScreenMock.EnablePasteFunc: method is nil but Screen.EnablePaste was just called")
	}
	callInfo := struct {
	}{}
	mock.lockEnablePaste.Lock()
	mock.calls.EnablePaste = append(mock.calls.EnablePaste, callInfo)
	mock.lockEnablePaste.Unlock()
	mock.EnablePasteFunc()
}

// EnablePasteCalls gets all the calls that were made to EnablePaste.
// Check the length with:
//     len(mockedScreen.EnablePasteCalls())
func (mock *ScreenMock) EnablePasteCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockEnablePaste.RLock()
	calls = mock.calls.EnablePaste
	mock.lockEnablePaste.RUnlock()
	return calls
}

// Fill calls FillFunc.
func (mock *ScreenMock) Fill(n rune, style tcell.Style) {
	if mock.FillFunc == nil {
		panic("ScreenMock.FillFunc: method is nil but Screen.Fill was just called")
	}
	callInfo := struct {
		N     rune
		Style tcell.Style
	}{
		N:     n,
		Style: style,
	}
	mock.lockFill.Lock()
	mock.calls.Fill = append(mock.calls.Fill, callInfo)
	mock.lockFill.Unlock()
	mock.FillFunc(n, style)
}

// FillCalls gets all the calls that were made to Fill.
// Check the length with:
//     len(mockedScreen.FillCalls())
func (mock *ScreenMock) FillCalls() []struct {
	N     rune
	Style tcell.Style
} {
	var calls []struct {
		N     rune
		Style tcell.Style
	}
	mock.lockFill.RLock()
	calls = mock.calls.Fill
	mock.lockFill.RUnlock()
	return calls
}

// Fini calls FiniFunc.
func (mock *ScreenMock) Fini() {
	if mock.FiniFunc == nil {
		panic("ScreenMock.FiniFunc: method is nil but Screen.Fini was just called")
	}
	callInfo := struct {
	}{}
	mock.lockFini.Lock()
	mock.calls.Fini = append(mock.calls.Fini, callInfo)
	mock.lockFini.Unlock()
	mock.FiniFunc()
}

// FiniCalls gets all the calls that were made to Fini.
// Check the length with:
//     len(mockedScreen.FiniCalls())
func (mock *ScreenMock) FiniCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockFini.RLock()
	calls = mock.calls.Fini
	mock.lockFini.RUnlock()
	return calls
}

// GetContent calls GetContentFunc.
func (mock *ScreenMock) GetContent(x int, y int) (rune, []rune, tcell.Style, int) {
	if mock.GetContentFunc == nil {
		panic("ScreenMock.GetContentFunc: method is nil but Screen.GetContent was just called")
	}
	callInfo := struct {
		X int
		Y int
	}{
		X: x,
		Y: y,
	}
	mock.lockGetContent.Lock()
	mock.calls.GetContent = append(mock.calls.GetContent, callInfo)
	mock.lockGetContent.Unlock()
	return mock.GetContentFunc(x, y)
}

// GetContentCalls gets all the calls that were made to GetContent.
// Check the length with:
//     len(mockedScreen.GetContentCalls())
func (mock *ScreenMock) GetContentCalls() []struct {
	X int
	Y int
} {
	var calls []struct {
		X int
		Y int
	}
	mock.lockGetContent.RLock()
	calls = mock.calls.GetContent
	mock.lockGetContent.RUnlock()
	return calls
}

// HasKey calls HasKeyFunc.
func (mock *ScreenMock) HasKey(key tcell.Key) bool {
	if mock.HasKeyFunc == nil {
		panic("ScreenMock.HasKeyFunc: method is nil but Screen.HasKey was just called")
	}
	callInfo := struct {
		Key tcell.Key
	}{
		Key: key,
	}
	mock.lockHasKey.Lock()
	mock.calls.HasKey = append(mock.calls.HasKey, callInfo)
	mock.lockHasKey.Unlock()
	return mock.HasKeyFunc(key)
}

// HasKeyCalls gets all the calls that were made to HasKey.
// Check the length with:
//     len(mockedScreen.HasKeyCalls())
func (mock *ScreenMock) HasKeyCalls() []struct {
	Key tcell.Key
} {
	var calls []struct {
		Key tcell.Key
	}
	mock.lockHasKey.RLock()
	calls = mock.calls.HasKey
	mock.lockHasKey.RUnlock()
	return calls
}

// HasMouse calls HasMouseFunc.
func (mock *ScreenMock) HasMouse() bool {
	if mock.HasMouseFunc == nil {
		panic("ScreenMock.HasMouseFunc: method is nil but Screen.HasMouse was just called")
	}
	callInfo := struct {
	}{}
	mock.lockHasMouse.Lock()
	mock.calls.HasMouse = append(mock.calls.HasMouse, callInfo)
	mock.lockHasMouse.Unlock()
	return mock.HasMouseFunc()
}

// HasMouseCalls gets all the calls that were made to HasMouse.
// Check the length with:
//     len(mockedScreen.HasMouseCalls())
func (mock *ScreenMock) HasMouseCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockHasMouse.RLock()
	calls = mock.calls.HasMouse
	mock.lockHasMouse.RUnlock()
	return calls
}

// HasPendingEvent calls HasPendingEventFunc.
func (mock *ScreenMock) HasPendingEvent() bool {
	if mock.HasPendingEventFunc == nil {
		panic("ScreenMock.HasPendingEventFunc: method is nil but Screen.HasPendingEvent was just called")
	}
	callInfo := struct {
	}{}
	mock.lockHasPendingEvent.Lock()
	mock.calls.HasPendingEvent = append(mock.calls.HasPendingEvent, callInfo)
	mock.lockHasPendingEvent.Unlock()
	return mock.HasPendingEventFunc()
}

// HasPendingEventCalls gets all the calls that were made to HasPendingEvent.
// Check the length with:
//     len(mockedScreen.HasPendingEventCalls())
func (mock *ScreenMock) HasPendingEventCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockHasPendingEvent.RLock()
	calls = mock.calls.HasPendingEvent
	mock.lockHasPendingEvent.RUnlock()
	return calls
}

// HideCursor calls HideCursorFunc.
func (mock *ScreenMock) HideCursor() {
	if mock.HideCursorFunc == nil {
		panic("ScreenMock.HideCursorFunc: method is nil but Screen.HideCursor was just called")
	}
	callInfo := struct {
	}{}
	mock.lockHideCursor.Lock()
	mock.calls.HideCursor = append(mock.calls.HideCursor, callInfo)
	mock.lockHideCursor.Unlock()
	mock.HideCursorFunc()
}

// HideCursorCalls gets all the calls that were made to HideCursor.
// Check the length with:
//     len(mockedScreen.HideCursorCalls())
func (mock *ScreenMock) HideCursorCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockHideCursor.RLock()
	calls = mock.calls.HideCursor
	mock.lockHideCursor.RUnlock()
	return calls
}

// Init calls InitFunc.
func (mock *ScreenMock) Init() error {
	if mock.InitFunc == nil {
		panic("ScreenMock.InitFunc: method is nil but Screen.Init was just called")
	}
	callInfo := struct {
	}{}
	mock.lockInit.Lock()
	mock.calls.Init = append(mock.calls.Init, callInfo)
	mock.lockInit.Unlock()
	return mock.InitFunc()
}

// InitCalls gets all the calls that were made to Init.
// Check the length with:
//     len(mockedScreen.InitCalls())
func (mock *ScreenMock) InitCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockInit.RLock()
	calls = mock.calls.Init
	mock.lockInit.RUnlock()
	return calls
}

// PollEvent calls PollEventFunc.
func (mock *ScreenMock) PollEvent() tcell.Event {
	if mock.PollEventFunc == nil {
		panic("ScreenMock.PollEventFunc: method is nil but Screen.PollEvent was just called")
	}
	callInfo := struct {
	}{}
	mock.lockPollEvent.Lock()
	mock.calls.PollEvent = append(mock.calls.PollEvent, callInfo)
	mock.lockPollEvent.Unlock()
	return mock.PollEventFunc()
}

// PollEventCalls gets all the calls that were made to PollEvent.
// Check the length with:
//     len(mockedScreen.PollEventCalls())
func (mock *ScreenMock) PollEventCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockPollEvent.RLock()
	calls = mock.calls.PollEvent
	mock.lockPollEvent.RUnlock()
	return calls
}

// PostEvent calls PostEventFunc.
func (mock *ScreenMock) PostEvent(ev tcell.Event) error {
	if mock.PostEventFunc == nil {
		panic("ScreenMock.PostEventFunc: method is nil but Screen.PostEvent was just called")
	}
	callInfo := struct {
		Ev tcell.Event
	}{
		Ev: ev,
	}
	mock.lockPostEvent.Lock()
	mock.calls.PostEvent = append(mock.calls.PostEvent, callInfo)
	mock.lockPostEvent.Unlock()
	return mock.PostEventFunc(ev)
}

// PostEventCalls gets all the calls that were made to PostEvent.
// Check the length with:
//     len(mockedScreen.PostEventCalls())
func (mock *ScreenMock) PostEventCalls() []struct {
	Ev tcell.Event
} {
	var calls []struct {
		Ev tcell.Event
	}
	mock.lockPostEvent.RLock()
	calls = mock.calls.PostEvent
	mock.lockPostEvent.RUnlock()
	return calls
}

// PostEventWait calls PostEventWaitFunc.
func (mock *ScreenMock) PostEventWait(ev tcell.Event) {
	if mock.PostEventWaitFunc == nil {
		panic("ScreenMock.PostEventWaitFunc: method is nil but Screen.PostEventWait was just called")
	}
	callInfo := struct {
		Ev tcell.Event
	}{
		Ev: ev,
	}
	mock.lockPostEventWait.Lock()
	mock.calls.PostEventWait = append(mock.calls.PostEventWait, callInfo)
	mock.lockPostEventWait.Unlock()
	mock.PostEventWaitFunc(ev)
}

// PostEventWaitCalls gets all the calls that were made to PostEventWait.
// Check the length with:
//     len(mockedScreen.PostEventWaitCalls())
func (mock *ScreenMock) PostEventWaitCalls() []struct {
	Ev tcell.Event
} {
	var calls []struct {
		Ev tcell.Event
	}
	mock.lockPostEventWait.RLock()
	calls = mock.calls.PostEventWait
	mock.lockPostEventWait.RUnlock()
	return calls
}

// RegisterRuneFallback calls RegisterRuneFallbackFunc.
func (mock *ScreenMock) RegisterRuneFallback(r rune, subst string) {
	if mock.RegisterRuneFallbackFunc == nil {
		panic("ScreenMock.RegisterRuneFallbackFunc: method is nil but Screen.RegisterRuneFallback was just called")
	}
	callInfo := struct {
		R     rune
		Subst string
	}{
		R:     r,
		Subst: subst,
	}
	mock.lockRegisterRuneFallback.Lock()
	mock.calls.RegisterRuneFallback = append(mock.calls.RegisterRuneFallback, callInfo)
	mock.lockRegisterRuneFallback.Unlock()
	mock.RegisterRuneFallbackFunc(r, subst)
}

// RegisterRuneFallbackCalls gets all the calls that were made to RegisterRuneFallback.
// Check the length with:
//     len(mockedScreen.RegisterRuneFallbackCalls())
func (mock *ScreenMock) RegisterRuneFallbackCalls() []struct {
	R     rune
	Subst string
} {
	var calls []struct {
		R     rune
		Subst string
	}
	mock.lockRegisterRuneFallback.RLock()
	calls = mock.calls.RegisterRuneFallback
	mock.lockRegisterRuneFallback.RUnlock()
	return calls
}

// Resize calls ResizeFunc.
func (mock *ScreenMock) Resize(n1 int, n2 int, n3 int, n4 int) {
	if mock.ResizeFunc == nil {
		panic("ScreenMock.ResizeFunc: method is nil but Screen.Resize was just called")
	}
	callInfo := struct {
		N1 int
		N2 int
		N3 int
		N4 int
	}{
		N1: n1,
		N2: n2,
		N3: n3,
		N4: n4,
	}
	mock.lockResize.Lock()
	mock.calls.Resize = append(mock.calls.Resize, callInfo)
	mock.lockResize.Unlock()
	mock.ResizeFunc(n1, n2, n3, n4)
}

// ResizeCalls gets all the calls that were made to Resize.
// Check the length with:
//     len(mockedScreen.ResizeCalls())
func (mock *ScreenMock) ResizeCalls() []struct {
	N1 int
	N2 int
	N3 int
	N4 int
} {
	var calls []struct {
		N1 int
		N2 int
		N3 int
		N4 int
	}
	mock.lockResize.RLock()
	calls = mock.calls.Resize
	mock.lockResize.RUnlock()
	return calls
}

// Resume calls ResumeFunc.
func (mock *ScreenMock) Resume() error {
	if mock.ResumeFunc == nil {
		panic("ScreenMock.ResumeFunc: method is nil but Screen.Resume was just called")
	}
	callInfo := struct {
	}{}
	mock.lockResume.Lock()
	mock.calls.Resume = append(mock.calls.Resume, callInfo)
	mock.lockResume.Unlock()
	return mock.ResumeFunc()
}

// ResumeCalls gets all the calls that were made to Resume.
// Check the length with:
//     len(mockedScreen.ResumeCalls())
func (mock *ScreenMock) ResumeCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockResume.RLock()
	calls = mock.calls.Resume
	mock.lockResume.RUnlock()
	return calls
}

// SetCell calls SetCellFunc.
func (mock *ScreenMock) SetCell(x int, y int, style tcell.Style, ch ...rune) {
	if mock.SetCellFunc == nil {
		panic("ScreenMock.SetCellFunc: method is nil but Screen.SetCell was just called")
	}
	callInfo := struct {
		X     int
		Y     int
		Style tcell.Style
		Ch    []rune
	}{
		X:     x,
		Y:     y,
		Style: style,
		Ch:    ch,
	}
	mock.lockSetCell.Lock()
	mock.calls.SetCell = append(mock.calls.SetCell, callInfo)
	mock.lockSetCell.Unlock()
	mock.SetCellFunc(x, y, style, ch...)
}

// SetCellCalls gets all the calls that were made to SetCell.
// Check the length with:
//     len(mockedScreen.SetCellCalls())
func (mock *ScreenMock) SetCellCalls() []struct {
	X     int
	Y     int
	Style tcell.Style
	Ch    []rune
} {
	var calls []struct {
		X     int
		Y     int
		Style tcell.Style
		Ch    []rune
	}
	mock.lockSetCell.RLock()
	calls = mock.calls.SetCell
	mock.lockSetCell.RUnlock()
	return calls
}

// SetContent calls SetContentFunc.
func (mock *ScreenMock) SetContent(x int, y int, mainc rune, combc []rune, style tcell.Style) {
	if mock.SetContentFunc == nil {
		panic("ScreenMock.SetContentFunc: method is nil but Screen.SetContent was just called")
	}
	callInfo := struct {
		X     int
		Y     int
		Mainc rune
		Combc []rune
		Style tcell.Style
	}{
		X:     x,
		Y:     y,
		Mainc: mainc,
		Combc: combc,
		Style: style,
	}
	mock.lockSetContent.Lock()
	mock.calls.SetContent = append(mock.calls.SetContent, callInfo)
	mock.lockSetContent.Unlock()
	mock.SetContentFunc(x, y, mainc, combc, style)
}

// SetContentCalls gets all the calls that were made to SetContent.
// Check the length with:
//     len(mockedScreen.SetContentCalls())
func (mock *ScreenMock) SetContentCalls() []struct {
	X     int
	Y     int
	Mainc rune
	Combc []rune
	Style tcell.Style
} {
	var calls []struct {
		X     int
		Y     int
		Mainc rune
		Combc []rune
		Style tcell.Style
	}
	mock.lockSetContent.RLock()
	calls = mock.calls.SetContent
	mock.lockSetContent.RUnlock()
	return calls
}

// SetCursorStyle calls SetCursorStyleFunc.
func (mock *ScreenMock) SetCursorStyle(cursorStyle tcell.CursorStyle) {
	if mock.SetCursorStyleFunc == nil {
		panic("ScreenMock.SetCursorStyleFunc: method is nil but Screen.SetCursorStyle was just called")
	}
	callInfo := struct {
		CursorStyle tcell.CursorStyle
	}{
		CursorStyle: cursorStyle,
	}
	mock.lockSetCursorStyle.Lock()
	mock.calls.SetCursorStyle = append(mock.calls.SetCursorStyle, callInfo)
	mock.lockSetCursorStyle.Unlock()
	mock.SetCursorStyleFunc(cursorStyle)
}

// SetCursorStyleCalls gets all the calls that were made to SetCursorStyle.
// Check the length with:
//     len(mockedScreen.SetCursorStyleCalls())
func (mock *ScreenMock) SetCursorStyleCalls() []struct {
	CursorStyle tcell.CursorStyle
} {
	var calls []struct {
		CursorStyle tcell.CursorStyle
	}
	mock.lockSetCursorStyle.RLock()
	calls = mock.calls.SetCursorStyle
	mock.lockSetCursorStyle.RUnlock()
	return calls
}

// SetStyle calls SetStyleFunc.
func (mock *ScreenMock) SetStyle(style tcell.Style) {
	if mock.SetStyleFunc == nil {
		panic("ScreenMock.SetStyleFunc: method is nil but Screen.SetStyle was just called")
	}
	callInfo := struct {
		Style tcell.Style
	}{
		Style: style,
	}
	mock.lockSetStyle.Lock()
	mock.calls.SetStyle = append(mock.calls.SetStyle, callInfo)
	mock.lockSetStyle.Unlock()
	mock.SetStyleFunc(style)
}

// SetStyleCalls gets all the calls that were made to SetStyle.
// Check the length with:
//     len(mockedScreen.SetStyleCalls())
func (mock *ScreenMock) SetStyleCalls() []struct {
	Style tcell.Style
} {
	var calls []struct {
		Style tcell.Style
	}
	mock.lockSetStyle.RLock()
	calls = mock.calls.SetStyle
	mock.lockSetStyle.RUnlock()
	return calls
}

// Show calls ShowFunc.
func (mock *ScreenMock) Show() {
	if mock.ShowFunc == nil {
		panic("ScreenMock.ShowFunc: method is nil but Screen.Show was just called")
	}
	callInfo := struct {
	}{}
	mock.lockShow.Lock()
	mock.calls.Show = append(mock.calls.Show, callInfo)
	mock.lockShow.Unlock()
	mock.ShowFunc()
}

// ShowCalls gets all the calls that were made to Show.
// Check the length with:
//     len(mockedScreen.ShowCalls())
func (mock *ScreenMock) ShowCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockShow.RLock()
	calls = mock.calls.Show
	mock.lockShow.RUnlock()
	return calls
}

// ShowCursor calls ShowCursorFunc.
func (mock *ScreenMock) ShowCursor(x int, y int) {
	if mock.ShowCursorFunc == nil {
		panic("ScreenMock.ShowCursorFunc: method is nil but Screen.ShowCursor was just called")
	}
	callInfo := struct {
		X int
		Y int
	}{
		X: x,
		Y: y,
	}
	mock.lockShowCursor.Lock()
	mock.calls.ShowCursor = append(mock.calls.ShowCursor, callInfo)
	mock.lockShowCursor.Unlock()
	mock.ShowCursorFunc(x, y)
}

// ShowCursorCalls gets all the calls that were made to ShowCursor.
// Check the length with:
//     len(mockedScreen.ShowCursorCalls())
func (mock *ScreenMock) ShowCursorCalls() []struct {
	X int
	Y int
} {
	var calls []struct {
		X int
		Y int
	}
	mock.lockShowCursor.RLock()
	calls = mock.calls.ShowCursor
	mock.lockShowCursor.RUnlock()
	return calls
}

// Size calls SizeFunc.
func (mock *ScreenMock) Size() (int, int) {
	if mock.SizeFunc == nil {
		panic("ScreenMock.SizeFunc: method is nil but Screen.Size was just called")
	}
	callInfo := struct {
	}{}
	mock.lockSize.Lock()
	mock.calls.Size = append(mock.calls.Size, callInfo)
	mock.lockSize.Unlock()
	return mock.SizeFunc()
}

// SizeCalls gets all the calls that were made to Size.
// Check the length with:
//     len(mockedScreen.SizeCalls())
func (mock *ScreenMock) SizeCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockSize.RLock()
	calls = mock.calls.Size
	mock.lockSize.RUnlock()
	return calls
}

// Suspend calls SuspendFunc.
func (mock *ScreenMock) Suspend() error {
	if mock.SuspendFunc == nil {
		panic("ScreenMock.SuspendFunc: method is nil but Screen.Suspend was just called")
	}
	callInfo := struct {
	}{}
	mock.lockSuspend.Lock()
	mock.calls.Suspend = append(mock.calls.Suspend, callInfo)
	mock.lockSuspend.Unlock()
	return mock.SuspendFunc()
}

// SuspendCalls gets all the calls that were made to Suspend.
// Check the length with:
//     len(mockedScreen.SuspendCalls())
func (mock *ScreenMock) SuspendCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockSuspend.RLock()
	calls = mock.calls.Suspend
	mock.lockSuspend.RUnlock()
	return calls
}

// Sync calls SyncFunc.
func (mock *ScreenMock) Sync() {
	if mock.SyncFunc == nil {
		panic("ScreenMock.SyncFunc: method is nil but Screen.Sync was just called")
	}
	callInfo := struct {
	}{}
	mock.lockSync.Lock()
	mock.calls.Sync = append(mock.calls.Sync, callInfo)
	mock.lockSync.Unlock()
	mock.SyncFunc()
}

// SyncCalls gets all the calls that were made to Sync.
// Check the length with:
//     len(mockedScreen.SyncCalls())
func (mock *ScreenMock) SyncCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockSync.RLock()
	calls = mock.calls.Sync
	mock.lockSync.RUnlock()
	return calls
}

// UnregisterRuneFallback calls UnregisterRuneFallbackFunc.
func (mock *ScreenMock) UnregisterRuneFallback(r rune) {
	if mock.UnregisterRuneFallbackFunc == nil {
		panic("ScreenMock.UnregisterRuneFallbackFunc: method is nil but Screen.UnregisterRuneFallback was just called")
	}
	callInfo := struct {
		R rune
	}{
		R: r,
	}
	mock.lockUnregisterRuneFallback.Lock()
	mock.calls.UnregisterRuneFallback = append(mock.calls.UnregisterRuneFallback, callInfo)
	mock.lockUnregisterRuneFallback.Unlock()
	mock.UnregisterRuneFallbackFunc(r)
}

// UnregisterRuneFallbackCalls gets all the calls that were made to UnregisterRuneFallback.
// Check the length with:
//     len(mockedScreen.UnregisterRuneFallbackCalls())
func (mock *ScreenMock) UnregisterRuneFallbackCalls() []struct {
	R rune
} {
	var calls []struct {
		R rune
	}
	mock.lockUnregisterRuneFallback.RLock()
	calls = mock.calls.UnregisterRuneFallback
	mock.lockUnregisterRuneFallback.RUnlock()
	return calls
}
