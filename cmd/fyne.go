package cmd

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
)

func labelWithWrapping(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	return label
}

func newText(text string, scale float32, bold bool) *canvas.Text {
	return &canvas.Text{
		Color:     GUI.Theme.TextColor(),
		Text:      text,
		TextSize:  int(float32(GUI.Theme.TextSize()) * scale),
		TextStyle: fyne.TextStyle{Bold: bold},
	}
}
