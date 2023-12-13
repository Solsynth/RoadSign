package sign

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type ProcessConfig struct {
	ID       string     `json:"id" yaml:"id"`
	Workdir  string     `json:"workdir" yaml:"workdir"`
	Command  []string   `json:"command" yaml:"command"`
	Prepares [][]string `json:"prepares" yaml:"prepares"`
	Preheat  bool       `json:"preheat" yaml:"preheat"`

	Cmd *exec.Cmd `json:"-"`
}

func (v *ProcessConfig) BootProcess() error {
	if v.Cmd != nil {
		return nil
	}
	if err := v.PrepareProcess(); err != nil {
		return err
	}
	if v.Cmd == nil {
		return v.StartProcess()
	}
	if v.Cmd.Process == nil || v.Cmd.ProcessState == nil {
		return v.StartProcess()
	}
	if v.Cmd.ProcessState.Exited() {
		return v.StartProcess()
	} else if v.Cmd.ProcessState.Exited() {
		return fmt.Errorf("process already dead")
	}
	if v.Cmd.ProcessState.Exited() {
		return fmt.Errorf("cannot start process")
	} else {
		return nil
	}
}

func (v *ProcessConfig) PrepareProcess() error {
	for _, script := range v.Prepares {
		if len(script) <= 0 {
			continue
		}
		cmd := exec.Command(script[0], script[1:]...)
		cmd.Dir = filepath.Join(v.Workdir)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (v *ProcessConfig) StartProcess() error {
	if len(v.Command) <= 0 {
		return fmt.Errorf("you need set the command for %s to enable process manager", v.ID)
	}

	v.Cmd = exec.Command(v.Command[0], v.Command[1:]...)
	v.Cmd.Dir = filepath.Join(v.Workdir)

	return v.Cmd.Start()
}

func (v *ProcessConfig) StopProcess() error {
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

func (v *RoadApp) PreheatProcesses(callback func(total int, success int)) {
	var processes []*ProcessConfig
	for _, site := range v.Sites {
		for _, process := range site.Processes {
			if process.Preheat {
				processes = append(processes, process)
			}
		}
	}

	success := 0
	for _, process := range processes {
		if process.BootProcess() == nil {
			success++
		}
	}

	callback(len(processes), success)
}
