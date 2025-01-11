package system

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/ultary/polykube/kluster/pkg/k8s"
)

type Manager struct {
	cluster *k8s.Cluster
	db      *gorm.DB

	done chan struct{}
	stop chan struct{}
}

func NewManager(cluster *k8s.Cluster, db *gorm.DB) *Manager {
	return &Manager{
		cluster: cluster,
		db:      db,
	}
}

func (m *Manager) Run() {
	m.done = make(chan struct{})
	m.stop = make(chan struct{})
	m.run()
}

func (m *Manager) Shutdown() {
	log.Info("[SystemManager] Stopping process")

	close(m.stop)
	<-m.done

	log.Info("[SystemManager] Stopped process")
}

func (m *Manager) run() {
	for {
		select {
		case _, ok := <-m.stop:
			if !ok {
				m.done <- struct{}{}
				return
			}
		default:
			time.Sleep(10 * time.Second)
		}

		// TODO: check applications
	}
}
