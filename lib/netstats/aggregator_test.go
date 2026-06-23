package netstats

import (
	"encoding/binary"
	"net"
	"testing"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"github.com/ikondratev/net-monitor/lib/netinterface"
)

func TestProcessPacketAggregatesIPv4TCPConnection(t *testing.T) {
	aggregator := NewAggregator()
	packet := newIPv4TCPPacket(t, "192.0.2.10", "198.51.100.20", 50000, 50001)

	aggregator.ProcessPacket(packet)
	aggregator.ProcessPacket(packet)

	rows := aggregator.ConnectionRows()
	if len(rows) != 1 {
		t.Fatalf("expected 1 connection row, got %d", len(rows))
	}

	row := rows[0]
	if row.SrcIP != "192.0.2.10" || row.DstIP != "198.51.100.20" {
		t.Fatalf("unexpected endpoints: %+v", row.Connection)
	}
	if row.SrcPort != "50000" || row.DstPort != "50001" || row.Proto != "TCP" {
		t.Fatalf("unexpected connection metadata: %+v", row.Connection)
	}
	if row.PacketCount != 2 {
		t.Fatalf("expected packet count 2, got %d", row.PacketCount)
	}
	if row.TotalBytes != 2*len(packet.Data()) {
		t.Fatalf("expected %d bytes, got %d", 2*len(packet.Data()), row.TotalBytes)
	}
	if row.LastSeen.IsZero() {
		t.Fatal("expected LastSeen to be set")
	}
}

func TestProcessPacketIgnoresNonIPv4Packets(t *testing.T) {
	aggregator := NewAggregator()
	packet := gopacket.NewPacket(make([]byte, 40), layers.LayerTypeIPv6, gopacket.Default)

	aggregator.ProcessPacket(packet)

	if rows := aggregator.ConnectionRows(); len(rows) != 0 {
		t.Fatalf("expected no rows for non-IPv4 packet, got %d", len(rows))
	}
}

func TestConnectionRowsSortsByBytesThenPacketsAndExpiresStaleRows(t *testing.T) {
	aggregator := NewAggregator()
	now := time.Now()
	aggregator.networkMap = map[netinterface.Connection]*netinterface.ConnStats{
		{SrcIP: "10.0.0.1", SrcPort: "1000", DstIP: "10.0.0.2", DstPort: "2000", Proto: "TCP"}: {
			PacketCount: 2,
			TotalBytes:  200,
			LastSeen:    now,
		},
		{SrcIP: "10.0.0.3", SrcPort: "1001", DstIP: "10.0.0.4", DstPort: "2001", Proto: "TCP"}: {
			PacketCount: 3,
			TotalBytes:  200,
			LastSeen:    now,
		},
		{SrcIP: "10.0.0.5", SrcPort: "1002", DstIP: "10.0.0.6", DstPort: "2002", Proto: "UDP"}: {
			PacketCount: 10,
			TotalBytes:  100,
			LastSeen:    now,
		},
		{SrcIP: "10.0.0.7", SrcPort: "1003", DstIP: "10.0.0.8", DstPort: "2003", Proto: "TCP"}: {
			PacketCount: 100,
			TotalBytes:  1000,
			LastSeen:    now.Add(-connectionTTL - time.Second),
		},
	}

	rows := aggregator.ConnectionRows()
	if len(rows) != 3 {
		t.Fatalf("expected 3 active rows, got %d", len(rows))
	}

	if rows[0].PacketCount != 3 || rows[0].TotalBytes != 200 {
		t.Fatalf("expected row with more packets first on byte tie, got %+v", rows[0].ConnStats)
	}
	if rows[1].PacketCount != 2 || rows[1].TotalBytes != 200 {
		t.Fatalf("expected second row with same bytes and fewer packets, got %+v", rows[1].ConnStats)
	}
	if rows[2].TotalBytes != 100 {
		t.Fatalf("expected lowest-volume row last, got %+v", rows[2].ConnStats)
	}
}

func newIPv4TCPPacket(t *testing.T, srcIP, dstIP string, srcPort, dstPort uint16) gopacket.Packet {
	t.Helper()

	const (
		ipHeaderLength  = 20
		tcpHeaderLength = 20
	)
	data := make([]byte, ipHeaderLength+tcpHeaderLength)
	data[0] = 0x45
	binary.BigEndian.PutUint16(data[2:4], uint16(len(data)))
	data[8] = 64
	data[9] = byte(layers.IPProtocolTCP)
	copy(data[12:16], net.ParseIP(srcIP).To4())
	copy(data[16:20], net.ParseIP(dstIP).To4())

	tcp := data[ipHeaderLength:]
	binary.BigEndian.PutUint16(tcp[0:2], srcPort)
	binary.BigEndian.PutUint16(tcp[2:4], dstPort)
	tcp[12] = 0x50

	packet := gopacket.NewPacket(data, layers.LayerTypeIPv4, gopacket.Default)
	packet.Metadata().CaptureLength = len(data)
	packet.Metadata().Length = len(data)
	return packet
}
