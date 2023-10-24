package cache

import (
	"testing"
	"time"
)

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2023/10/20 0:10
 * @file: cache_test.go
 * @description:  基于sync.Map实现的线程安全本地缓存单元测试
 */

func TestSafeCache_SetAndGet(t *testing.T) {
	cache := &SafeCache{}
	key := "test_key"
	value := "test_value"
	ttl := 2 * time.Second

	cache.Set(key, value, ttl)

	// 获取缓存项
	cachedValue, found := cache.Get(key)

	if !found {
		t.Errorf("Expected to find value for key %s, but it was not found", key)
	}

	if cachedValue != value {
		t.Errorf("Expected cached value to be %s, but got %v", value, cachedValue)
	}
}

func TestSafeCache_Delete(t *testing.T) {
	cache := &SafeCache{}
	key := "test_key"
	value := "test_value"
	ttl := 2 * time.Second

	cache.Set(key, value, ttl)
	cache.Delete(key)

	// 尝试获取已删除的缓存项
	cachedValue, found := cache.Get(key)

	if found {
		t.Errorf("Expected not to find a value for key %s after deletion, but found %v", key, cachedValue)
	}
}

func TestSafeCache_Expiration(t *testing.T) {
	cache := &SafeCache{}
	key := "test_key"
	value := "test_value"
	ttl := 2 * time.Second

	cache.Set(key, value, ttl)

	// 等待过期
	time.Sleep(3 * time.Second)

	// 尝试获取已过期的缓存项
	cachedValue, found := cache.Get(key)

	if found {
		t.Errorf("Expected not to find a value for key %s after expiration, but found %v", key, cachedValue)
	}
}
