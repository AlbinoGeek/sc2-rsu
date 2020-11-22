package gui

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"

	"github.com/spf13/viper"

	"github.com/AlbinoGeek/sc2-rsu/types"
)

// DarkerTheme is a darker variant of fyne.DarkTheme()
func DarkerTheme() fyne.Theme {
	return &darkerTheme{
		base: theme.DarkTheme(),
	}
}

type darkerTheme struct {
	base fyne.Theme
}

// colors //

func (darkerTheme) BackgroundColor() color.Color {
	return types.GetColor("theme.color.background", color.NRGBA{0x20, 0x20, 0x22, 0xff})
}

func (darkerTheme) ButtonColor() color.Color {
	return types.GetColor("theme.color.button", color.Transparent)
}

func (darkerTheme) DisabledButtonColor() color.Color {
	return types.GetColor("theme.color.disabledButton", color.NRGBA{0x16, 0x16, 0x18, 0xff})
}

func (theme darkerTheme) DisabledIconColor() color.Color {
	return types.GetColor("theme.color.disabled", theme.base.DisabledIconColor())
}

func (theme darkerTheme) DisabledTextColor() color.Color {
	return types.GetColor("theme.color.disabled", theme.base.DisabledTextColor())
}

func (theme darkerTheme) FocusColor() color.Color {
	return types.GetColor("theme.color.focus", theme.base.FocusColor())
}

func (theme darkerTheme) HoverColor() color.Color {
	return types.GetColor("theme.color.hover", theme.base.HoverColor())
}

func (theme darkerTheme) HyperlinkColor() color.Color {
	return types.GetColor("theme.color.primary", theme.base.HyperlinkColor())
}

func (theme darkerTheme) IconColor() color.Color {
	return types.GetColor("theme.color.text", theme.base.IconColor())
}

func (theme darkerTheme) PlaceHolderColor() color.Color {
	return types.GetColor("theme.color.placeholder", theme.base.PlaceHolderColor())
}

func (theme darkerTheme) PrimaryColor() color.Color {
	return types.GetColor("theme.color.primary", theme.base.PrimaryColor())
}

func (theme darkerTheme) ScrollBarColor() color.Color {
	return types.GetColor("theme.color.scrollBar", theme.base.ScrollBarColor())
}

func (theme darkerTheme) ShadowColor() color.Color {
	return types.GetColor("theme.color.shadow", theme.base.ShadowColor())
}

func (theme darkerTheme) TextColor() color.Color {
	return types.GetColor("theme.color.text", theme.base.TextColor())
}

// integers //

func (darkerTheme) IconInlineSize() int {
	return viper.GetInt("theme.iconInlineSize")
}

func (darkerTheme) Padding() int {
	return viper.GetInt("theme.padding")
}

func (darkerTheme) ScrollBarSize() int {
	return viper.GetInt("theme.scrollBarSize")
}

func (darkerTheme) ScrollBarSmallSize() int {
	return viper.GetInt("theme.scrollBarSmallSize")
}

func (darkerTheme) TextSize() int {
	return viper.GetInt("theme.textSize")
}

// fonts //

func (theme darkerTheme) TextFont() fyne.Resource {
	return theme.base.TextFont()
}

func (theme darkerTheme) TextBoldFont() fyne.Resource {
	return theme.base.TextBoldFont()
}

func (theme darkerTheme) TextBoldItalicFont() fyne.Resource {
	return theme.base.TextBoldItalicFont()
}

func (theme darkerTheme) TextItalicFont() fyne.Resource {
	return theme.base.TextItalicFont()
}

func (theme darkerTheme) TextMonospaceFont() fyne.Resource {
	return theme.base.TextMonospaceFont()
}
