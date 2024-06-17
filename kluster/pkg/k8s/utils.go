package k8s

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"

	certmanager "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	log "github.com/sirupsen/logrus"
	istio "istio.io/client-go/pkg/apis/networking/v1beta1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

type Resource interface {
	Apply(ctx Context, namespace string) error
}

func Exec(ctx Context, namespace, podName, containerName string, command []string) (string, string, error) {

	config := ctx.Config()
	client := ctx.KubernetesClientset()

	req := client.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Namespace(namespace).
		Name(podName).
		SubResource("exec").
		VersionedParams(&core.PodExecOptions{
			Container: containerName,
			Command:   command,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, http.MethodPost, req.URL())
	if err != nil {
		fmt.Errorf("error while creating SPDY executor: %v", err)
		return "", "", err
	}

	var stdout, stderr bytes.Buffer
	err = exec.StreamWithContext(
		ctx,
		remotecommand.StreamOptions{
			Stdin:  nil,
			Stdout: &stdout,
			Stderr: &stderr,
			Tty:    false,
		})
	if err != nil {
		return "", "", fmt.Errorf("error in Stream: %v", err)
	}

	outs := stdout.String()
	errs := stderr.String()
	return outs, errs, nil
}

func ApplyCertificate(ctx Context, certificate *certmanager.Certificate, namespace string) error {

	client := ctx.CertManagerClientset()

	if certificate.Namespace == "" {
		certificate.Namespace = namespace
	} else {
		namespace = certificate.Namespace
	}

	current, err := client.CertmanagerV1().Certificates(namespace).Get(ctx, certificate.Name, meta.GetOptions{})
	if err != nil {
		e, ok := err.(*errors.StatusError)
		if !ok || e.Status().Code != http.StatusNotFound {
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

func ApplyGateway(ctx Context, gateway *istio.Gateway, namespace string) error {

	client := ctx.IstioClientset()

	if gateway.Namespace == "" {
		gateway.Namespace = namespace
	} else {
		namespace = gateway.Namespace
	}

	current, err := client.NetworkingV1beta1().Gateways(namespace).Get(ctx, gateway.Name, meta.GetOptions{})
	if err != nil {
		e, ok := err.(*errors.StatusError)
		if !ok || e.Status().Code != http.StatusNotFound {
			return err
		}

		if _, err = client.NetworkingV1beta1().Gateways(namespace).Create(ctx, gateway, meta.CreateOptions{}); err != nil {
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
	_, err = client.NetworkingV1beta1().Gateways(namespace).Update(ctx, gateway, meta.UpdateOptions{})
	if err == nil {
		return nil
	}

	return nil
}

func ApplyVirtualService(ctx Context, virtualService *istio.VirtualService, namespace string) error {

	client := ctx.IstioClientset()

	if virtualService.Namespace == "" {
		virtualService.Namespace = namespace
	} else {
		namespace = virtualService.Namespace
	}

	current, err := client.NetworkingV1beta1().VirtualServices(namespace).Get(ctx, virtualService.Name, meta.GetOptions{})
	if err != nil {
		e, ok := err.(*errors.StatusError)
		if !ok || e.Status().Code != http.StatusNotFound {
			return err
		}

		if _, err = client.NetworkingV1beta1().VirtualServices(namespace).Create(ctx, virtualService, meta.CreateOptions{}); err != nil {
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
	_, err = client.NetworkingV1beta1().VirtualServices(namespace).Update(ctx, virtualService, meta.UpdateOptions{})
	if err == nil {
		return nil
	}

	return nil
}
