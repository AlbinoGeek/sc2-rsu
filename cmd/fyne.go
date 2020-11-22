package cmd

import (
	"fmt"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"

	"github.com/spf13/viper"
)

func getConfigColor(key string, def color.Color) color.Color {
	slice := viper.Get(fmt.Sprintf("%s.r", key))
	if slice == nil {
		return def
	}

	clr := color.NRGBA{
		uint8(viper.GetUint(fmt.Sprintf("%s.r", key))),
		uint8(viper.GetUint(fmt.Sprintf("%s.g", key))),
		uint8(viper.GetUint(fmt.Sprintf("%s.b", key))),
		uint8(viper.GetUint(fmt.Sprintf("%s.a", key))),
	}

	if clr.A == 0 {
		return def
	}

	return def
}

func labelWithWrapping(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	return label
}

func newText(text string, scale float32, bold bool) *canvas.Text {
	return &canvas.Text{
		Color:     GUI.theme.TextColor(),
		Text:      text,
		TextSize:  int(float32(GUI.theme.TextSize()) * scale),
		TextStyle: fyne.TextStyle{Bold: bold},
	}
}
