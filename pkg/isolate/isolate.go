package isolate

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/joshjms/firefly-executor/config"
)

const BASE_UID int = 60000

// Represents a sandbox.
// The UID of the sandbox is BASE_UID + id.
type Sandbox struct {
	// Path to the root of the sandbox.
	Id string

	Root   string
	Cgroup bool // Whether to use cgroups.
}

func (s *Sandbox) UID() string {
	id, _ := strconv.Atoi(s.Id)

	return fmt.Sprintf("%d", BASE_UID+id)
}

// Options for running a command in the sandbox.
type RunOptions struct {
	// Command to run in the sandbox.
	Args []string
	// Mount points to set up in the sandbox. Nil means no mounts.
	Mount *Mount
	Envs  []Env

	Stdin          string // Path to read stdin from.
	Stdout         string // Path to write stdout to.
	Stderr         string // Path to write stderr to.
	StderrToStdout bool   // Whether to redirect stderr to stdout. This will ignore Stderr.

	Meta string // Path to write the metadata of the execution to.

	// Resource limits
	Processes int   // Maximum number of processes.
	MemLimit  int64 // Memory limit in bytes.
	TimeLimit int64 // Time limit in milliseconds.
}

// Represents a mount point.
type Mount struct {
	// Source of the mount from outside the sandbox.
	Source string
	// Destination of the mount inside the sandbox.
	Destination string

	// Options for the mount.
	Options string
}

type Env struct {
	Key   string
	Value string
}

func NewSandbox(id int64, cg bool) (*Sandbox, error) {
	idStr := fmt.Sprintf("%d", id)

	return &Sandbox{
		Id:     idStr,
		Cgroup: cg,
	}, nil
}

func (s *Sandbox) Init() error {
	cmd := exec.Command(config.GetIsolate())

	if s.Cgroup {
		cmd.Args = append(cmd.Args, "--cg")
	}

	cmd.Args = append(cmd.Args, "--box-id", s.Id)

	cmd.Args = append(cmd.Args, "--init")

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to initialize sandbox: %v\n%s", err, out)
	}

	s.Root = string(out)

	return nil
}

func (s *Sandbox) Cleanup() error {
	cmd := exec.Command(config.GetIsolate())

	if s.Cgroup {
		cmd.Args = append(cmd.Args, "--cg")
	}

	cmd.Args = append(cmd.Args, "--box-id", s.Id)

	cmd.Args = append(cmd.Args, "--cleanup")

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to cleanup sandbox: %v\n%s", err, out)
	}

	return nil
}

func (s *Sandbox) Run(opt *RunOptions) error {
	cmd := exec.Command(config.GetIsolate())

	if s.Cgroup {
		cmd.Args = append(cmd.Args, "--cg")
	}

	cmd.Args = append(cmd.Args, "--box-id", s.Id)

	if opt.Mount != nil {
		if opt.Mount.Options != "" {
			cmd.Args = append(cmd.Args, "--dir", fmt.Sprintf("%s=%s:%s",
				opt.Mount.Destination,
				opt.Mount.Source,
				opt.Mount.Options,
			))
		} else {
			cmd.Args = append(cmd.Args, "--dir", fmt.Sprintf("%s=%s",
				opt.Mount.Destination,
				opt.Mount.Source,
			))
		}
	}

	for _, env := range opt.Envs {
		cmd.Args = append(cmd.Args, "--env", fmt.Sprintf("%s=%s", env.Key, env.Value))
	}

	if opt.Stdin != "" {
		cmd.Args = append(cmd.Args, "--stdin", opt.Stdin)
	}

	if opt.Stdout != "" {
		cmd.Args = append(cmd.Args, "--stdout", opt.Stdout)
	}

	if opt.StderrToStdout {
		cmd.Args = append(cmd.Args, "--stderr-to-stdout")
	} else if opt.Stderr != "" {
		cmd.Args = append(cmd.Args, "--stderr", opt.Stderr)
	}

	if opt.Meta != "" {
		cmd.Args = append(cmd.Args, "-M", opt.Meta)
	}

	if opt.Processes > 0 {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--processes=%d", opt.Processes))
	}

	if opt.MemLimit > 0 {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--mem=%d", opt.MemLimit))
	}

	if opt.TimeLimit > 0 {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--time=%f", float64(opt.TimeLimit/1000)))
	}

	cmd.Args = append(cmd.Args, "--run")

	cmd.Args = append(cmd.Args, "--")
	cmd.Args = append(cmd.Args, opt.Args...)

	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to run command in sandbox: %v\n%s", err, out)
	}

	return nil
}
