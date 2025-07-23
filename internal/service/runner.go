package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nurzhanova2/ci-platform/internal/repository"
	"github.com/nurzhanova2/ci-platform/internal/logger"
	"github.com/nurzhanova2/ci-platform/internal/model"

)

type PipelineService struct {
	GitRepo     *repository.GitRepository
	Docker      *repository.DockerRunner
	Notifier    *NotifierService
	LogWriter   *logger.LogWriter
}

func NewPipelineService(
	git *repository.GitRepository,
	docker *repository.DockerRunner,
	notifier *NotifierService,
	logWriter *logger.LogWriter,
) *PipelineService {
	return &PipelineService{
		GitRepo:   git,
		Docker:    docker,
		Notifier:  notifier,
		LogWriter: logWriter,
	}
}

func (s *PipelineService) RunPipeline(repoURL string, branch string) error {
	log.Printf("Starting pipeline for %s [%s]", repoURL, branch)

	localPath, err := s.GitRepo.Clone(repoURL, branch)
	if err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}
	defer os.RemoveAll(localPath) // cleanup

	pipelinePath := filepath.Join(localPath, ".ci.yaml")
	pipeline, err := model.LoadPipeline(pipelinePath)
	if err != nil {
		return fmt.Errorf("failed to parse pipeline: %w", err)
	}

	for _, step := range pipeline.Steps {
		log.Printf("Running step: %s", step.Name)
		out, err := s.Docker.RunStep(step, localPath)
		s.LogWriter.Write(step.Name, out)

		if err != nil {
			log.Printf("step %s failed: %v", step.Name, err)
			s.Notifier.Send(repoURL, false)
			return fmt.Errorf("step failed: %w", err)
		}
	}

	s.Notifier.Send(repoURL, true)
	return nil
}
