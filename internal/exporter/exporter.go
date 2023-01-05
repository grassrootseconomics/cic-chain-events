package exporter

import (
	"github.com/VictoriaMetrics/metrics"
	"github.com/grassrootseconomics/cic-chain-events/internal/syncer"
)

func Register(stats *syncer.Stats) {
	metrics.NewGauge("indexer_head_cursor", func() float64 {
		return float64(stats.GetHeadCursor())
	})

	metrics.NewGauge("indexer_lower_bound", func() float64 {
		return float64(stats.GetLowerBound())
	})

	metrics.NewGauge("indexer_missing_blocks", func() float64 {
		if stats.GetHeadCursor()-stats.GetLowerBound() < 10 {
			return float64(0)
		} else {
			return float64(stats.GetHeadCursor() - stats.GetLowerBound())
		}
	})
}
