package cmd

import (
	"fmt"
	"image/png"
	"os"

	"github.com/nicjohnson145/posterize/config"
	"github.com/nicjohnson145/posterize/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Root() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "posterize PATH_TO_IMAGE",
		Short: "Split image into printable slices",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// So we don't print usage messages on execution errors
			cmd.SilenceUsage = true
			// So we dont double report errors
			cmd.SilenceErrors = true
			return config.InitializeConfig(cmd)
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			log := config.InitLogger()

			slicer := internal.NewSlicer(internal.SlicerConfig{
				Logger: log,
				PixelsPerInch: viper.GetFloat64(config.PixelsPerInch),
				MaxPageLongSide: viper.GetFloat64(config.MaxImageLongSide),
				MaxPageShortSide: viper.GetFloat64(config.MaxImageShortSide),
			})

			content, err := os.ReadFile(args[0])
			if err != nil {
				log.Err(err).Str("path", args[0]).Msg("error reading input file")
				return err
			}

			images, err := slicer.Slice(content)
			if err != nil {
				log.Err(err).Msg("error slicing image")
				return err
			}

			log.Debug().Msg("writing images")
			for i, img := range images {
				name := fmt.Sprintf("%v-%02d.png", viper.GetString(config.OutputPrefix), i)
				f, err := os.Create(name)
				if err != nil {
					log.Err(err).Str("name", name).Msg("error creating output file")
					return err
				}
				defer f.Close()

				err = png.Encode(f, img)
				if err != nil {
					log.Err(err).Msg("error encoding to png")
					return err
				}
			}

			return nil
		},
	}
	rootCmd.PersistentFlags().BoolP(config.Debug, "d", false, "Enable debug logging")
	rootCmd.Flags().Float64P(config.MaxImageLongSide, "l", config.DefaultMaxImageLongSide, "The max value in inches that a sub-image should try to occupy, relative to the long side of a piece of paper")
	rootCmd.Flags().Float64P(config.MaxImageShortSide, "s", config.DefaultMaxImageShortSide, "The max value in inches that a sub-image should try to occupy, relative to the short side of a piece of paper")
	rootCmd.Flags().Float64P(config.PixelsPerInch, "p", config.DefaultPixelsPerInch, "The number of pixels per inch the source image has")
	rootCmd.Flags().StringP(config.OutputPrefix, "o", config.DefaultOutputPrefix, "The prefix each sub image should get when written out")


	rootCmd.AddCommand(
		versionCmd(),
	)

	return rootCmd
}
