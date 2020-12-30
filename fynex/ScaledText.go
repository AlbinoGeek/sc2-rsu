package fynex

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

// TextSize refers to a Typgraphical Sizing defined by Material Design
type TextSize uint8

const (
	// TextSizeBody1 is 16dp and represents regular body type
	TextSizeBody1 TextSize = iota

	// TextSizeBody2 is 14dp and represents condensed body type
	TextSizeBody2

	// TextSizeSubtitle1 is 16dp and represents a medium subtitle
	TextSizeSubtitle1

	// TextSizeSubtitle2 is 14dp and represents a light subtitle
	TextSizeSubtitle2

	// TextSizeHeading1 is 96dp and represents hero type
	TextSizeHeading1

	// TextSizeHeading2 is 60dp and represents a primary heading
	TextSizeHeading2

	// TextSizeHeading3 is 48dp and represents a secondary heading
	TextSizeHeading3

	// TextSizeHeading4 is 34dp
	TextSizeHeading4

	// TextSizeHeading5 is 24dp
	TextSizeHeading5

	// TextSizeHeading6 is 20dp and used by component Titles
	TextSizeHeading6
)

var styleSize = map[TextSize]float32{
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
	return newText(text, styleSize[level], false)
}

// Style is a wrapped fyne.TextStyle with a TextSize
type Style struct {
	Size TextSize
	fyne.TextStyle
}

// StyleHeading5 is the Material Design standard Heading5 text style
func StyleHeading5() Style {
	return Style{
		Size: TextSizeHeading5,
		TextStyle: fyne.TextStyle{
			Bold: true,
		},
	}
}

// NewTextWithStyle returns a canvas.Text element with a given alignment
// and Material Design Style applied to it. Mimics widget.NewLabelWithStyle
func NewTextWithStyle(text string, alignment fyne.TextAlign, style Style) *canvas.Text {
	t := newText(text, styleSize[style.Size], false)
	t.Alignment = alignment
	t.TextStyle = style.TextStyle
	return t
}

// ! Material Design base font size is 16 but fyne is 14 ...
// ! until fyne changes this, we're going to scale them up.
func newText(text string, scale float32, bold bool) *canvas.Text {
	t := canvas.NewText(text, theme.TextColor())

	if t.TextSize == 14 {
		t.TextSize = 1.14 * t.TextSize * scale
	} else {
		t.TextSize = t.TextSize * scale
	}

	t.TextStyle.Bold = bold

	return t
}
