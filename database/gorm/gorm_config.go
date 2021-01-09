package gorm

import (
	"fmt"
	"sync"
	"time"
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
	Name                 string        // Default used database name.
	Type                 string        // Database type: mysql, sqlite, mssql, pgsql, oracle.
	Role                 string        // (Optional, "master" in default) Node role, used for master-slave mode: master, slave.
	Debug                bool          // (Optional) Debug mode enables debug information logging and output.
	Prefix               string        // (Optional) Table prefix.
	DryRun               bool          // (Optional) Dry run, which does SELECT but no INSERT/UPDATE/DELETE statements.
	Weight               int           // (Optional) Weight for load balance calculating, it's useless if there's just one node.
	Charset              string        // (Optional, "utf8mb4" in default) Custom charset when operating on database.
	LinkInfo             string        `json:"link"`        // (Optional) Custom link information, when it is used, configuration Host/Port/User/Pass/Name are ignored.
	MaxIdleConnCount     int           `json:"maxidle"`     // (Optional) Max idle connection configuration for underlying connection pool.
	MaxOpenConnCount     int           `json:"maxopen"`     // (Optional) Max open connection configuration for underlying connection pool.
	MaxConnLifetime      time.Duration `json:"maxlifetime"` // (Optional) Max connection TTL configuration for underlying connection pool.
	CreatedAt            string        // (Optional) The filed name of table for automatic-filled created datetime.
	UpdatedAt            string        // (Optional) The filed name of table for automatic-filled updated datetime.
	DeletedAt            string        // (Optional) The filed name of table for automatic-filled updated datetime.
	TimeMaintainDisabled bool          // (Optional) Disable the automatic time maintaining feature.
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
		node.Name, node.Type, node.Role, node.Charset, node.Debug,
		node.MaxIdleConnCount,
		node.MaxOpenConnCount,
		node.MaxConnLifetime,
		node.LinkInfo,
	)
}