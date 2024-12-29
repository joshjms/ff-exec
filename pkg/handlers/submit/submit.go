package submit

import (
	"strings"

	"github.com/joshjms/firefly-executor/pkg/handlers/submit/c"
	"github.com/joshjms/firefly-executor/pkg/handlers/submit/cpp"
	"github.com/joshjms/firefly-executor/pkg/handlers/submit/python"
)

type Job struct {
	Code      string `json:"code"`
	Language  string `json:"language"`
	MemLimit  int64  `json:"mem_limit"`
	TimeLimit int64  `json:"time_limit"`
	Stdin     string `json:"stdin"`
}

type JobResponse struct {
	Stdout  string  `json:"stdout"`
	Stderr  string  `json:"stderr"`
	Details Details `json:"details"`
	Error   string  `json:"error"`
}

type Details struct {
	Time         int64 `json:"time"`
	TimeWall     int64 `json:"time_wall"`
	MaxRss       int64 `json:"max_rss"`
	CswVoluntary int64 `json:"csw_voluntary"`
	CswForced    int64 `json:"csw_forced"`
	CgMem        int64 `json:"cg_mem"`
	ExitCode     int   `json:"exit_code"`
}

func Submit(j *Job) *JobResponse {
	var stdout, stderr string
	var time, timeWall, maxRss, cswVoluntary, cswForced, cgMem int64
	var exitCode int

	var err error

	switch j.Language {
	case "c":
		stdout, stderr, time, timeWall, maxRss, cswVoluntary, cswForced, cgMem, exitCode, err = c.Submit(j.Code, j.Stdin, j.MemLimit, j.TimeLimit)
	case "cpp":
		stdout, stderr, time, timeWall, maxRss, cswVoluntary, cswForced, cgMem, exitCode, err = cpp.Submit(j.Code, j.Stdin, j.MemLimit, j.TimeLimit)
	case "python":
		stdout, stderr, time, timeWall, maxRss, cswVoluntary, cswForced, cgMem, exitCode, err = python.Submit(j.Code, j.Stdin, j.MemLimit, j.TimeLimit)
	default:
		return &JobResponse{
			Error: "unsupported language",
		}
	}

	if err != nil {
		errMsg := strings.Split(strings.TrimSpace(err.Error()), "\n")
		lastErr := errMsg[len(errMsg)-1]

		return &JobResponse{
			Error: lastErr,
		}
	}

	return &JobResponse{
		Stdout: stdout,
		Stderr: stderr,
		Details: Details{
			Time:         time,
			TimeWall:     timeWall,
			MaxRss:       maxRss,
			CswVoluntary: cswVoluntary,
			CswForced:    cswForced,
			CgMem:        cgMem,
			ExitCode:     exitCode,
		},
	}
}
