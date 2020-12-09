package fynemd

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"github.com/AlbinoGeek/sc2-rsu/cmd/gui"
)

// Pane ...
//
// Implements: NavigationItem
type Pane interface {
	GetWindow() gui.Window

	GetContent() fyne.CanvasObject
	SetContent(fyne.CanvasObject)
	GetLabel() fyne.CanvasObject
	GetIcon() fyne.Resource
	SetIcon(fyne.Resource)
	GetTitle() string
	SetTitle(string)
}

// NewPane returns a pane to be used with navigation, specifying a title
func NewPane(title string, window gui.Window) Pane {
	return NewPaneWithIcon(title, nil, window)
}

// NewPaneWithIcon returns a pane to be used with navigation, specifying
// title and icon
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

// PaneBase implements Pane
type PaneBase struct {
	content fyne.CanvasObject
	icon    fyne.Resource
	title   string
	label   fyne.CanvasObject
	window  gui.Window
}

// GetWindow returns the parent window this pane is shown in
func (p *PaneBase) GetWindow() gui.Window {
	return p.window
}

// GetContent returns the content of this pane
func (p *PaneBase) GetContent() fyne.CanvasObject {
	return p.content
}

// SetContent changes the element to be shown when this pane is selected
func (p *PaneBase) SetContent(content fyne.CanvasObject) {
	p.content = content
}

// GetLabel returns the assemebled label element for this pane
func (p *PaneBase) GetLabel() fyne.CanvasObject {
	return p.label
}

// GetIcon returns the icon shown in labels associated with this pane
func (p *PaneBase) GetIcon() fyne.Resource {
	return p.icon
}

// SetIcon changes the icon shown in labels associated with this pane
func (p *PaneBase) SetIcon(icon fyne.Resource) {
	p.icon = icon
	p.label.(*widget.Button).SetIcon(icon)
}

// GetTitle returns the text shown in labels associated with this pane
func (p *PaneBase) GetTitle() string {
	return p.title
}

// SetTitle changes the text shown in labels associated with this pane
func (p *PaneBase) SetTitle(title string) {
	p.title = title
	p.label.(*widget.Button).SetText(title)
}
