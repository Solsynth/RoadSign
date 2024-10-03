package warden

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

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

	Cmd    *exec.Cmd       `json:"-"`
	Logger strings.Builder `json:"-"`

	Status AppStatus `json:"status"`
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
	v.Cmd.Stdout = &v.Logger
	v.Cmd.Stderr = &v.Logger

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
			} else {
				v.Cmd = nil
			}
			return err
		} else {
			v.Cmd = nil
		}
	}

	v.Status = AppExited
	return nil
}

func (v *AppInstance) Logs() string {
	return v.Logger.String()
}
