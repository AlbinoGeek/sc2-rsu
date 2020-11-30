package fynewidget

import "fyne.io/fyne/canvas"

// NewHeader ...
func NewHeader(text string) *canvas.Text {
	return NewText(text, 1.7, true)
}
