package receivers

import "k8s.io/client-go/informers"

type NewFunc func(factory informers.SharedInformerFactory) Receiver

type Receiver interface {
	Run()
}
