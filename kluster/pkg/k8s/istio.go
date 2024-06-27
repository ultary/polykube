package k8s

import (
	"reflect"

	log "github.com/sirupsen/logrus"
	istio "istio.io/client-go/pkg/apis/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ApplyGateway(ctx Context, gateway *istio.Gateway, namespace string) error {

	client := ctx.IstioClientset()

	if gateway.Namespace != "" {
		namespace = gateway.Namespace
	}

	current, err := client.NetworkingV1beta1().Gateways(namespace).Get(ctx, gateway.Name, meta.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
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

	if virtualService.Namespace != "" {
		namespace = virtualService.Namespace
	}

	current, err := client.NetworkingV1beta1().VirtualServices(namespace).Get(ctx, virtualService.Name, meta.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
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
