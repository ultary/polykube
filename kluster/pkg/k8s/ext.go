package k8s

import (
	"context"
	"reflect"

	certmanager "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	log "github.com/sirupsen/logrus"
	istio "istio.io/client-go/pkg/apis/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

////////////////////////////////////////////////////////////////
//
//  Cert Manager
//

func (c *Client) ApplyCertificate(ctx context.Context, certificate *certmanager.Certificate, namespace string) error {

	client := c.CertManagerClientset()

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

////////////////////////////////////////////////////////////////
//
//  Istio
//

func (c *Client) ApplyGateway(ctx context.Context, gateway *istio.Gateway, namespace string) error {

	client := c.IstioClientset()

	if gateway.Namespace != "" {
		namespace = gateway.Namespace
	}

	current, err := client.NetworkingV1().Gateways(namespace).Get(ctx, gateway.Name, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}

		if _, err = client.NetworkingV1().Gateways(namespace).Create(ctx, gateway, metav1.CreateOptions{}); err != nil {
			log.Errorf("error creating Gateway: %v", err)
			return err
		}

		return nil
	}

	if reflect.DeepEqual(gateway.Spec, current.Spec) {
		log.Println("Gateway is already up-to-date")
		return nil
	}

	gateway.ResourceVersion = current.ResourceVersion
	_, err = client.NetworkingV1().Gateways(namespace).Update(ctx, gateway, metav1.UpdateOptions{})
	if err == nil {
		return nil
	}

	return nil
}

func (c *Client) ApplyVirtualService(ctx context.Context, virtualService *istio.VirtualService, namespace string) error {

	client := c.IstioClientset()

	if virtualService.Namespace != "" {
		namespace = virtualService.Namespace
	}

	current, err := client.NetworkingV1().VirtualServices(namespace).Get(ctx, virtualService.Name, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}

		if _, err = client.NetworkingV1().VirtualServices(namespace).Create(ctx, virtualService, metav1.CreateOptions{}); err != nil {
			log.Errorf("error creating VirtualService: %v", err)
			return err
		}

		return nil
	}

	if reflect.DeepEqual(virtualService.Spec, current.Spec) {
		log.Println("VirtualService is already up-to-date")
		return nil
	}

	virtualService.ResourceVersion = current.ResourceVersion
	_, err = client.NetworkingV1().VirtualServices(namespace).Update(ctx, virtualService, metav1.UpdateOptions{})
	if err == nil {
		return nil
	}

	return nil
}
