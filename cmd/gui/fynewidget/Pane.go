package fynewidget

import (
	"fyne.io/fyne"
	"github.com/AlbinoGeek/sc2-rsu/cmd/gui"
)

type Pane interface {
	GetWindow() gui.Window
	GetContent() fyne.CanvasObject
	SetContent(fyne.CanvasObject)
}

func NewPane(title string, window gui.Window) Pane {
	return &PaneBase{
		title:  title,
		window: window,
	}
}

func NewPaneWithIcon(title string, icon fyne.Resource, window gui.Window) Pane {
	return &PaneBase{
		icon:   icon,
		title:  title,
		window: window,
	}
}

type PaneBase struct {
	content fyne.CanvasObject
	icon    fyne.Resource
	title   string
	window  gui.Window
}

func (p *PaneBase) GetWindow() gui.Window {
	return p.window
}

func (p *PaneBase) GetContent() fyne.CanvasObject {
	return p.content
}

func (p *PaneBase) SetContent(content fyne.CanvasObject) {
	p.content = content
}
