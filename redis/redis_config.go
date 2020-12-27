package redis

import (
	"fmt"
	"sync"
)

const (
	DefaultGroupName   = "default" // Default group name.
)

// Config is the configuration management object.
type Config map[string]ConfigGroup

// ConfigGroup is a slice of configuration node for specified named group.
type ConfigGroup []ConfigNode

// ConfigNode is configuration for one node.
type ConfigNode struct {
	Host                 string        // Host of server, ip or domain like: 127.0.0.1, localhost
	Port                 string        // Port, it's commonly 3306.
	User                 string        // Authentication username.
	Pass                 string        // Authentication password.
	Db                   string        // Default used database name.
	Prefix               string        // (Optional) Table prefix.
	PoolSize     		int           `json:"pool"`     // Maximum number of socket connections.
	MinIdleConns 		int            // Minimum number of idle connections which is useful when establishing
}

// configs is internal used configuration object.
var configs struct {
	sync.RWMutex
	config Config // All configurations.
	group  string // Default configuration group.
}

func init() {
	configs.config = make(Config)
	configs.group = DefaultGroupName
}

// SetConfig sets the global configuration for package.
// It will overwrite the old configuration of package.
func SetConfig(config Config) {
	configs.Lock()
	defer configs.Unlock()
	configs.config = config
}

// SetConfigGroup sets the configuration for given group.
func SetConfigGroup(group string, nodes ConfigGroup) {
	configs.Lock()
	defer configs.Unlock()
	configs.config[group] = nodes
}

// AddConfigNode adds one node configuration to configuration of given group.
func AddConfigNode(group string, node ConfigNode) {
	configs.Lock()
	defer configs.Unlock()
	configs.config[group] = append(configs.config[group], node)
}

// AddDefaultConfigNode adds one node configuration to configuration of default group.
func AddDefaultConfigNode(node ConfigNode) {
	AddConfigNode(DefaultGroupName, node)
}

// AddDefaultConfigGroup adds multiple node configurations to configuration of default group.
func AddDefaultConfigGroup(nodes ConfigGroup) {
	SetConfigGroup(DefaultGroupName, nodes)
}

// GetConfig retrieves and returns the configuration of given group.
func GetConfig(group string) ConfigGroup {
	configs.RLock()
	defer configs.RUnlock()
	return configs.config[group]
}

// SetDefaultGroup sets the group name for default configuration.
func SetDefaultGroup(name string) {
	configs.Lock()
	defer configs.Unlock()
	configs.group = name
}

// GetDefaultGroup returns the { name of default configuration.
func GetDefaultGroup() string {
	configs.RLock()
	defer configs.RUnlock()
	return configs.group
}

// IsConfigured checks and returns whether the database configured.
// It returns true if any configuration exists.
func IsConfigured() bool {
	configs.RLock()
	defer configs.RUnlock()
	return len(configs.config) > 0
}

// String returns the node as string.
func (node *ConfigNode) String() string {
	return fmt.Sprintf(
		`%s@%s:%s,%s,%s,%s,%s,%v,%d-%d-%d#%s`,
		node.User, node.Host, node.Port,
		node.MinIdleConns, node.PoolSize,
	)
}