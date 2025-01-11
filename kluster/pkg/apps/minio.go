package apps

/*
import (
	log "github.com/sirupsen/logrus"
	istio "istio.io/client-go/pkg/apis/networking/v1"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/yaml"

	"github.com/ultary/polykube/kluster/pkg/k8s"
	"github.com/ultary/polykube/kluster/pkg/utils"
)

type MinIO struct {
	secret core.Secret
	sts    apps.StatefulSet
	sv     core.Service
	vsv    istio.VirtualService
}

func NewMinIO(manifests Manifests) *MinIO {

	const name = "minio"

	var retval MinIO

	m := manifests.Get("Secret", name)
	if err := yaml.Unmarshal(m, &retval.secret); err != nil {
		log.Fatalf("Error unmarshalling YAML to Secret: %v", err)
	}

	m = manifests.Get("StatefulSet", name)
	if err := yaml.Unmarshal(m, &retval.sts); err != nil {
		log.Fatalf("Error unmarshalling YAML to StatefulSet: %v", err)
	}

	m = manifests.Get("Service", name)
	if err := yaml.Unmarshal(m, &retval.sv); err != nil {
		log.Fatalf("Error unmarshalling YAML to Service: %v", err)
	}

	m = manifests.Get("VirtualService", name)
	if err := yaml.Unmarshal(m, &retval.vsv); err != nil {
		log.Fatalf("Error unmarshalling YAML to VirtualService: %v", err)
	}

	return &retval
}

func (m *MinIO) Apply(ctx k8s.Context, namespace string) error {

	{
		var result *core.Secret
		name := m.secret.Name
		result, err := k8s.GetSecret(ctx, name, namespace)
		if err != nil {

			if !errors.IsNotFound(err) {
				log.Fatalf("Error getting MinIO Secret: %v", err)
			}

			password := utils.NewPassword()
			m.secret.StringData = map[string]string{
				"MINIO_ROOT_PASSWORD": password,
			}

			result, err = k8s.CreateSecret(ctx, namespace, &m.secret)
			if err != nil {
				log.Fatalf("Error creating MinIO Secret: %v", err)
			}

			log.Debug(result)
		}
	}

	if err := k8s.ApplyStatefulSet(ctx, &m.sts, namespace); err != nil {
		log.Fatalf("Error applying MinIO StatefulSet: %v", err)
	}

	if err := k8s.ApplyService(ctx, &m.sv, namespace); err != nil {
		log.Fatalf("Error applying MinIO Service: %v", err)
	}

	if err := k8s.ApplyVirtualService(ctx, &m.vsv, namespace); err != nil {
		log.Fatalf("Error applying MinIO VirtualService: %v", err)
	}

	return nil
}

func (m *MinIO) CreateAccessKey(ctx k8s.Context, name, bucket string) (accessKey, secretKey string, err error) {
	return "", "", nil
}
*/
