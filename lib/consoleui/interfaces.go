package consoleui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/ikondratev/net-monitor/lib/netdevice"
)

var (
	interfacesTitleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	interfacesHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	interfacesCellStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	interfacesHelpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

func PrintInterfaces(stdout io.Writer, devices []netdevice.Device) error {
	if _, err := fmt.Fprintln(stdout, interfacesTitleStyle.Render("Available Network Interfaces")); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(stdout); err != nil {
		return err
	}
	for _, line := range interfaceLines(devices) {
		if _, err := fmt.Fprintln(stdout, line); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(stdout); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(stdout, interfacesHelpStyle.Render("Run: net-monitor -i <interface>")); err != nil {
		return err
	}
	return nil
}

func interfaceLines(devices []netdevice.Device) []string {
	lines := make([]string, 0, len(devices)+1)
	lines = append(lines, lipgloss.JoinHorizontal(
		lipgloss.Left,
		interfacesHeaderStyle.Width(16).Render("NAME"),
		interfacesHeaderStyle.Width(36).Render("DESCRIPTION"),
		interfacesHeaderStyle.Render("ADDRESSES"),
	))
	for _, device := range devices {
		description := device.Description
		if description == "" {
			description = "-"
		}
		addresses := "-"
		if len(device.Addresses) > 0 {
			addresses = strings.Join(device.Addresses, ", ")
		}
		lines = append(lines, lipgloss.JoinHorizontal(
			lipgloss.Left,
			interfacesCellStyle.Width(16).Render(device.Name),
			interfacesCellStyle.Width(36).Render(description),
			interfacesCellStyle.Render(addresses),
		))
	}
	return lines
}
