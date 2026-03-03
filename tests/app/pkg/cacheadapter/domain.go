// Package cacheadapter provides helpers and adapters for namespaced caching operations.
package cacheadapter

func makeID(namespace, key string) string {
	return namespace + "/" + key
}
