package post_process

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

type PostProcessService interface {
	Apply(ctx context.Context, file string, args ...string) (string, error)
}

type ffmpegPostProcess struct {
	logOutWriter io.Writer
	logErrWriter io.Writer
}

func NewFfmpegPostProcess(logOutWriter, logErrWriter io.Writer) PostProcessService {
	return &ffmpegPostProcess{
		logOutWriter: logOutWriter,
		logErrWriter: logErrWriter,
	}
}

func (s *ffmpegPostProcess) Apply(ctx context.Context, file string, args ...string) (string, error) {
	commandPart := []string{
		"ffmpeg",
		"-i", file,
	}
	commandPart = append(commandPart, args...)

	output := getOutputFile(file)
	commandPart = append(commandPart, output, "-y")
	command := strings.Join(commandPart, " ")

	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = s.logOutWriter
	cmd.Stderr = s.logErrWriter
	log.Printf("[DEBUG] executing command: %s", command)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to execute command: %v", err)
	}
	if _, err := os.Stat(output); err != nil {
		return output, err
	}
	return output, nil
}

func getOutputFile(file string) string {
	dir, fileName := path.Split(file)
	return fmt.Sprintf("%sfpp_%s", dir, fileName)
}
