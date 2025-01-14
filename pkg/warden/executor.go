package warden

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/samber/lo"
)

var InstancePool []*AppInstance

func GetFromPool(id string) *AppInstance {
	val, ok := lo.Find(InstancePool, func(item *AppInstance) bool {
		return item.Manifest.ID == id
	})
	return lo.Ternary(ok, val, nil)
}

func StartPool() []error {
	var errors []error
	for _, instance := range InstancePool {
		if err := instance.Wake(); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

type AppStatus = int8

const (
	AppCreated = AppStatus(iota)
	AppStarting
	AppStarted
	AppExited
	AppFailure
)

type AppInstance struct {
	Manifest Application `json:"manifest"`
	Status   AppStatus   `json:"status"`

	Cmd *exec.Cmd `json:"-"`

	LogPath string             `json:"-"`
	Logger  *lumberjack.Logger `json:"-"`
}

func (v *AppInstance) Wake() error {
	if v.Cmd != nil {
		return nil
	}
	if v.Cmd == nil {
		return v.Start()
	}
	if v.Cmd.Process == nil || v.Cmd.ProcessState == nil {
		return v.Start()
	}
	if v.Cmd.ProcessState.Exited() {
		return v.Start()
	}
	return nil
}

func (v *AppInstance) Start() error {
	manifest := v.Manifest

	if len(manifest.Command) <= 0 {
		return fmt.Errorf("you need set the command for %s to enable process manager", manifest.ID)
	}

	v.Cmd = exec.Command(manifest.Command[0], manifest.Command[1:]...)
	v.Cmd.Dir = filepath.Join(manifest.Workdir)
	v.Cmd.Env = append(v.Cmd.Env, manifest.Environment...)

	logBasePath := viper.GetString("logging.warden_apps")
	logPath := filepath.Join(logBasePath, fmt.Sprintf("%s.log", manifest.ID))

	v.LogPath = logPath
	v.Logger = &lumberjack.Logger{
		Filename:   v.LogPath,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}

	v.Cmd.Stdout = v.Logger
	v.Cmd.Stderr = v.Logger

	// Monitor
	go func() {
		for {
			if v.Cmd != nil && v.Cmd.Process == nil {
				v.Status = AppStarting
			} else if v.Cmd != nil && v.Cmd.ProcessState == nil {
				v.Status = AppStarted
			} else {
				v.Status = AppFailure
				v.Cmd = nil
				return
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	return v.Cmd.Start()
}

func (v *AppInstance) Stop() error {
	if v.Cmd != nil && v.Cmd.Process != nil {
		if err := v.Cmd.Process.Signal(syscall.SIGTERM); err != nil {
			log.Warn().Int("pid", v.Cmd.Process.Pid).Err(err).Msgf("Failed to send SIGTERM to process...")
			if err = v.Cmd.Process.Kill(); err != nil {
				log.Error().Int("pid", v.Cmd.Process.Pid).Err(err).Msgf("Failed to kill process...")
				return err
			}
		}
	}

	// We need to wait for the process to exit
	// The wait syscall will read the exit status of the process
	// So that we don't produce defunct processes
	// Refer to https://stackoverflow.com/questions/46293435/golang-exec-command-cause-a-lot-of-defunct-processes
	_ = v.Cmd.Wait()

	v.Cmd = nil
	v.Status = AppExited
	v.Logger.Close()
	return nil
}

func (v *AppInstance) Logs() string {
	file, err := os.Open(v.LogPath)
	if err != nil {
		return ""
	}
	raw, _ := io.ReadAll(file)
	return string(raw)
}
