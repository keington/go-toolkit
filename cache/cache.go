package cache

import (
	"sync"
	"time"
)

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2023/10/17 0:49
 * @file: cache.go
 * @description: 基于sync.Map实现的线程安全本地缓存
 */

// Cache 缓存接口
type Cache interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}

// Entry 条目
type Entry struct {
	value      interface{}
	expiration int64
}

type SafeCache struct {
	syncMap sync.Map
}

// Set 添加key value
func (sc *SafeCache) Set(key string, value interface{}, ttl time.Duration) {
	expiration := time.Now().Add(ttl).UnixNano()
	sc.syncMap.Store(key, Entry{
		value:      value,
		expiration: expiration,
	})
}

// Get 获取指定key的value
func (sc *SafeCache) Get(key string) (interface{}, bool) {
	entry, found := sc.syncMap.Load(key)
	if !found {
		return nil, false
	}

	cacheEntry := entry.(Entry)
	if time.Now().UnixNano() > cacheEntry.expiration {
		sc.syncMap.Delete(key)
		return nil, false
	}
	return cacheEntry.value, true
}

// Delete 删除key
func (sc *SafeCache) Delete(key string) {
	sc.syncMap.Delete(key)
}

// Clean 清理
// 要运行 Clean 方法，可以在初始化缓存时启动一个单独的 Goroutine
// e.g.
//	cache := &SafeCache{}
//	go cache.Clean()
func (sc *SafeCache) Clean() {
	for {
		time.Sleep(1 * time.Minute)
		sc.syncMap.Range(func(key, entry interface{}) bool {
			cacheEntry := entry.(Entry)
			if time.Now().UnixNano() > cacheEntry.expiration {
				sc.syncMap.Delete(key)
			}
			return true
		})
	}
}
