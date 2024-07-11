package otlp

import (
	"context"
	_ "embed"
	"github.com/jackc/pgx/v5/pgxpool"

	log "github.com/sirupsen/logrus"
	"github.com/ultary/monokube/kluster/pkg/k8s"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	"sigs.k8s.io/yaml"

	"github.com/ultary/monokube/kluster/pkg/apps"
	"github.com/ultary/monokube/kluster/pkg/helm"
)

type ApplyCommand struct {
	sa      core.ServiceAccount
	cr      rbac.ClusterRole
	crb     rbac.ClusterRoleBinding
	agent   *helm.Chart
	gateway *helm.Chart
}

func NewApplyCommand(manifests apps.Manifests) *ApplyCommand {

	var retval ApplyCommand

	m := manifests.Get("ServiceAccount", "otel-collector")
	if err := yaml.Unmarshal(m, &retval.sa); err != nil {
		log.Fatalf("Error unmarshalling YAML to ServiceAccount: %v", err)
	}

	m = manifests.Get("ClusterRole", "otel-collector")
	if err := yaml.Unmarshal(m, &retval.cr); err != nil {
		log.Fatalf("Error unmarshalling YAML to ClusterRole: %v", err)
	}

	m = manifests.Get("ClusterRoleBinding", "otel-collector")
	if err := yaml.Unmarshal(m, &retval.cr); err != nil {
		log.Fatalf("Error unmarshalling YAML to ClusterRole: %v", err)
	}

	return &retval
}

func (c *ApplyCommand) Do(ctx k8s.Context, namespace string) error {
	if err := k8s.ApplyServiceAccount(ctx, &c.sa, namespace); err != nil {
		log.SetReportCaller(true)
		log.Fatal(err)
		return err
	}
	if err := k8s.ApplyClusterRole(ctx, &c.cr); err != nil {
		log.SetReportCaller(true)
		log.Fatal(err)
		return err
	}
	if err := k8s.ApplyClusterRoleBiding(ctx, &c.crb); err != nil {
		log.SetReportCaller(true)
		log.Fatal(err)
		return err
	}
	if err := c.agent.Apply(ctx, namespace); err != nil {
		log.SetReportCaller(true)
		log.Fatal(err)
		return err
	}
	if err := c.gateway.Apply(ctx, namespace); err != nil {
		log.SetReportCaller(true)
		log.Fatal(err)
		return err
	}
	return nil
}

//go:embed manifests/sa.yaml
var saYaml []byte

//go:embed manifests/cr.yaml
var crYaml []byte

//go:embed manifests/crb.yaml
var crbYaml []byte

func Store(ctx context.Context, pool *pgxpool.Pool) {

}

func Sync(ctx k8s.Context, namespace string) error {

	var sa core.ServiceAccount
	if err := yaml.Unmarshal(saYaml, &sa); err != nil {
		return err
	}
	if err := k8s.ApplyServiceAccount(ctx, &sa, namespace); err != nil {
		log.SetReportCaller(true)
		log.Error(err)
		return err
	}

	var cr rbac.ClusterRole
	if err := yaml.Unmarshal(crYaml, &cr); err != nil {
		return err
	}
	if err := k8s.ApplyClusterRole(ctx, &cr); err != nil {
		log.SetReportCaller(true)
		log.Error(err)
		return err
	}

	var crb rbac.ClusterRoleBinding
	if err := yaml.Unmarshal(crbYaml, &crb); err != nil {
		return err
	}
	if err := k8s.ApplyClusterRoleBiding(ctx, &crb); err != nil {
		log.SetReportCaller(true)
		log.Error(err)
		return err
	}

	return nil
}
