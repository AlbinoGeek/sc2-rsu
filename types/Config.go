package types

import (
	"fmt"
	"image/color"

	"github.com/spf13/viper"
)

// GetColor returns an NRGBA key with the defined prefix, if set -- otherwise
// the color passed in as def is returned instead.
func GetColor(key string, def color.Color) color.Color {
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
