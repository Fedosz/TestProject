package telemetry

import (
	grpcotel "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc/stats"
)

func ServerStatsHandler() stats.Handler {
	return grpcotel.NewServerHandler()
}
