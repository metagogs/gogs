package snow

import (
	"math/rand"

	"github.com/bwmarrin/snowflake"
)

func getNodeID() int64 {
	var nodeID int64
	ip, err := GetLocalIP()
	if err != nil || len(ip) == 0 {
		nodeID = rand.Int63n(2 << 14) //nolint
	} else {
		nodeID = IP4toInt16(ip)
	}

	return nodeID
}

func NewSnowNode() (*snowflake.Node, error) {
	nodeID := getNodeID()
	snowflake.NodeBits = 16
	snowflake.StepBits = 6
	sf, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, err
	}

	return sf, nil
}
