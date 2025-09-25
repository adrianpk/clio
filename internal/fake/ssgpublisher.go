package fake

import (
	"context"

	"github.com/adrianpk/clio/internal/feat/ssg"
)

// SsgPublisher is a fake implementation of ssg.Publisher for testing.
type SsgPublisher struct {
	// Expected results
	ValidateFn func(cfg ssg.PublisherConfig) error
	PublishFn  func(ctx context.Context, cfg ssg.PublisherConfig, sourceDir string) (string, error)
	PlanFn     func(ctx context.Context, cfg ssg.PublisherConfig, sourceDir string) (ssg.PlanReport, error)

	// Captured arguments
	ValidateCalls []struct{ Cfg ssg.PublisherConfig }
	PublishCalls  []struct{ Ctx context.Context; Cfg ssg.PublisherConfig; SourceDir string }
	PlanCalls     []struct{ Ctx context.Context; Cfg ssg.PublisherConfig; SourceDir string }
}

// NewSsgPublisher creates a new fake SsgPublisher.
func NewSsgPublisher() *SsgPublisher {
	return &SsgPublisher{}
}

func (f *SsgPublisher) Validate(cfg ssg.PublisherConfig) error {
	f.ValidateCalls = append(f.ValidateCalls, struct{ Cfg ssg.PublisherConfig }{Cfg: cfg})
	if f.ValidateFn != nil {
		return f.ValidateFn(cfg)
	}
	return nil
}

func (f *SsgPublisher) Publish(ctx context.Context, cfg ssg.PublisherConfig, sourceDir string) (string, error) {
	f.PublishCalls = append(f.PublishCalls, struct{ Ctx context.Context; Cfg ssg.PublisherConfig; SourceDir string }{Ctx: ctx, Cfg: cfg, SourceDir: sourceDir})
	if f.PublishFn != nil {
		return f.PublishFn(ctx, cfg, sourceDir)
	}
	return "fake-commit-url", nil
}

func (f *SsgPublisher) Plan(ctx context.Context, cfg ssg.PublisherConfig, sourceDir string) (ssg.PlanReport, error) {
	f.PlanCalls = append(f.PlanCalls, struct{ Ctx context.Context; Cfg ssg.PublisherConfig; SourceDir string }{Ctx: ctx, Cfg: cfg, SourceDir: sourceDir})
	if f.PlanFn != nil {
		return f.PlanFn(ctx, cfg, sourceDir)
	}
	return ssg.PlanReport{Summary: "fake plan"}, nil
}
