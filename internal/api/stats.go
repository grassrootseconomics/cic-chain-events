package api

import (
	"net/http"

	"github.com/alitto/pond"
	"github.com/grassrootseconomics/cic-chain-events/internal/syncer"
	"github.com/labstack/echo/v4"
)

type statsResponse struct {
	HeadCursor            uint64 `json:"headCursor"`
	LowerBound            uint64 `json:"lowerBound"`
	MissingBlocks         uint64 `json:"missingBlocks"`
	WorkerQueueSize       uint64 `json:"workerQueueSize"`
	WorkerCount           int    `json:"workerCount"`
	WorkerSuccessfulTasks uint64 `json:"workerSuccessfulTasks"`
	WorkerFailedTasks     uint64 `json:"workerFailedTasks"`
}

func StatsHandler(
	syncerStats *syncer.Stats,
	poolStats *pond.WorkerPool,
) func(echo.Context) error {
	return func(ctx echo.Context) error {
		headCursor := syncerStats.GetHeadCursor()
		lowerBound := syncerStats.GetLowerBound()

		stats := statsResponse{
			HeadCursor:            headCursor,
			LowerBound:            lowerBound,
			WorkerCount:           poolStats.RunningWorkers(),
			WorkerQueueSize:       poolStats.WaitingTasks(),
			WorkerSuccessfulTasks: poolStats.SuccessfulTasks(),
			WorkerFailedTasks:     poolStats.FailedTasks(),
		}

		if headCursor-lowerBound < 10 {
			stats.MissingBlocks = 0
		} else {
			stats.MissingBlocks = headCursor - lowerBound
		}

		return ctx.JSON(http.StatusOK, okResp{
			Ok:   true,
			Data: stats,
		})
	}
}
