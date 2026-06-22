package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/ikondratev/net-monitor/lib/app"
	"github.com/ikondratev/net-monitor/lib/netcapture"
	"github.com/ikondratev/net-monitor/lib/netdevice"
	"github.com/ikondratev/net-monitor/lib/netstats"
)

type Config struct {
	ShowInterfaces bool
	Device         string
}

func Run(args []string, stdout io.Writer) error {
	cfg, err := parseFlags(args)
	if err != nil {
		return err
	}

	if cfg.ShowInterfaces {
		return printInterfaces(stdout)
	}

	device := cfg.Device
	if device == "" {
		device, err = netdevice.FindActiveDevice()
		if err != nil {
			return err
		}
	}

	capture, err := netcapture.Open(device)
	if err != nil {
		return fmt.Errorf("failed to open interface %q: %w", device, err)
	}
	defer capture.Close()

	aggregator := netstats.NewAggregator()
	application := app.New(device, capture, aggregator)
	application.Run()

	return nil
}

func parseFlags(args []string) (Config, error) {
	fs := flag.NewFlagSet("net-monitor", flag.ContinueOnError)

	var cfg Config

	fs.BoolVar(&cfg.ShowInterfaces, "si", false, "show available network interfaces")
	fs.StringVar(&cfg.Device, "i", "", "network interface to capture")
	fs.StringVar(&cfg.Device, "interface", "", "network interface to capture")

	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func printInterfaces(stdout io.Writer) error {
	devices, err := netdevice.ListDevices()
	if err != nil {
		return fmt.Errorf("failed to list interfaces: %w", err)
	}
	for _, device := range devices {
		fmt.Fprint(stdout, device.Name)
		if device.Description != "" {
			fmt.Fprintf(stdout, " - %s", device.Description)
		}
		if len(device.Addresses) > 0 {
			fmt.Fprintf(stdout, " [%s]", strings.Join(device.Addresses, ", "))
		}
		fmt.Fprintln(stdout)
	}
	return nil
}
