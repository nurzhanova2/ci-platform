package repository

import (
    "bytes"
    "context"
    "errors"
    "fmt"
    "log"
    "os/exec"
    "strings"

    "github.com/nurzhanova2/ci-platform/internal/model"
)

type DockerRunner struct{}

func NewDockerRunner() *DockerRunner {
    return &DockerRunner{}
}

func (d *DockerRunner) RunStep(step model.Step, workDir string) (string, error) {

    if step.Image == "" {
        return "", errors.New("step.Image is required")
    }
    if len(step.Commands) == 0 {
        return "", fmt.Errorf("no commands specified for step %s", step.Name)
    }

    script := strings.Join(step.Commands, " && ")
    containerName := "ci-step-" + step.Name

    cmdArgs := []string{
        "run", "--rm",
        "--name", containerName,
        "-v", fmt.Sprintf("%s:/workspace", workDir),
        "-w", "/workspace",
        step.Image,
        "sh", "-c", script,
    }

    log.Printf("Running Docker command: docker %s\n", strings.Join(cmdArgs, " "))

    cmd := exec.CommandContext(context.Background(), "docker", cmdArgs...)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out

    err := cmd.Run()
    if err != nil {
        return out.String(), fmt.Errorf("docker run failed: %v\nOutput:\n%s", err, out.String())
    }

    return out.String(), nil
}