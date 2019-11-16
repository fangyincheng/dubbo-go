package dubbo

import (
	"sync"
	"time"
)

var poolMap *PoolMap

func init() {
	poolMap = new(PoolMap)
}

type PoolMap struct {
	pools map[string]*gettyRPCClientPool
	lock  sync.RWMutex
}

func GetPool(key string, rpcClient *Client) *gettyRPCClientPool {
	poolMap.lock.RLock()
	p := poolMap.pools[key]
	poolMap.lock.RUnlock()
	if p == nil {
		poolMap.lock.Lock()
		p := poolMap.pools[key]
		if p == nil {
			p = newGettyRPCClientConnPool(rpcClient, clientConf.PoolSize, time.Duration(int(time.Second)*clientConf.PoolTTL))
			poolMap.pools[key] = p
		}
		poolMap.lock.Unlock()
	}
	return p
}
