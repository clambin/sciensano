package reporter

import (
	"context"
	"github.com/clambin/go-common/tabulator"
	"github.com/clambin/sciensano/internal/reports/store"
	"github.com/clambin/sciensano/internal/sciensano"
	"log/slog"
)

type Summary[T summarizer] struct {
	Name   string
	Source Publisher[T]
	Mode   sciensano.SummaryColumn
	Store  *store.Store
	Logger *slog.Logger
}

//go:generate mockery --name Publisher --with-expecter=true
type Publisher[T any] interface {
	Register(chan T)
	Unregister(chan T)
}

type summarizer interface {
	Summarize(column sciensano.SummaryColumn) (*tabulator.Tabulator, error)
}

func (s *Summary[T]) Run(ctx context.Context) error {
	ch := make(chan T)
	s.Source.Register(ch)
	defer func() {
		s.Source.Unregister(ch)
		close(ch)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data := <-ch:
			s.createReport(data)
		}
	}
}

func (s *Summary[T]) createReport(data T) {
	s.Logger.Debug("data received")
	summarized, err := data.Summarize(s.Mode)
	if err != nil {
		s.Logger.Error("failed to generate report", "err", err)
		return
	}
	s.Store.Put(s.Name, summarized)
	s.Logger.Debug("report stored")
}
