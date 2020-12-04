package fynemd

import (
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

// NewText ...
func NewText(text string, scale float32, bold bool) *canvas.Text {
	t := canvas.NewText(text, theme.TextColor())
	t.TextSize = int(float32(t.TextSize) * scale)
	t.TextStyle.Bold = bold
	return t
}
