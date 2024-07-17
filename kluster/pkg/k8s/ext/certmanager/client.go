package certmanager

import (
	"context"
	"reflect"

	certmanager "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Client struct {
	clientset *versioned.Clientset
	informer  *Informer
}

func NewClient(clientset *versioned.Clientset) *Client {
	return &Client{
		clientset: clientset,
	}
}

func (c *Client) ApplyCertificate(ctx context.Context, certificate *certmanager.Certificate, namespace string) error {

	client := c.clientset

	if certificate.Namespace != "" {
		namespace = certificate.Namespace
	}

	current, err := client.CertmanagerV1().Certificates(namespace).Get(ctx, certificate.Name, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}

		if _, err = client.CertmanagerV1().Certificates(namespace).Create(ctx, certificate, metav1.CreateOptions{}); err != nil {
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
	_, err = client.CertmanagerV1().Certificates(namespace).Update(ctx, certificate, metav1.UpdateOptions{})
	if err == nil {
		return nil
	}

	return nil
}
