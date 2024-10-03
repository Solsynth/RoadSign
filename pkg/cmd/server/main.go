package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	roadsign "git.solsynth.dev/goatworks/roadsign/pkg"
	"git.solsynth.dev/goatworks/roadsign/pkg/hypertext"
	"git.solsynth.dev/goatworks/roadsign/pkg/navi"
	"git.solsynth.dev/goatworks/roadsign/pkg/sideload"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func main() {
	// Configure settings
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("/")
	viper.SetConfigName("settings")
	viper.SetConfigType("toml")

	// Load settings
	if err := viper.ReadInConfig(); err != nil {
		log.Panic().Err(err).Msg("An error occurred when loading settings.")
	}

	// Present settings
	if len(viper.GetString("security.credential")) <= 0 {
		credential := strings.ReplaceAll(uuid.NewString(), "-", "")
		viper.Set("security.credential", credential)
		_ = viper.WriteConfig()

		log.Warn().Msg("There isn't any api credential configured in settings.yml, auto generated a credential for api accessing.")
		log.Warn().Msgf("RoadSign auto generated api credential is %s", credential)
	}

	// Load & init navigator
	if err := navi.ReadInConfig(viper.GetString("paths.configs")); err != nil {
		log.Panic().Err(err).Msg("An error occurred when loading configurations.")
	} else {
		log.Info().Int("count", len(navi.R.Regions)).Msg("All configuration has been loaded.")
	}

	// Init warden
	navi.InitializeWarden(navi.R.Regions)

	// Init hypertext server
	hypertext.RunServer(
		hypertext.InitServer(),
		viper.GetStringSlice("hypertext.ports"),
		viper.GetStringSlice("hypertext.secured_ports"),
	)

	// Init sideload server
	hypertext.RunServer(
		sideload.InitSideload(),
		viper.GetStringSlice("sideload.ports"),
		[]string{},
	)

	log.Info().Msgf("RoadSign v%s is started...", roadsign.AppVersion)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msgf("RoadSign v%s is quitting...", roadsign.AppVersion)
}
