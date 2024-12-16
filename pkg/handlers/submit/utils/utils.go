package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ReadStdout(jobId string) (string, error) {
	stdoutBytes, err := os.ReadFile(filepath.Join("mounts", jobId, "stdout"))
	if err != nil {
		return "", fmt.Errorf("failed to read stdout: %v", err)
	}

	return string(stdoutBytes), nil
}

func ReadStderr(jobId string) (string, error) {
	stderrBytes, err := os.ReadFile(filepath.Join("mounts", jobId, "stderr"))
	if err != nil {
		return "", fmt.Errorf("failed to read stderr: %v", err)
	}

	return string(stderrBytes), nil
}

func ReadMetadata(jobId string) (time, timeWall, maxRss, cswVoluntary, cswForced, cgMem int64, exitCode int, err error) {
	metadataBytes, err := os.ReadFile(filepath.Join("metadata", jobId))
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("failed to read metadata: %v", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(metadataBytes))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		keyValue := strings.Split(line, ":")
		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])

		switch key {
		case "time":
			timeFloat, _ := strconv.ParseFloat(value, 64)
			time = int64(timeFloat * 1000)
		case "time-wall":
			timeWallFloat, _ := strconv.ParseFloat(value, 64)
			timeWall = int64(timeWallFloat * 1000)
		case "max-rss":
			maxRss, _ = strconv.ParseInt(value, 10, 64)
		case "csw-voluntary":
			cswVoluntary, _ = strconv.ParseInt(value, 10, 64)
		case "csw-forced":
			cswForced, _ = strconv.ParseInt(value, 10, 64)
		case "cg-mem":
			cgMem, _ = strconv.ParseInt(value, 10, 64)
		case "exitcode":
			exitCode, _ = strconv.Atoi(value)
		}
	}

	return time, timeWall, maxRss, cswVoluntary, cswForced, cgMem, exitCode, nil
}
