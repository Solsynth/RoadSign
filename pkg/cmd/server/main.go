package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	roadsign "code.smartsheep.studio/goatworks/roadsign/pkg"
	"code.smartsheep.studio/goatworks/roadsign/pkg/hypertext"
	"code.smartsheep.studio/goatworks/roadsign/pkg/sideload"
	"code.smartsheep.studio/goatworks/roadsign/pkg/sign"
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
	viper.SetConfigName("settings")
	viper.SetConfigType("yaml")

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

	// Load & init sign
	if err := sign.ReadInConfig(viper.GetString("paths.configs")); err != nil {
		log.Panic().Err(err).Msg("An error occurred when loading configurations.")
	} else {
		log.Info().Int("count", len(sign.App.Sites)).Msg("All configuration has been loaded.")
	}

	// Preheat processes
	go func() {
		log.Info().Msg("Preheating processes...")
		sign.App.PreheatProcesses(func(total int, success int) {
			log.Info().Int("requested", total).Int("succeed", success).Msgf("Preheat processes completed!")
		})
	}()

	// Init hypertext server
	hypertext.RunServer(
		hypertext.InitServer(),
		viper.GetStringSlice("hypertext.ports"),
		viper.GetStringSlice("hypertext.secured_ports"),
		viper.GetString("hypertext.certificate.pem"),
		viper.GetString("hypertext.certificate.key"),
	)

	// Init sideload server
	hypertext.RunServer(
		sideload.InitSideload(),
		viper.GetStringSlice("hypertext.sideload_ports"),
		viper.GetStringSlice("hypertext.sideload_secured_ports"),
		viper.GetString("hypertext.certificate.sideload_pem"),
		viper.GetString("hypertext.certificate.sideload_key"),
	)

	log.Info().Msgf("RoadSign v%s is started...", roadsign.AppVersion)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msgf("RoadSign v%s is quitting...", roadsign.AppVersion)
}
