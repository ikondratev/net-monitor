package netdevice

import (
	"errors"

	"github.com/google/gopacket/pcap"
)

type Device struct {
	Name        string
	Description string
	Addresses   []string
}

func ListDevices() ([]Device, error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}

	result := make([]Device, 0, len(devices))
	for _, d := range devices {
		addresses := make([]string, 0, len(d.Addresses))
		for _, addr := range d.Addresses {
			if addr.IP != nil {
				addresses = append(addresses, addr.IP.String())
			}
		}

		result = append(result, Device{
			Name:        d.Name,
			Description: d.Description,
			Addresses:   addresses,
		})
	}

	return result, nil
}

func FindActiveDevice() (string, error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return "", err
	}
	for _, d := range devices {
		if len(d.Addresses) > 0 {
			for _, addr := range d.Addresses {
				if !addr.IP.IsLoopback() {
					return d.Name, nil
				}
			}
		}
	}
	return "", errors.New("no active network interface found")
}
