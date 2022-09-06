package snow

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/bwmarrin/snowflake"
)

func Test_getNodeID(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		nodeID  int64
	}{
		{
			name:    "getNodeID",
			wantErr: false,
			nodeID:  getNodeID(),
		},
		{
			name:    "rand rand id",
			wantErr: false,
			nodeID:  rand.Int63n(2 << 14),
		},
		{
			name:    "max rand id",
			wantErr: true,
			nodeID:  2 << 15,
		},
		{
			name:    "get ip node",
			wantErr: false,
			nodeID:  int64(0)<<8 + int64(1),
		},
		{
			name:    "get ip node",
			wantErr: false,
			nodeID:  int64(255)<<8 + int64(255),
		},
		{
			name:    "get ip node",
			wantErr: false,
			nodeID:  int64(128)<<8 + int64(128),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snowflake.NodeBits = 16
			sf, err := snowflake.NewNode(tt.nodeID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetLocalIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				fmt.Println(sf.Generate().Int64())
			}
		})
	}
}
