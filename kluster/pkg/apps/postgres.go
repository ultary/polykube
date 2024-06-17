package apps

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/yaml"

	"ultary.co/kluster/pkg/helm"
	"ultary.co/kluster/pkg/k8s"
	"ultary.co/kluster/pkg/utils"
)

type PostgreSQL struct {
	secret core.Secret
	sts    apps.StatefulSet
	sv     core.Service
}

func NewPostgreSQL(chart *helm.Chart) (retval PostgreSQL) {

	const name = "postgres"

	m := chart.Get("Secret", name)
	if err := yaml.Unmarshal(m, &retval.secret); err != nil {
		log.Fatalf("Error unmarshalling YAML to Secret: %v", err)
	}

	m = chart.Get("StatefulSet", name)
	if err := yaml.Unmarshal(m, &retval.sts); err != nil {
		log.Fatalf("Error unmarshalling YAML to StatefulSet: %v", err)
	}

	m = chart.Get("Service", name)
	if err := yaml.Unmarshal(m, &retval.sv); err != nil {
		log.Fatalf("Error unmarshalling YAML to Service: %v", err)
	}

	return
}

func (p *PostgreSQL) Apply(ctx k8s.Context, namespace string) error {

	{
		var result *core.Secret
		name := p.secret.Name
		result, err := k8s.GetSecret(ctx, name, namespace)
		if err != nil {

			status := err.(*errors.StatusError).ErrStatus
			if status.Code != http.StatusNotFound {
				log.Fatalf("Error getting PostgreSQL Secret: %v", err)
			}

			password := utils.NewPassword()
			p.secret.StringData = map[string]string{
				"POSTGRES_PASSWORD": password,
			}

			result, err = k8s.CreateSecret(ctx, namespace, &p.secret)
			if err != nil {
				log.Fatalf("Error creating PostgreSQL Secret: %v", err)
			}

			log.Debug(result)
		}
	}

	if err := k8s.ApplyStatefulSet(ctx, &p.sts, namespace); err != nil {
		log.Fatalf("Error applying PostgreSQL StatefulSet: %v", err)
	}

	if err := k8s.ApplyService(ctx, &p.sv, namespace); err != nil {
		log.Fatalf("Error applying PostgreSQL Service: %v", err)
	}

	return nil
}
