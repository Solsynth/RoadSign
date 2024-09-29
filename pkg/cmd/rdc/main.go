package main

import (
	"os"

	roadsign "git.solsynth.dev/goatworks/roadsign/pkg"
	"git.solsynth.dev/goatworks/roadsign/pkg/cmd/rdc/conn"
	"git.solsynth.dev/goatworks/roadsign/pkg/cmd/rdc/deploy"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func main() {
	// Configure settings
	viper.AddConfigPath("$HOME")
	viper.SetConfigName(".roadsignrc")
	viper.SetConfigType("toml")

	// Load settings
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SafeWriteConfig()
			viper.ReadInConfig()
		} else {
			log.Panic().Err(err).Msg("An error occurred when loading settings.")
		}
	}

	// Configure CLI
	app := &cli.App{
		Name:     "RoadSign CLI",
		Version:  roadsign.AppVersion,
		Suggest:  true,
		Commands: append(append([]*cli.Command{}, conn.CliCommands...), deploy.DeployCommands...),
	}

	// Run CLI
	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when running cli.")
	}
}
