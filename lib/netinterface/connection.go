package netinterface

import "time"

type Connection struct {
	SrcIP   string
	SrcPort string
	DstIP   string
	DstPort string
	Proto   string
}

type ConnStats struct {
	PacketCount int
	TotalBytes  int
	LastSeen    time.Time
}

type ConnRow struct {
	Connection
	ConnStats
}
