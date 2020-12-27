package gorm

import (
	"fmt"
	"github.com/go-ecosystem/mysql"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/gins"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/grand"
	"github.com/gogf/gf/util/gutil"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	frameCoreComponentNameDatabase = "gfplus.core.component.database"
	configNodeNameDatabase         = "database"
)

func Database(name ...string) *gorm.DB {
	group := DefaultGroupName
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	var (
		configMap     map[string]interface{}
		configNodeKey string
	)
	// It firstly searches the configuration of the instance name.
		configNodeKey, _ = gutil.MapPossibleItemByKey(
			gins.Config().GetMap("."),
			configNodeNameDatabase,
		)
		if configNodeKey == "" {
			configNodeKey = configNodeNameDatabase
		}
		configMap = gins.Config().GetMap(configNodeKey)

	if len(configMap) == 0 && !IsConfigured() {
		panic(fmt.Sprintf(`database init failed: "%s" node not found, is config file or configuration missing?`, configNodeNameDatabase))
	}
	if len(configMap) == 0 {
		configMap = make(map[string]interface{})
	}
	// Parse <m> as map-slice and adds it to gdb's global configurations.
	for g, groupConfig := range configMap {
		cg := ConfigGroup{}
		switch value := groupConfig.(type) {
		case []interface{}:
			for _, v := range value {
				if node := parseDBConfigNode(v); node != nil {
					cg = append(cg, *node)
				}
			}
		case map[string]interface{}:
			if node := parseDBConfigNode(value); node != nil {
				cg = append(cg, *node)
			}
		}
		if len(cg) > 0 {
			if GetConfig(group) == nil {
				glog.Printf("add configuration for group: %s, %#v", g, cg)
				SetConfigGroup(g, cg)
			} else {
				glog.Printf("ignore configuration as it already exists for group: %s, %#v", g, cg)
				glog.Printf("%s, %#v", g, cg)
			}
		}
	}
	// Parse <m> as a single node configuration,
	// which is the default group configuration.
	if node := parseDBConfigNode(configMap); node != nil {
		cg := ConfigGroup{}
		if node.LinkInfo != "" || node.Host != "" {
			cg = append(cg, *node)
		}

		if len(cg) > 0 {
			if GetConfig(group) == nil {
				glog.Printf("add configuration for group: %s, %#v", DefaultGroupName, cg)
				SetConfigGroup(DefaultGroupName, cg)
			} else {
				glog.Printf("ignore configuration as it already exists for group: %s, %#v", DefaultGroupName, cg)
				glog.Printf("%s, %#v", DefaultGroupName, cg)
			}
		}
	}
	gormDb := mysql.GetDBByKey(group)
	if gormDb != nil {
		return gormDb
	}
	// Create a new ORM object with given configurations.
	if db, err := New(name...); err == nil {
		return db
	}
	return nil
}

// New creates and returns an ORM object with global configurations.
// The parameter <name> specifies the configuration group name,
// which is DefaultGroupName in default.
func New(group ...string) (db *gorm.DB, err error) {
	groupName := configs.group
	if len(group) > 0 && group[0] != "" {
		groupName = group[0]
	}
	configs.RLock()
	defer configs.RUnlock()

	if len(configs.config) < 1 {
		return nil, gerror.New("empty database configuration")
	}
	if _, ok := configs.config[groupName]; ok {
		if node, err := getConfigNodeByGroup(groupName, true); err == nil {
			cnf := mysql.NewConfig(
				node.User,
				node.Pass,
				node.Host,
				node.Port,
				node.Name,
				node.Charset,
				logger.Error,
				mysql.WithMaxOpenConns(node.MaxOpenConnCount),
				mysql.WithMaxIdleConns(node.MaxIdleConnCount))
			// register default db
			mysql.RegisterByKey(cnf, groupName)
			db := mysql.GetDBByKey(groupName)
			if db != nil {
				return db, nil
			} else {
				return nil, gerror.Newf(`database node "%s" is not found`, groupName)
			}
		} else {
			return nil, err
		}
	} else {
		return nil, gerror.New(fmt.Sprintf(`database configuration node "%s" is not found`, groupName))
	}
}

// getConfigNodeByGroup calculates and returns a configuration node of given group. It
// calculates the value internally using weight algorithm for load balance.
//
// The parameter <master> specifies whether retrieving a master node, or else a slave node
// if master-slave configured.
func getConfigNodeByGroup(group string, master bool) (*ConfigNode, error) {
	if list, ok := configs.config[group]; ok {
		// Separates master and slave configuration nodes array.
		masterList := make(ConfigGroup, 0)
		slaveList := make(ConfigGroup, 0)
		for i := 0; i < len(list); i++ {
			if list[i].Role == "slave" {
				slaveList = append(slaveList, list[i])
			} else {
				masterList = append(masterList, list[i])
			}
		}
		if len(masterList) < 1 {
			return nil, gerror.New("at least one master node configuration's need to make sense")
		}
		if len(slaveList) < 1 {
			slaveList = masterList
		}
		if master {
			return getConfigNodeByWeight(masterList), nil
		} else {
			return getConfigNodeByWeight(slaveList), nil
		}
	} else {
		return nil, gerror.New(fmt.Sprintf("empty database configuration for item name '%s'", group))
	}
}

// getConfigNodeByWeight calculates the configuration weights and randomly returns a node.
//
// Calculation algorithm brief:
// 1. If we have 2 nodes, and their weights are both 1, then the weight range is [0, 199];
// 2. Node1 weight range is [0, 99], and node2 weight range is [100, 199], ratio is 1:1;
// 3. If the random number is 99, it then chooses and returns node1;
func getConfigNodeByWeight(cg ConfigGroup) *ConfigNode {
	if len(cg) < 2 {
		return &cg[0]
	}
	var total int
	for i := 0; i < len(cg); i++ {
		total += cg[i].Weight * 100
	}
	// If total is 0 means all of the nodes have no weight attribute configured.
	// It then defaults each node's weight attribute to 1.
	if total == 0 {
		for i := 0; i < len(cg); i++ {
			cg[i].Weight = 1
			total += cg[i].Weight * 100
		}
	}
	// Exclude the right border value.
	r := grand.N(0, total-1)
	min := 0
	max := 0
	for i := 0; i < len(cg); i++ {
		max = min + cg[i].Weight*100
		//fmt.Printf("r: %d, min: %d, max: %d\n", r, min, max)
		if r >= min && r < max {
			return &cg[i]
		} else {
			min = max
		}
	}
	return nil
}

func parseDBConfigNode(value interface{}) *ConfigNode {
	nodeMap, ok := value.(map[string]interface{})
	if !ok {
		return nil
	}
	node := &ConfigNode{}
	err := gconv.Struct(nodeMap, node)
	if err != nil {
		panic(err)
	}
	if _, v := gutil.MapPossibleItemByKey(nodeMap, "link"); v != nil {
		node.LinkInfo = gconv.String(v)
	}
	// Parse link syntax.
	if node.LinkInfo != "" && node.Type == "" {
		match, _ := gregex.MatchString(`([a-z]+):(.+)`, node.LinkInfo)
		if len(match) == 3 {
			node.Type = gstr.Trim(match[1])
			node.LinkInfo = gstr.Trim(match[2])
		}
	}
	return node
}