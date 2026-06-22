package netcapture

import (
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

const (
	snapshotLength = 256
	promiscuous    = false
	readTimeout    = 1 * time.Second
	packetFilter   = "ip and (tcp or udp or icmp)"
)

type Capture struct {
	handle *pcap.Handle
	source *gopacket.PacketSource
}

func Open(device string) (*Capture, error) {
	handle, err := pcap.OpenLive(device, snapshotLength, promiscuous, readTimeout)
	if err != nil {
		return nil, err
	}

	if err := handle.SetBPFFilter(packetFilter); err != nil {
		handle.Close()
		return nil, err
	}

	return &Capture{
		handle: handle,
		source: gopacket.NewPacketSource(handle, handle.LinkType()),
	}, nil
}

func (c *Capture) Packets() chan gopacket.Packet {
	return c.source.Packets()
}

func (c *Capture) Close() {
	c.handle.Close()
}
