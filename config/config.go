package config

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	Debug = "debug"

	MaxImageLongSide = "max-paper-long-side"
	MaxImageShortSide = "max-paper-short-side"
	PixelsPerInch = "pixels-per-inch"
	OutputPrefix = "output-prefix"

	// These measurements discovered by just printing on my printer and measuring the gap
	DefaultMaxImageLongSide = 10.625
	DefaultMaxImageShortSide = 8.125
	DefaultPixelsPerInch = 100
	DefaultOutputPrefix = "img"
)

func InitializeConfig(cmd *cobra.Command) error {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.BindPFlags(cmd.Flags())

	viper.SetDefault(MaxImageLongSide, DefaultMaxImageLongSide)
	viper.SetDefault(MaxImageShortSide, DefaultMaxImageShortSide)
	viper.SetDefault(PixelsPerInch, DefaultPixelsPerInch)
	viper.SetDefault(OutputPrefix, DefaultOutputPrefix)

	return nil
}
