package repository

import (
	"context"

	"github.com/nurzhanova2/ci-platform/internal/model"
)

type JobRepository interface {
	SavePipelineRun(ctx context.Context, run model.PipelineRun) error
	GetLastRuns(ctx context.Context, limit int) ([]model.PipelineRun, error)
}
