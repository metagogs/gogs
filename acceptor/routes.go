package acceptor

import "time"

type AcceptroConfig struct {
	Name     string                 // the name of the acceptor
	Groups   []*AcceptorGroupConfig // the groups of the acceptor
	Type     string                 // the type of the acceptor
	HttpPort int
	UdpPort  int
}

type AcceptorGroupConfig struct {
	GroupName          string        // the name of the group
	BucketFillInterval time.Duration // the interval of the bucket fill
	BucketCapacity     int64         // the capacity of the bucket
	Ordered            bool          // whether the message is ordered, only for webrtc
}
