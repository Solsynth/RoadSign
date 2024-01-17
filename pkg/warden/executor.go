package warden

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/samber/lo"
)

type AppStatus = int8

const (
	AppCreated = AppStatus(iota)
	AppStarting
	AppStarted
	AppExited
	AppFailure
)

type WardenInstance struct {
	Manifest WardenApplication `json:"manifest"`

	Cmd    *exec.Cmd       `json:"-"`
	Logger strings.Builder `json:"-"`

	Status AppStatus `json:"status"`
}

func (v *WardenInstance) Wake() error {
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
	} else if v.Cmd.ProcessState.Exited() {
		return fmt.Errorf("process already dead")
	}
	if v.Cmd.ProcessState.Exited() {
		return fmt.Errorf("cannot start process")
	} else {
		return nil
	}
}

func (v *WardenInstance) Start() error {
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
			if v.Cmd.Process == nil || v.Cmd.ProcessState == nil {
				v.Status = AppStarting
			} else if !v.Cmd.ProcessState.Exited() {
				v.Status = AppStarted
			} else {
				v.Status = lo.Ternary(v.Cmd.ProcessState.Success(), AppExited, AppFailure)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return v.Cmd.Start()
}

func (v *WardenInstance) Stop() error {
	if v.Cmd != nil && v.Cmd.Process != nil {
		if err := v.Cmd.Process.Signal(os.Interrupt); err != nil {
			v.Cmd.Process.Kill()
			return err
		} else {
			v.Cmd = nil
		}
	}

	return nil
}

func (v *WardenInstance) Logs() string {
	return v.Logger.String()
}
