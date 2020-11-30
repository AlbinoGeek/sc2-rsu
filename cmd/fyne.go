package cmd

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func labelWithWrapping(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord

	return label
}
