package tui

import (
	"fmt"
	"strings"

	"github.com/AvengeMedia/danklinux/internal/distros"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) viewGentooUseFlags() string {
	var b strings.Builder

	b.WriteString(m.renderBanner())
	b.WriteString("\n")

	title := m.styles.Title.Render("Gentoo Global USE Flags")
	b.WriteString(title)
	b.WriteString("\n\n")

	info := m.styles.Normal.Render("The following global USE flags will be enabled in /etc/portage/make.conf:")
	b.WriteString(info)
	b.WriteString("\n\n")

	for _, flag := range distros.GentooGlobalUseFlags {
		flagLine := m.styles.Success.Render(fmt.Sprintf("  • %s", flag))
		b.WriteString(flagLine)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	note := m.styles.Subtle.Render("These flags ensure proper Qt6, Wayland, and compositor support.")
	b.WriteString(note)
	b.WriteString("\n\n")

	help := m.styles.Subtle.Render("Press Enter to continue, Esc to go back")
	b.WriteString(help)

	return b.String()
}

func (m Model) updateGentooUseFlagsState(msg tea.Msg) (tea.Model, tea.Cmd) {
	if gccMsg, ok := msg.(gccVersionCheckMsg); ok {
		if gccMsg.err != nil || gccMsg.major < 15 {
			m.state = StateGentooGCCCheck
			return m, nil
		}
		if checkFingerprintEnabled() {
			m.state = StateAuthMethodChoice
			m.selectedConfig = 0
		} else {
			m.state = StatePasswordPrompt
			m.passwordInput.Focus()
		}
		return m, nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "enter":
			if m.selectedWM == 1 {
				return m, m.checkGCCVersion()
			}
			if checkFingerprintEnabled() {
				m.state = StateAuthMethodChoice
				m.selectedConfig = 0
			} else {
				m.state = StatePasswordPrompt
				m.passwordInput.Focus()
			}
			return m, nil
		case "esc":
			m.state = StateDependencyReview
			return m, nil
		}
	}
	return m, m.listenForLogs()
}
