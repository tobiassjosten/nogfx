.PHONY: all
all: pkg/mock/pkg_client.go pkg/mock/pkg_module.go pkg/mock/pkg_ui.go pkg/mock/tcell_screen.go

pkg/mock/pkg_client.go: pkg/services.go
	~/go/bin/moq -pkg mock pkg/ Client:ClientMock > pkg/mock/pkg_client.go

pkg/mock/pkg_module.go: pkg/modules.go
	~/go/bin/moq -pkg mock pkg/ Module:ModuleMock > pkg/mock/pkg_module.go

pkg/mock/pkg_ui.go: pkg/services.go
	~/go/bin/moq -pkg mock pkg/ UI:UIMock > pkg/mock/pkg_ui.go

pkg/mock/pkg_world.go: pkg/services.go
	~/go/bin/moq -pkg mock pkg/ World:WorldMock > pkg/mock/pkg_world.go

pkg/mock/tcell_screen.go: pkg/services.go
	~/go/bin/moq -pkg mock ~/go/pkg/mod/github.com/gdamore/tcell/v2@v2.6.0/ Screen:ScreenMock > pkg/mock/tcell_screen.go
