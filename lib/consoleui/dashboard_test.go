package consoleui

import (
	"strconv"
	"testing"
	"time"

	"github.com/ikondratev/net-monitor/lib/netinterface"
)

func TestToTableRowsLimitsRowsAndFormatsValues(t *testing.T) {
	rows := make([]netinterface.ConnRow, 0, maxTableRows+1)
	for i := 0; i < maxTableRows+1; i++ {
		rows = append(rows, netinterface.ConnRow{
			Connection: netinterface.Connection{
				SrcIP:   "192.0.2." + strconv.Itoa(i+1),
				SrcPort: "1000",
				DstIP:   "198.51.100.1",
				DstPort: "443",
				Proto:   "TCP",
			},
			ConnStats: netinterface.ConnStats{
				PacketCount: i + 1,
				TotalBytes:  1536,
			},
		})
	}

	tableRows := toTableRows(rows)
	if len(tableRows) != maxTableRows {
		t.Fatalf("expected %d table rows, got %d", maxTableRows, len(tableRows))
	}

	first := tableRows[0]
	if first[0] != "TCP" {
		t.Fatalf("expected protocol TCP, got %q", first[0])
	}
	if first[1] != "192.0.2.1:1000" {
		t.Fatalf("unexpected source column: %q", first[1])
	}
	if first[2] != "198.51.100.1:443" {
		t.Fatalf("unexpected destination column: %q", first[2])
	}
	if first[3] != "1" {
		t.Fatalf("unexpected packet count column: %q", first[3])
	}
	if first[4] != "1.5 KB" {
		t.Fatalf("unexpected volume column: %q", first[4])
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want string
	}{
		{name: "bytes", in: 512, want: "512 B"},
		{name: "kilobytes", in: 1536, want: "1.5 KB"},
		{name: "megabytes", in: 2 * 1024 * 1024, want: "2.0 MB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatBytes(tt.in); got != tt.want {
				t.Fatalf("formatBytes(%d) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestSecondsUntilNextRefreshReturnsZeroAfterInterval(t *testing.T) {
	lastRefresh := time.Now().Add(-refreshInterval)

	if got := secondsUntilNextRefresh(lastRefresh); got != 0 {
		t.Fatalf("expected 0 seconds until next refresh, got %d", got)
	}
}
