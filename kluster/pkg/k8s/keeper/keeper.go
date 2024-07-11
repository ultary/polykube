package keeper

import (
	"github.com/ultary/monokube/kluster/pkg/k8s"
	"sync"
	"time"

	"github.com/ultary/monokube/kluster/pkg/k8s/keeper/receivers"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
)

type listener struct {
	runners []receivers.Receiver
}

func Init(ctx k8s.Context) *listener {

	// var kubeconfig *string
	// if home := homedir.HomeDir(); home != "" {
	// 	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	// } else {
	// 	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	// }
	// flag.Parse()

	// configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: *kubeconfig}
	// configOverrides := &clientcmd.ConfigOverrides{CurrentContext: "docker-desktop"}

	// config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides).ClientConfig()
	// if err != nil {
	// 	klog.Fatalf("Error building kubeconfig: %s", err.Error())
	// }

	// clientset, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	// }

	clientset := ctx.KubernetesClientset()
	informerFactory := informers.NewSharedInformerFactory(clientset, time.Second*30)
	factories := [...]receivers.NewFunc{
		receivers.NewSecretReceiver,
		receivers.NewConfigMapReceiver,
		receivers.NewDeploymentReceiver,
		receivers.NewStatefulSetReceiver,
		receivers.NewDaemonSetReceiver,
	}
	runners := make([]receivers.Receiver, 0, len(factories))
	for _, factory := range factories {
		runner := factory(informerFactory)
		runners = append(runners, runner)
	}

	informerFactory.Start(wait.NeverStop)

	return &listener{
		runners: runners,
	}
}

func (l *listener) Listen() {

	var wg sync.WaitGroup

	runners := l.runners

	wg.Add(len(runners))
	for _, runner := range runners {
		go func() {
			defer wg.Done()
			runner.Run()
		}()
	}

	wg.Wait()
}
