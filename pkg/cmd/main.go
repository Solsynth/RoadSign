package main

import (
	roadsign "code.smartsheep.studio/goatworks/roadsign/pkg"
	"code.smartsheep.studio/goatworks/roadsign/pkg/administration"
	"code.smartsheep.studio/goatworks/roadsign/pkg/configurator"
	"code.smartsheep.studio/goatworks/roadsign/pkg/hypertext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func main() {
	// Configure settings
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.SetConfigName("settings")
	viper.SetConfigType("yaml")

	// Load settings
	if err := viper.ReadInConfig(); err != nil {
		log.Panic().Err(err).Msg("An error occurred when loading settings.")
	}

	// Load configurations
	if err := configurator.ReadInConfig(viper.GetString("paths.configs")); err != nil {
		log.Panic().Err(err).Msg("An error occurred when loading configurations.")
	} else {
		log.Debug().Any("sites", configurator.C).Msg("All configuration has been loaded.")
	}

	// Init hypertext server
	hypertext.RunServer(
		hypertext.InitServer(),
		viper.GetStringSlice("hypertext.ports"),
		viper.GetStringSlice("hypertext.secured_ports"),
		viper.GetString("hypertext.certificate.pem"),
		viper.GetString("hypertext.certificate.key"),
	)

	// Init administration server
	hypertext.RunServer(
		administration.InitAdministration(),
		viper.GetStringSlice("hypertext.administration_ports"),
		viper.GetStringSlice("hypertext.administration_secured_ports"),
		viper.GetString("hypertext.certificate.administration_pem"),
		viper.GetString("hypertext.certificate.administration_key"),
	)

	log.Info().Msgf("RoadSign v%s is started...", roadsign.AppVersion)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msgf("RoadSign v%s is quitting...", roadsign.AppVersion)
}
