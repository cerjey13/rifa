package background

import (
	"context"
	"time"

	"rifa/backend/pkg/db"
	"rifa/backend/pkg/logx"
)

func StartIdempotencyCleanup(
	ctx context.Context,
	database db.DB,
	logger logx.Logger,
) {
	ticker := time.NewTicker(time.Hour)
	go func() {
		defer ticker.Stop()
		logger.Info(ctx, "Started idempotency cleanup worker")

		for {
			select {
			case <-ticker.C:
				cleanupCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
				defer cancel()

				err := database.ExecContext(
					cleanupCtx,
					`DELETE FROM idempotency_keys WHERE created_at < NOW() - INTERVAL '24 hours'`,
				)
				if err != nil {
					logger.Error(
						ctx,
						"Failed to cleanup idempotency_keys",
						"err",
						err,
					)
				}

			case <-ctx.Done():
				logger.Info(ctx, "Stopping idempotency cleanup worker")
				return
			}
		}
	}()
}
