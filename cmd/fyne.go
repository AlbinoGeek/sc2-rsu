package cmd

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func labelWithWrapping(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	return label
}

func newText(text string, size int, bold bool) *canvas.Text {
	return &canvas.Text{
		Color:     theme.TextColor(),
		Text:      text,
		TextSize:  size,
		TextStyle: fyne.TextStyle{Bold: bold},
	}
}
