package python

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/joshjms/firefly-executor/pkg/controller"
	"github.com/joshjms/firefly-executor/pkg/handlers/submit/utils"
	"github.com/joshjms/firefly-executor/pkg/isolate"
)

func Submit(code, stdin string, memLimit, timeLimit int64) (stdout, stderr string, time, timeWall, maxRss, cswVoluntary, cswForced, cgMem int64, exitCode int, err error) {
	jobId := uuid.New().String()

	os.Mkdir(filepath.Join("mounts", jobId), 0755)

	// defer os.RemoveAll(filepath.Join("mounts", jobId))
	// defer os.RemoveAll(filepath.Join("metadata", jobId))

	boxId := controller.GetEmptyBox()
	if boxId == -1 {
		return "", "", 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("no available box")
	}
	defer controller.ReleaseBox(boxId)

	sandbox, err := isolate.NewSandbox(boxId, true)
	if err != nil {
		return "", "", 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("failed to create sandbox: %v", err)
	}

	if err := sandbox.Init(); err != nil {
		return "", "", 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("failed to initialize sandbox: %v", err)
	}

	if err := exec.Command("setfacl", "-R", "-m", fmt.Sprintf("u:%s:rwx", sandbox.UID()), filepath.Join("mounts", jobId)).Run(); err != nil {
		return "", "", 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("failed to setfacl: %v", err)
	}

	os.Create(filepath.Join("mounts", jobId, "a.py"))
	os.WriteFile(filepath.Join("mounts", jobId, "a.py"), []byte(code), 0644)
	os.Create(filepath.Join("mounts", jobId, "stdin"))
	os.WriteFile(filepath.Join("mounts", jobId, "stdin"), []byte(stdin), 0644)

	currentDir, err := os.Getwd()
	if err != nil {
		return "", "", 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("failed to get current working directory: %v", err)
	}

	runErr := sandbox.Run(&isolate.RunOptions{
		Envs: []isolate.Env{
			{
				Key:   "PYTHONUNBUFFERED",
				Value: "1",
			},
			{
				Key:   "PATH",
				Value: "/bin",
			},
		},
		Args: []string{"/usr/bin/python3", "a.py"},
		Mount: &isolate.Mount{
			Source:      filepath.Join(currentDir, "mounts", jobId),
			Destination: "/box",
			Options:     "rw",
		},
		Stdin:     "stdin",
		Stdout:    "stdout",
		Stderr:    "stderr",
		Meta:      filepath.Join("metadata", jobId),
		Processes: 1,
		MemLimit:  memLimit,
		TimeLimit: timeLimit,
	})

	if err := sandbox.Cleanup(); err != nil {
		return "", "", 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("failed to cleanup sandbox: %v", err)
	}

	stdout, err = utils.ReadStdout(jobId)
	if err != nil {
		return "", "", 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("failed to read stdout: %v", err)
	}

	stderr, err = utils.ReadStderr(jobId)
	if err != nil {
		return "", "", 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("failed to read stderr: %v", err)
	}

	time, timeWall, maxRss, cswVoluntary, cswForced, cgMem, exitCode, err = utils.ReadMetadata(jobId)
	if err != nil {
		return "", "", 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("failed to read metadata: %v", err)
	}

	return stdout, stderr, time, timeWall, maxRss, cswVoluntary, cswForced, cgMem, exitCode, runErr
}
