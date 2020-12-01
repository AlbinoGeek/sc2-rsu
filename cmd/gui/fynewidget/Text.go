package fynewidget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

// NewText ...
func NewText(text string, scale float32, bold bool) *canvas.Text {
	return &canvas.Text{
		Color:     theme.TextColor(),
		Text:      text,
		TextSize:  int(float32(theme.TextSize()) * scale),
		TextStyle: fyne.TextStyle{Bold: bold},
	}
}
