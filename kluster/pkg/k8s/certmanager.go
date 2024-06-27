package k8s

import (
	"reflect"

	certmanager "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ApplyCertificate(ctx Context, certificate *certmanager.Certificate, namespace string) error {

	client := ctx.CertManagerClientset()

	if certificate.Namespace != "" {
		namespace = certificate.Namespace
	}

	current, err := client.CertmanagerV1().Certificates(namespace).Get(ctx, certificate.Name, meta.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}

		if _, err = client.CertmanagerV1().Certificates(namespace).Create(ctx, certificate, meta.CreateOptions{}); err != nil {
			log.Errorf("error creating Certificate: %v", err)
			return err
		}

		return nil
	}

	if reflect.DeepEqual(certificate.Spec, current.Spec) {
		log.Println("Certificate is already up-to-date")
		return nil
	}

	certificate.ResourceVersion = current.ResourceVersion
	_, err = client.CertmanagerV1().Certificates(namespace).Update(ctx, certificate, meta.UpdateOptions{})
	if err == nil {
		return nil
	}

	return nil
}
