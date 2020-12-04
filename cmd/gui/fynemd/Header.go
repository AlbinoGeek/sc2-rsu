package fynemd

import "fyne.io/fyne/canvas"

// NewHeader ...
func NewHeader(text string) *canvas.Text {
	return NewText(text, 1.43, true) // approx 20dp
}
