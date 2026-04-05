package gradingujian_service

import (
	"context"
	"errors"
)

type GradingUjianPilganExecutor interface {
	GradingUjianPilgan(ctx context.Context, idAttempt int) error
}

type StatistikUjianExecutor interface {
	StatistikUjian(ctx context.Context, idAttempt int) error
}

type CompositeGradingUjianExecutor struct {
	gradingSvc   GradingUjianPilganExecutor
	statistikSvc StatistikUjianExecutor
}

func NewCompositeGradingUjianExecutor(
	gradingSvc GradingUjianPilganExecutor,
	statistikSvc StatistikUjianExecutor,
) *CompositeGradingUjianExecutor {
	return &CompositeGradingUjianExecutor{
		gradingSvc:   gradingSvc,
		statistikSvc: statistikSvc,
	}
}

func (e *CompositeGradingUjianExecutor) GradingUjianPilgan(ctx context.Context, idAttempt int) error {
	if e == nil || e.gradingSvc == nil {
		return errors.New("service grading ujian belum terpasang")
	}

	return e.gradingSvc.GradingUjianPilgan(ctx, idAttempt)
}

func (e *CompositeGradingUjianExecutor) StatistikUjian(ctx context.Context, idAttempt int) error {
	if e == nil || e.statistikSvc == nil {
		return errors.New("service statistik ujian belum terpasang")
	}

	return e.statistikSvc.StatistikUjian(ctx, idAttempt)
}

var _ GradingUjianExecutor = (*CompositeGradingUjianExecutor)(nil)
