# Net Monitor

Simple terminal network monitor. It captures packets from an active network interface and shows aggregated connection statistics in the console.

## Requirements

- Go 1.24.2 or newer
- libpcap on macOS/Linux, or Npcap on Windows
- Permissions to capture network packets

On macOS/Linux, packet capture usually requires `sudo`.

## Install

```bash
go install github.com/ikondratev/net-monitor/cmd/net-monitor@latest
```

Make sure Go's bin directory is in your `PATH`:

```bash
export PATH="$HOME/go/bin:$PATH"
```

Run:

```bash
sudo net-monitor
```

List available interfaces:

```bash
net-monitor -si
```

Capture from a specific interface:

```bash
sudo net-monitor -i en0
```

## Build Locally

```bash
go build -o bin/net-monitor ./cmd/net-monitor
```

Run:

```bash
sudo ./bin/net-monitor
```