package grafana

import (
	"context"
	"ultary.co/kluster/pkg/k8s"
	"ultary.co/kluster/pkg/repo"
)

func Apply(ctx context.Context, client *k8s.Client, manifests *repo.Manifests) {

	// postgres.CreateDatabase(ctx, client, "grafana")
	// postgres.CreateRole(client)
}
