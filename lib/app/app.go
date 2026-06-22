package app

import (
	"github.com/ikondratev/net-monitor/lib/consoleui"
	"github.com/ikondratev/net-monitor/lib/netcapture"
	"github.com/ikondratev/net-monitor/lib/netstats"
)

type App struct {
	device     string
	capture    *netcapture.Capture
	aggregator *netstats.Aggregator
}

func New(device string, capture *netcapture.Capture, aggregator *netstats.Aggregator) *App {
	return &App{
		device:     device,
		capture:    capture,
		aggregator: aggregator,
	}
}

func (a *App) Run() {
	go func() {
		for packet := range a.capture.Packets() {
			a.aggregator.ProcessPacket(packet)
		}
	}()

	consoleui.RunDashboard(a.device, a.aggregator)
}
