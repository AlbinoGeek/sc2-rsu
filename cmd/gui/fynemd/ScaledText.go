package fynemd

import (
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

// ! Material Design base font size is 16 but fyne is 14 ...
// ! until fyne changes this, we're going scale them up.

type TextSize uint8

const (
	TextSizeBody1 TextSize = iota
	TextSizeBody2
	TextSizeSubtitle1
	TextSizeSubtitle2
	TextSizeHeading1
	TextSizeHeading2
	TextSizeHeading3
	TextSizeHeading4
	TextSizeHeading5
	TextSizeHeading6
)

var styleSize = map[TextSize]float64{
	TextSizeBody1:     1,
	TextSizeBody2:     .875,
	TextSizeSubtitle1: 1,
	TextSizeSubtitle2: .875,
	TextSizeHeading1:  6,
	TextSizeHeading2:  3.75,
	TextSizeHeading3:  3,
	TextSizeHeading4:  2.125,
	TextSizeHeading5:  1.5,
	TextSizeHeading6:  1.25,
}

// NewScaledText returns a canvas.Text element with a given Material
// Design type scale applied to it.
func NewScaledText(level TextSize, text string) *canvas.Text {
	return newText(text, 1.14*styleSize[level], false)
}

func newText(text string, scale float64, bold bool) *canvas.Text {
	t := canvas.NewText(text, theme.TextColor())
	t.TextSize = int(float64(t.TextSize) * scale)
	t.TextStyle.Bold = bold
	return t
}
