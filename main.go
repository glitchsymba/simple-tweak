package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ══════════════════════════════════════════════════════════════════════════════
//  STYLES
// ══════════════════════════════════════════════════════════════════════════════

var (
	green    = lipgloss.Color("#00FF9C")
	dimGreen = lipgloss.Color("#00AA66")
	blue     = lipgloss.Color("#44AAFF")
	yellow   = lipgloss.Color("#FFD700")
	red      = lipgloss.Color("#FF4C4C")
	gray     = lipgloss.Color("#555555")
	bg       = lipgloss.Color("#0D0D0D")
	white    = lipgloss.Color("#EEEEEE")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(green).
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(green).
			Padding(0, 3)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(gray).
			Italic(true)

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(green)

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(red)

	warnStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(yellow)

	infoStyle = lipgloss.NewStyle().
			Foreground(blue)

	dimStyle = lipgloss.NewStyle().
			Foreground(gray)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gray).
			Padding(1, 2)

	confirmBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(yellow).
			Padding(1, 3)

	keyStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	catStyle = lipgloss.NewStyle().
			Foreground(blue).
			Bold(true)

	appStyle = lipgloss.NewStyle().
			Padding(1, 3)
)

// ══════════════════════════════════════════════════════════════════════════════
//  TWEAK
// ══════════════════════════════════════════════════════════════════════════════

type tweak struct {
	title       string
	description string
	category    string
	risk        string // "safe" | "moderate" | "advanced"
	applyFn     func() error
	revertFn    func() error
}

func (t tweak) Title() string { return t.title }
func (t tweak) Description() string {
	var riskLabel string
	switch t.risk {
	case "safe":
		riskLabel = lipgloss.NewStyle().Foreground(green).Render("[safe]")
	case "moderate":
		riskLabel = lipgloss.NewStyle().Foreground(yellow).Render("[moderate]")
	case "advanced":
		riskLabel = lipgloss.NewStyle().Foreground(red).Render("[advanced]")
	}
	return fmt.Sprintf("[%s]  %s  %s", t.category, t.description, riskLabel)
}
func (t tweak) FilterValue() string { return t.title + " " + t.category + " " + t.risk }

// ══════════════════════════════════════════════════════════════════════════════
//  TWEAKS
// ══════════════════════════════════════════════════════════════════════════════

func buildTweaks() []tweak {
	return []tweak{

		// ── PERFORMANCE ───────────────────────────────────────────────────────
		{
			title:       "⚡  High Performance Power Plan",
			description: "Switch to High Performance for maximum CPU speed",
			category:    "Performance", risk: "safe",
			applyFn:  func() error { return ps(`powercfg /setactive 8c5e7fda-e8bf-4a96-9a85-a6e23a8c635c`) },
			revertFn: func() error { return ps(`powercfg /setactive 381b4222-f694-41f0-9685-ff5bb260df2e`) },
		},
		{
			title:       "🚀  Ultimate Performance Power Plan",
			description: "Enable hidden Ultimate Performance plan – best for desktops",
			category:    "Performance", risk: "moderate",
			applyFn: func() error {
				_ = ps(`powercfg -duplicatescheme e9a42b02-d5df-448d-aa00-03f14749eb61`)
				return ps(`powercfg /setactive e9a42b02-d5df-448d-aa00-03f14749eb61`)
			},
		},
		{
			title:       "🔧  Disable SysMain (Superfetch)",
			description: "Stop pre-loading apps into RAM – helps HDDs & low-RAM PCs",
			category:    "Performance", risk: "moderate",
			applyFn:  func() error { return ps(`Stop-Service SysMain -Force; Set-Service SysMain -StartupType Disabled`) },
			revertFn: func() error { return ps(`Set-Service SysMain -StartupType Automatic; Start-Service SysMain`) },
		},
		{
			title:       "🖥️  Disable Visual Effects",
			description: "Turn off animations & shadows for a snappier UI",
			category:    "Performance", risk: "safe",
			applyFn: func() error {
				return reg(`add "HKCU\Software\Microsoft\Windows\CurrentVersion\Explorer\VisualEffects" /v VisualFXSetting /t REG_DWORD /d 2 /f`)
			},
			revertFn: func() error {
				return reg(`add "HKCU\Software\Microsoft\Windows\CurrentVersion\Explorer\VisualEffects" /v VisualFXSetting /t REG_DWORD /d 1 /f`)
			},
		},
		{
			title:       "🧠  Prioritize Foreground Apps",
			description: "Tell the CPU scheduler to favor the active window",
			category:    "Performance", risk: "safe",
			applyFn: func() error {
				return reg(`add "HKLM\SYSTEM\CurrentControlSet\Control\PriorityControl" /v Win32PrioritySeparation /t REG_DWORD /d 38 /f`)
			},
		},
		{
			title:       "⏱️  Disable Startup Delay",
			description: "Remove the 10-second delay before startup apps launch",
			category:    "Performance", risk: "safe",
			applyFn: func() error {
				return reg(`add "HKCU\Software\Microsoft\Windows\CurrentVersion\Explorer\Serialize" /v StartupDelayInMSec /t REG_DWORD /d 0 /f`)
			},
		},
		{
			title:       "🗑️  Disable Hibernation",
			description: "Delete hiberfil.sys and reclaim several GB of disk space",
			category:    "Performance", risk: "moderate",
			applyFn:  func() error { return ps(`powercfg /hibernate off`) },
			revertFn: func() error { return ps(`powercfg /hibernate on`) },
		},

		// ── GAMING ────────────────────────────────────────────────────────────
		{
			title:       "🎮  Enable Game Mode",
			description: "Prioritize GPU/CPU resources when a game is running",
			category:    "Gaming", risk: "safe",
			applyFn: func() error {
				return reg(`add "HKCU\Software\Microsoft\GameBar" /v AllowAutoGameMode /t REG_DWORD /d 1 /f`)
			},
		},
		{
			title:       "🖱️  Disable Mouse Acceleration",
			description: "Enable true 1:1 mouse input – essential for FPS games",
			category:    "Gaming", risk: "safe",
			applyFn: func() error {
				_ = reg(`add "HKCU\Control Panel\Mouse" /v MouseSpeed /t REG_SZ /d 0 /f`)
				_ = reg(`add "HKCU\Control Panel\Mouse" /v MouseThreshold1 /t REG_SZ /d 0 /f`)
				return reg(`add "HKCU\Control Panel\Mouse" /v MouseThreshold2 /t REG_SZ /d 0 /f`)
			},
		},
		{
			title:       "📊  Enable HW-Accelerated GPU Scheduling",
			description: "Reduce GPU latency (Windows 10 2004+ and a modern GPU required)",
			category:    "Gaming", risk: "moderate",
			applyFn: func() error {
				return reg(`add "HKLM\SYSTEM\CurrentControlSet\Control\GraphicsDrivers" /v HwSchMode /t REG_DWORD /d 2 /f`)
			},
			revertFn: func() error {
				return reg(`add "HKLM\SYSTEM\CurrentControlSet\Control\GraphicsDrivers" /v HwSchMode /t REG_DWORD /d 1 /f`)
			},
		},
		{
			title:       "🔔  Disable Xbox Game Bar",
			description: "Stop Game Bar from running in the background to save RAM",
			category:    "Gaming", risk: "safe",
			applyFn: func() error {
				return reg(`add "HKCU\Software\Microsoft\Windows\CurrentVersion\GameDVR" /v AppCaptureEnabled /t REG_DWORD /d 0 /f`)
			},
		},
		{
			title:       "🕹️  Set Platform Timer Resolution",
			description: "Use 0.5ms timer for reduced frame latency in games",
			category:    "Gaming", risk: "moderate",
			applyFn:  func() error { return ps(`bcdedit /set useplatformtick yes`) },
			revertFn: func() error { return ps(`bcdedit /deletevalue useplatformtick`) },
		},

		// ── NETWORK ───────────────────────────────────────────────────────────
		{
			title:       "🌐  Flush DNS Cache",
			description: "Clear DNS resolver cache – fixes many connection problems",
			category:    "Network", risk: "safe",
			applyFn: func() error { return ps(`ipconfig /flushdns`) },
		},
		{
			title:       "📡  Set DNS to Cloudflare (1.1.1.1)",
			description: "Use fast & private Cloudflare DNS on all active adapters",
			category:    "Network", risk: "moderate",
			applyFn: func() error {
				return ps(`Get-NetAdapter | Where-Object {$_.Status -eq 'Up'} | ForEach-Object { Set-DnsClientServerAddress -InterfaceIndex $_.ifIndex -ServerAddresses ("1.1.1.1","1.0.0.1") }`)
			},
			revertFn: func() error {
				return ps(`Get-NetAdapter | Where-Object {$_.Status -eq 'Up'} | ForEach-Object { Set-DnsClientServerAddress -InterfaceIndex $_.ifIndex -ResetServerAddresses }`)
			},
		},
		{
			title:       "⚡  Disable Nagle's Algorithm",
			description: "Reduce TCP latency – great for competitive online gaming",
			category:    "Network", risk: "moderate",
			applyFn: func() error {
				_ = reg(`add "HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces" /v TcpAckFrequency /t REG_DWORD /d 1 /f`)
				return reg(`add "HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces" /v TCPNoDelay /t REG_DWORD /d 1 /f`)
			},
		},
		{
			title:       "🔒  Disable IPv6",
			description: "Turn off IPv6 on all adapters – can fix VPN & connection issues",
			category:    "Network", risk: "moderate",
			applyFn:  func() error { return ps(`Disable-NetAdapterBinding -Name "*" -ComponentID ms_tcpip6`) },
			revertFn: func() error { return ps(`Enable-NetAdapterBinding -Name "*" -ComponentID ms_tcpip6`) },
		},

		// ── PRIVACY ───────────────────────────────────────────────────────────
		{
			title:       "🔇  Disable Windows Telemetry",
			description: "Stop Microsoft diagnostic data collection in the background",
			category:    "Privacy", risk: "moderate",
			applyFn: func() error {
				_ = reg(`add "HKLM\SOFTWARE\Policies\Microsoft\Windows\DataCollection" /v AllowTelemetry /t REG_DWORD /d 0 /f`)
				return ps(`Stop-Service DiagTrack -Force; Set-Service DiagTrack -StartupType Disabled`)
			},
		},
		{
			title:       "📢  Disable Advertising ID",
			description: "Stop apps from using your personalized advertising ID",
			category:    "Privacy", risk: "safe",
			applyFn: func() error {
				return reg(`add "HKCU\Software\Microsoft\Windows\CurrentVersion\AdvertisingInfo" /v Enabled /t REG_DWORD /d 0 /f`)
			},
		},
		{
			title:       "🔍  Disable Cortana & Web Search",
			description: "Remove Cortana and stop web results appearing in Start search",
			category:    "Privacy", risk: "moderate",
			applyFn: func() error {
				_ = reg(`add "HKLM\SOFTWARE\Policies\Microsoft\Windows\Windows Search" /v AllowCortana /t REG_DWORD /d 0 /f`)
				return reg(`add "HKLM\SOFTWARE\Policies\Microsoft\Windows\Windows Search" /v DisableWebSearch /t REG_DWORD /d 1 /f`)
			},
		},
		{
			title:       "📷  Disable Activity History",
			description: "Stop Windows from logging the apps and files you open",
			category:    "Privacy", risk: "safe",
			applyFn: func() error {
				return reg(`add "HKLM\SOFTWARE\Policies\Microsoft\Windows\System" /v EnableActivityFeed /t REG_DWORD /d 0 /f`)
			},
		},

		// ── CLEANUP ───────────────────────────────────────────────────────────
		{
			title:       "🧹  Clear Temp Files",
			description: "Delete all files in %TEMP% and C:\\Windows\\Temp",
			category:    "Cleanup", risk: "safe",
			applyFn: func() error {
				tmp := os.Getenv("TEMP")
				_ = ps(fmt.Sprintf(`Remove-Item -Path "%s\*" -Recurse -Force -ErrorAction SilentlyContinue`, tmp))
				return ps(`Remove-Item -Path "C:\Windows\Temp\*" -Recurse -Force -ErrorAction SilentlyContinue`)
			},
		},
		{
			title:       "🗑️  Empty Recycle Bin",
			description: "Permanently remove all items sitting in the Recycle Bin",
			category:    "Cleanup", risk: "safe",
			applyFn: func() error { return ps(`Clear-RecycleBin -Force -ErrorAction SilentlyContinue`) },
		},
		{
			title:       "📦  Clean Windows Update Cache",
			description: "Remove old update download files to reclaim disk space",
			category:    "Cleanup", risk: "safe",
			applyFn: func() error {
				_ = ps(`Stop-Service wuauserv -Force`)
				_ = ps(`Remove-Item -Path "C:\Windows\SoftwareDistribution\Download\*" -Recurse -Force -ErrorAction SilentlyContinue`)
				return ps(`Start-Service wuauserv`)
			},
		},

		// ── SECURITY ──────────────────────────────────────────────────────────
		{
			title:       "🛡️  Enable Defender Real-Time Protection",
			description: "Make sure Windows Defender is actively scanning files",
			category:    "Security", risk: "safe",
			applyFn: func() error { return ps(`Set-MpPreference -DisableRealtimeMonitoring $false`) },
		},
		{
			title:       "🔥  Enable Firewall (All Profiles)",
			description: "Turn on the firewall for Domain, Private and Public networks",
			category:    "Security", risk: "safe",
			applyFn: func() error { return ps(`Set-NetFirewallProfile -Profile Domain,Public,Private -Enabled True`) },
		},
		{
			title:       "🔑  Enable Credential Guard",
			description: "Protect login credentials with virtualization-based security",
			category:    "Security", risk: "advanced",
			applyFn: func() error {
				return reg(`add "HKLM\SYSTEM\CurrentControlSet\Control\DeviceGuard" /v EnableVirtualizationBasedSecurity /t REG_DWORD /d 1 /f`)
			},
		},
	}
}

// ══════════════════════════════════════════════════════════════════════════════
//  HELPERS
// ══════════════════════════════════════════════════════════════════════════════

func ps(command string) error {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func reg(args string) error {
	parts := strings.Fields(args)
	cmd := exec.Command("reg", parts...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ══════════════════════════════════════════════════════════════════════════════
//  MESSAGES
// ══════════════════════════════════════════════════════════════════════════════

type applyDoneMsg struct{ err error }
type revertDoneMsg struct{ err error }

// ══════════════════════════════════════════════════════════════════════════════
//  STATE
// ══════════════════════════════════════════════════════════════════════════════

type screen int

const (
	screenMenu screen = iota
	screenConfirm
	screenApplying
	screenResult
)

// ══════════════════════════════════════════════════════════════════════════════
//  MODEL
// ══════════════════════════════════════════════════════════════════════════════

type model struct {
	list      list.Model
	spinner   spinner.Model
	screen    screen
	selected  tweak
	result    string
	isError   bool
	isRevert  bool
	confirmYN bool
	width     int
	height    int
}

// ══════════════════════════════════════════════════════════════════════════════
//  INIT
// ══════════════════════════════════════════════════════════════════════════════

func initialModel() model {
	tweaks := buildTweaks()
	items := make([]list.Item, len(tweaks))
	for i, t := range tweaks {
		items[i] = t
	}

	del := list.NewDefaultDelegate()
	del.ShowDescription = true
	del.Styles.SelectedTitle = del.Styles.SelectedTitle.
		Foreground(green).BorderForeground(green).Bold(true)
	del.Styles.SelectedDesc = del.Styles.SelectedDesc.
		Foreground(dimGreen).BorderForeground(green)
	del.Styles.NormalTitle = del.Styles.NormalTitle.Foreground(white)
	del.Styles.NormalDesc = del.Styles.NormalDesc.Foreground(gray)

	l := list.New(items, del, 72, 22)
	l.Title = "WINDOWS PC TWEAKER"
	l.Styles.Title = titleStyle
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(green)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(green)
	l.Styles.StatusBar = lipgloss.NewStyle().Foreground(gray)

	sp := spinner.New()
	sp.Spinner = spinner.Points
	sp.Style = lipgloss.NewStyle().Foreground(green).Bold(true)

	return model{list: l, spinner: sp, screen: screenMenu}
}

func (m model) Init() tea.Cmd { return nil }

// ══════════════════════════════════════════════════════════════════════════════
//  UPDATE
// ══════════════════════════════════════════════════════════════════════════════

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.list.SetWidth(msg.Width - 6)
		m.list.SetHeight(msg.Height - 12)
		return m, nil

	case tea.KeyMsg:
		switch m.screen {
		case screenMenu:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter":
				if sel, ok := m.list.SelectedItem().(tweak); ok {
					m.selected = sel
					m.confirmYN = false
					m.screen = screenConfirm
					return m, nil
				}
			}

		case screenConfirm:
			switch msg.String() {
			case "left", "h", "right", "l", "tab":
				m.confirmYN = !m.confirmYN
			case "y", "Y":
				m.confirmYN = true
			case "n", "N":
				m.confirmYN = false
			case "enter":
				if m.confirmYN {
					m.isRevert = false
					m.screen = screenApplying
					sel := m.selected
					return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
						return applyDoneMsg{err: sel.applyFn()}
					})
				}
				m.screen = screenMenu
			case "r", "R":
				if m.selected.revertFn != nil {
					m.isRevert = true
					m.screen = screenApplying
					sel := m.selected
					return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
						return revertDoneMsg{err: sel.revertFn()}
					})
				}
			case "esc", "q":
				m.screen = screenMenu
			}

		case screenResult:
			switch msg.String() {
			case "enter", "esc", "q", " ":
				m.screen = screenMenu
				m.result = ""
			}
		}

	case applyDoneMsg:
		m.screen = screenResult
		m.isError = msg.err != nil
		if msg.err != nil {
			m.result = "Failed: " + msg.err.Error()
		} else {
			m.result = "Tweak applied successfully!"
		}
		return m, nil

	case revertDoneMsg:
		m.screen = screenResult
		m.isError = msg.err != nil
		if msg.err != nil {
			m.result = "Revert failed: " + msg.err.Error()
		} else {
			m.result = "Tweak reverted successfully!"
		}
		return m, nil

	case spinner.TickMsg:
		if m.screen == screenApplying {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	if m.screen == screenMenu {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
	return m, nil
}

// ══════════════════════════════════════════════════════════════════════════════
//  VIEW
// ══════════════════════════════════════════════════════════════════════════════

func (m model) View() string {
	header := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("  ⚙  SYMBA TWEAKS  "),
		subtitleStyle.Render("   Performance · Gaming · Privacy · Network · Security"),
		"",
	)

	var body string

	switch m.screen {

	case screenMenu:
		legend := strings.Join([]string{
			"  ",
			keyStyle.Render("[↑↓]") + dimStyle.Render(" move"),
			"  " + keyStyle.Render("[enter]") + dimStyle.Render(" select"),
			"  " + keyStyle.Render("[/]") + dimStyle.Render(" filter"),
			"  " + keyStyle.Render("[q]") + dimStyle.Render(" quit"),
		}, "")
		riskLegend := "  " +
			lipgloss.NewStyle().Foreground(green).Render("● safe") + "   " +
			lipgloss.NewStyle().Foreground(yellow).Render("● moderate") + "   " +
			lipgloss.NewStyle().Foreground(red).Render("● advanced")

		body = lipgloss.JoinVertical(lipgloss.Left,
			m.list.View(), "", legend, dimStyle.Render(riskLegend))

	case screenConfirm:
		t := m.selected
		riskColor := green
		switch t.risk {
		case "moderate":
			riskColor = yellow
		case "advanced":
			riskColor = red
		}

		var yesBtn, noBtn string
		if m.confirmYN {
			yesBtn = lipgloss.NewStyle().Foreground(bg).Background(green).Bold(true).Padding(0, 3).Render("YES – Apply")
			noBtn = lipgloss.NewStyle().Foreground(gray).Border(lipgloss.RoundedBorder()).BorderForeground(gray).Padding(0, 2).Render("NO – Cancel")
		} else {
			yesBtn = lipgloss.NewStyle().Foreground(green).Border(lipgloss.RoundedBorder()).BorderForeground(green).Padding(0, 2).Render("YES – Apply")
			noBtn = lipgloss.NewStyle().Foreground(bg).Background(red).Bold(true).Padding(0, 3).Render("NO – Cancel")
		}

		revertLine := ""
		if t.revertFn != nil {
			revertLine = "\n" + dimStyle.Render("  [r] revert this tweak instead")
		}

		box := confirmBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
			catStyle.Render("Category: "+t.category)+"   "+
				lipgloss.NewStyle().Foreground(riskColor).Bold(true).Render("Risk: "+t.risk),
			"",
			warnStyle.Render(t.title),
			infoStyle.Render(t.description),
			"",
			dimStyle.Render("Apply this tweak now?"),
			"",
			lipgloss.JoinHorizontal(lipgloss.Top, yesBtn, "   ", noBtn),
			revertLine,
		))
		body = lipgloss.JoinVertical(lipgloss.Left, box, "",
			dimStyle.Render("  [←/→ or tab] switch  [enter] confirm  [esc] back"))

	case screenApplying:
		action := "Applying"
		if m.isRevert {
			action = "Reverting"
		}
		box := boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
			m.spinner.View()+"  "+infoStyle.Render(action+": "+m.selected.title),
			"",
			dimStyle.Render("Please wait, this may take a moment..."),
		))
		body = box

	case screenResult:
		icon, st := "✓  ", successStyle
		if m.isError {
			icon, st = "✗  ", errorStyle
		}
		box := boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
			st.Render(icon+m.result),
			"",
			dimStyle.Render("[enter / esc] back to menu"),
		))
		body = box
	}

	return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left, header, body))
}

// ══════════════════════════════════════════════════════════════════════════════
//  MAIN
// ══════════════════════════════════════════════════════════════════════════════

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Fehler: %v\n", err)
		os.Exit(1)
	}
}
