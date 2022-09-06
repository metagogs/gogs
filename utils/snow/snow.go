package snow

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/snowflake"
)

func getNodeID() int64 {
	var nodeID int64
	ip, err := GetLocalIP()
	if err != nil || len(ip) == 0 {
		rand.Seed(time.Now().UnixNano())
		nodeID = rand.Int63n(2 << 14) //nolint
	} else {
		nodeID = IP4toInt16(ip)
	}

	return nodeID
}

func NewSnowNode() (*snowflake.Node, error) {
	nodeID := getNodeID()
	snowflake.NodeBits = 16
	sf, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, err
	}

	return sf, nil
}
