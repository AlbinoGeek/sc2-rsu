package cmd

import (
	"image/color"

	"fyne.io/fyne"
	"github.com/spf13/viper"
)

type guiTheme struct {
	Base fyne.Theme
}

// colors //

func (guiTheme) BackgroundColor() color.Color {
	return getConfigColor("theme.color.background", color.NRGBA{0x20, 0x20, 0x22, 0xff})
}

func (guiTheme) ButtonColor() color.Color {
	return getConfigColor("theme.color.button", color.Transparent)
}

func (guiTheme) DisabledButtonColor() color.Color {
	return getConfigColor("theme.color.disabledButton", color.NRGBA{0x16, 0x16, 0x18, 0xff})
}

func (t guiTheme) DisabledIconColor() color.Color {
	return getConfigColor("theme.color.disabled", t.Base.DisabledIconColor())
}

func (t guiTheme) DisabledTextColor() color.Color {
	return getConfigColor("theme.color.disabled", t.Base.DisabledTextColor())
}

func (t guiTheme) FocusColor() color.Color {
	return getConfigColor("theme.color.focus", t.Base.FocusColor())
}

func (t guiTheme) HoverColor() color.Color {
	return getConfigColor("theme.color.hover", t.Base.HoverColor())
}

func (t guiTheme) HyperlinkColor() color.Color {
	return getConfigColor("theme.color.primary", t.Base.HyperlinkColor())
}

func (t guiTheme) IconColor() color.Color {
	return getConfigColor("theme.color.text", t.Base.IconColor())
}

func (t guiTheme) PlaceHolderColor() color.Color {
	return getConfigColor("theme.color.placeholder", t.Base.PlaceHolderColor())
}

func (t guiTheme) PrimaryColor() color.Color {
	return getConfigColor("theme.color.primary", t.Base.PrimaryColor())
}

func (t guiTheme) ScrollBarColor() color.Color {
	return getConfigColor("theme.color.scrollBar", t.Base.ScrollBarColor())
}

func (t guiTheme) ShadowColor() color.Color {
	return getConfigColor("theme.color.shadow", t.Base.ShadowColor())
}

func (t guiTheme) TextColor() color.Color {
	return getConfigColor("theme.color.text", t.Base.TextColor())
}

// integers //

func (guiTheme) IconInlineSize() int {
	return viper.GetInt("theme.iconInlineSize")
}

func (guiTheme) Padding() int {
	return viper.GetInt("theme.padding")
}

func (guiTheme) ScrollBarSize() int {
	return viper.GetInt("theme.scrollBarSize")
}

func (guiTheme) ScrollBarSmallSize() int {
	return viper.GetInt("theme.scrollBarSmallSize")
}

func (guiTheme) TextSize() int {
	return viper.GetInt("theme.textSize")
}

// fonts //

func (t guiTheme) TextFont() fyne.Resource {
	return t.Base.TextFont()
}

func (t guiTheme) TextBoldFont() fyne.Resource {
	return t.Base.TextBoldFont()
}

func (t guiTheme) TextBoldItalicFont() fyne.Resource {
	return t.Base.TextBoldItalicFont()
}

func (t guiTheme) TextItalicFont() fyne.Resource {
	return t.Base.TextItalicFont()
}

func (t guiTheme) TextMonospaceFont() fyne.Resource {
	return t.Base.TextMonospaceFont()
}
