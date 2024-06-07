package utils

import (
	"github.com/bwmarrin/snowflake"
)

var snowNode *snowflake.Node

func init() {
	dataCenterID := int64(0)
	machineID := int64(1)
	nodeID := (dataCenterID << 5) | machineID
	snowNode, _ = snowflake.NewNode(nodeID)
	// if err != nil {
	// 	return 0, fmt.Errorf("Error creating new Node: %s ", err.Error())
	// }
}

func GenerateSnowId() (int64, error) {
	return snowNode.Generate().Int64(), nil
}
