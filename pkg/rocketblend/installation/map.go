package installation

import (
	"sync"

	"github.com/rocketblend/rocketblend/pkg/rocketblend/reference"
)

type InstallationMap struct {
	sync.RWMutex
	internal map[reference.Reference]*Installation
}

func NewInstallationMap() *InstallationMap {
	return &InstallationMap{
		internal: make(map[reference.Reference]*Installation),
	}
}

func (cm *InstallationMap) Load(key reference.Reference) (*Installation, bool) {
	cm.RLock()
	result, ok := cm.internal[key]
	cm.RUnlock()
	return result, ok
}

func (cm *InstallationMap) Store(key reference.Reference, value *Installation) {
	cm.Lock()
	cm.internal[key] = value
	cm.Unlock()
}

func (cm *InstallationMap) Delete(key reference.Reference) {
	cm.Lock()
	delete(cm.internal, key)
	cm.Unlock()
}

func (im *InstallationMap) Range(f func(key reference.Reference, value *Installation) bool) {
	im.RLock()
	defer im.RUnlock()
	for k, v := range im.internal {
		if !f(k, v) {
			break
		}
	}
}

func (im *InstallationMap) ToRegularMap() map[reference.Reference]*Installation {
	regMap := make(map[reference.Reference]*Installation)
	im.Range(func(key reference.Reference, value *Installation) bool {
		regMap[key] = value
		return true
	})
	return regMap
}
