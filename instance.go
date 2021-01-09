package easy

import (
	"github.com/gogf/gf/container/gmap"
)

var (
	// Instances is the instance map for common used components.
	Instances = gmap.NewStrAnyMap(true)
)

// Get returns the instance by given name.
func Get(name string) interface{} {
	return Instances.Get(name)
}

// Set sets a instance object to the instance manager with given name.
func Set(name string, instance interface{}) {
	Instances.Set(name, instance)
}

// GetOrSet returns the instance by name,
// or set instance to the instance manager if it does not exist and returns this instance.
func GetOrSet(name string, instance interface{}) interface{} {
	return Instances.GetOrSet(name, instance)
}

// GetOrSetFunc returns the instance by name,
// or sets instance with returned value of callback function <f> if it does not exist
// and then returns this instance.
func GetOrSetFunc(name string, f func() interface{}) interface{} {
	return Instances.GetOrSetFunc(name, f)
}

// GetOrSetFuncLock returns the instance by name,
// or sets instance with returned value of callback function <f> if it does not exist
// and then returns this instance.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function <f>
// with mutex.Lock of the hash map.
func GetOrSetFuncLock(name string, f func() interface{}) interface{} {
	return Instances.GetOrSetFuncLock(name, f)
}

// SetIfNotExist sets <instance> to the map if the <name> does not exist, then returns true.
// It returns false if <name> exists, and <instance> would be ignored.
func SetIfNotExist(name string, instance interface{}) bool {
	return Instances.SetIfNotExist(name, instance)
}

