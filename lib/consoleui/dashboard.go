package consoleui

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ikondratev/net-monitor/lib/netinterface"
	"github.com/ikondratev/net-monitor/lib/netstats"
)

const (
	maxTableRows    = 15
	refreshInterval = 5 * time.Second
)

var (
	titleStyle           = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	metaStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	helpStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	_          tea.Model = dashboardModel{}
)

type dashboardModel struct {
	device              string
	aggregator          *netstats.Aggregator
	spinner             spinner.Model
	spinerSecond        spinner.Model
	table               table.Model
	lastDataRefresh     time.Time
	totalRows           int
	isWaitingForTraffic bool
}

type refreshDataMsg struct{}

func newDashboardModel(device string, aggregator *netstats.Aggregator) dashboardModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))

	s2 := spinner.New()
	s2.Spinner = spinner.Points
	s2.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))

	styles := table.DefaultStyles()
	styles.Header = styles.Header.Bold(true).Foreground(lipgloss.Color("212"))
	styles.Cell = styles.Cell.Foreground(lipgloss.Color("252"))

	t := table.New(
		table.WithColumns([]table.Column{
			{Title: "PROTO", Width: 8},
			{Title: "SOURCE (IP:PORT)", Width: 32},
			{Title: "DESTINATION (IP:PORT)", Width: 32},
			{Title: "PACKETS", Width: 10},
			{Title: "VOLUME", Width: 12},
		}),
		table.WithFocused(false),
		table.WithHeight(maxTableRows),
	)

	rows := aggregator.ConnectionRows()
	t.SetRows(toTableRows(rows))
	t.SetStyles(styles)
	return dashboardModel{
		device:              device,
		aggregator:          aggregator,
		spinner:             s,
		spinerSecond:        s2,
		table:               t,
		lastDataRefresh:     time.Now(),
		totalRows:           len(rows),
		isWaitingForTraffic: len(rows) == 0,
	}
}

func (m dashboardModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.spinerSecond.Tick,
		refreshData(),
	)
}

func (m dashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case refreshDataMsg:
		rows := m.aggregator.ConnectionRows()
		m.table.SetRows(toTableRows(rows))
		m.totalRows = len(rows)
		m.isWaitingForTraffic = len(rows) == 0
		m.lastDataRefresh = time.Now()
		return m, refreshData()
	case spinner.TickMsg:
		var cmd1 tea.Cmd
		var cmd2 tea.Cmd

		m.spinner, cmd1 = m.spinner.Update(msg)
		m.spinerSecond, cmd2 = m.spinerSecond.Update(msg)
		return m, tea.Batch(cmd1, cmd2)
	}
	return m, nil
}

func (m dashboardModel) View() string {
	hiddenRows := m.totalRows - maxTableRows
	if hiddenRows < 0 {
		hiddenRows = 0
	}
	title := titleStyle.Render(fmt.Sprintf("%s Network Monitor", m.spinner.View()))
	meta := metaStyle.Render(fmt.Sprintf(
		"Interface: %s | Last refresh: %s | Next refresh in: %ds | Hidden connections: %d",
		m.device,
		m.lastDataRefresh.Format("15:04:05"),
		secondsUntilNextRefresh(m.lastDataRefresh),
		hiddenRows,
	))
	content := m.table.View()
	if m.isWaitingForTraffic {
		content = startupBoxStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				titleStyle.Render(fmt.Sprintf("Interface: %s %s", m.device, m.spinerSecond.View())),
			),
		)
	}
	footer := helpStyle.Render("Press q or Ctrl+C to quit")
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		meta,
		"",
		content,
		"",
		footer,
	)
}

func RunDashboard(device string, aggregator *netstats.Aggregator) error {
	p := tea.NewProgram(
		newDashboardModel(device, aggregator),
		tea.WithAltScreen(),
	)
	_, err := p.Run()
	return err
}

func refreshData() tea.Cmd {
	return tea.Tick(refreshInterval, func(time.Time) tea.Msg {
		return refreshDataMsg{}
	})
}

func toTableRows(rows []netinterface.ConnRow) []table.Row {
	if len(rows) > maxTableRows {
		rows = rows[:maxTableRows]
	}

	tableRows := make([]table.Row, 0, len(rows))
	for _, row := range rows {
		src := fmt.Sprintf("%s:%s", row.SrcIP, row.SrcPort)
		dst := fmt.Sprintf("%s:%s", row.DstIP, row.DstPort)
		tableRows = append(tableRows, table.Row{
			row.Proto,
			src,
			dst,
			strconv.Itoa(row.PacketCount),
			formatBytes(row.TotalBytes),
		})
	}

	return tableRows
}

func secondsUntilNextRefresh(lastDataRefresh time.Time) int {
	remaining := refreshInterval - time.Since(lastDataRefresh)
	if remaining <= 0 {
		return 0
	}
	return int(remaining.Seconds()) + 1
}

func formatBytes(b int) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
