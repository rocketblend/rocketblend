package rocketpack

import (
	"sync"

	"github.com/rocketblend/rocketblend/pkg/rocketblend/reference"
)

type RocketPackMap struct {
	sync.RWMutex
	internal map[reference.Reference]*RocketPack
}

func NewRocketPackMap() *RocketPackMap {
	return &RocketPackMap{
		internal: make(map[reference.Reference]*RocketPack),
	}
}

func (cm *RocketPackMap) Load(key reference.Reference) (*RocketPack, bool) {
	cm.RLock()
	result, ok := cm.internal[key]
	cm.RUnlock()
	return result, ok
}

func (cm *RocketPackMap) Store(key reference.Reference, value *RocketPack) {
	cm.Lock()
	cm.internal[key] = value
	cm.Unlock()
}

func (cm *RocketPackMap) Delete(key reference.Reference) {
	cm.Lock()
	delete(cm.internal, key)
	cm.Unlock()
}

func (im *RocketPackMap) Range(f func(key reference.Reference, value *RocketPack) bool) {
	im.RLock()
	defer im.RUnlock()
	for k, v := range im.internal {
		if !f(k, v) {
			break
		}
	}
}

func (im *RocketPackMap) ToRegularMap() map[reference.Reference]*RocketPack {
	regMap := make(map[reference.Reference]*RocketPack)
	im.Range(func(key reference.Reference, value *RocketPack) bool {
		regMap[key] = value
		return true
	})
	return regMap
}
