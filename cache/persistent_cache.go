package cache

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2023/10/20 0:21
 * @file: persistent_cache.go
 * @description:
 */

type Persistent interface {
	Cache
	SaveToFile(filename string) error
	LoadFromFile(filename string) error
}

// SaveToFile PersistentFile is a file-based implementation of Persistent
func SaveToFile(fileName string) error {
	return nil
}
