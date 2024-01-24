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
	} else if v.Cmd.ProcessState.Exited() {
		return fmt.Errorf("process already dead")
	}
	if v.Cmd.ProcessState.Exited() {
		return fmt.Errorf("cannot start process")
	} else {
		return nil
	}
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

func (v *AppInstance) Stop() error {
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

func (v *AppInstance) Logs() string {
	return v.Logger.String()
}
