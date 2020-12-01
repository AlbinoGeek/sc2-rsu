package fynewidget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"github.com/AlbinoGeek/sc2-rsu/cmd/gui"
)

type Pane interface {
	GetWindow() gui.Window

	// GetContent() implements fynewidget.NavigationItem
	GetContent() fyne.CanvasObject
	SetContent(fyne.CanvasObject)

	// GetLabel() implements fynewidget.NavigationItem
	GetLabel() fyne.CanvasObject

	GetIcon() fyne.Resource
	SetIcon(fyne.Resource)
	GetTitle() string
	SetTitle(string)
}

func NewPane(title string, window gui.Window) Pane {
	return NewPaneWithIcon(title, nil, window)
}

func NewPaneWithIcon(title string, icon fyne.Resource, window gui.Window) Pane {
	pane := &PaneBase{
		icon:   icon,
		title:  title,
		window: window,
	}

	pane.label = newNavigationLabel(title, icon)
	return pane
}

func newNavigationLabel(title string, icon fyne.Resource) fyne.CanvasObject {
	b := widget.NewButtonWithIcon(title, icon, nil)
	b.Alignment = widget.ButtonAlignLeading
	b.HideShadow = true
	b.Importance = widget.LowImportance
	return b
}

type PaneBase struct {
	content fyne.CanvasObject
	icon    fyne.Resource
	title   string
	label   fyne.CanvasObject
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

func (p *PaneBase) GetLabel() fyne.CanvasObject {
	return p.label
}

func (p *PaneBase) GetIcon() fyne.Resource {
	return p.icon
}

func (p *PaneBase) SetIcon(icon fyne.Resource) {
	p.icon = icon
	p.label.(*widget.Button).SetIcon(icon)
}

func (p *PaneBase) GetTitle() string {
	return p.title
}

func (p *PaneBase) SetTitle(title string) {
	p.title = title
	p.label.(*widget.Button).SetText(title)
}
