package setup

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"sub-muse/internal/config"
	"sub-muse/internal/keyring"
	"sub-muse/internal/subsonic"
)

type Step int

const (
	StepWelcome Step = iota
	StepServerURL
	StepUsername
	StepPassword
	StepTesting
	StepSaving
	StepDone
	StepError
)

type connectionTestResult struct {
	Success bool
	Error   error
}

type saveResult struct {
	Success bool
	Error   error
}

type Model struct {
	step       Step
	serverURL  textinput.Model
	username   textinput.Model
	password   textinput.Model
	spinner    spinner.Model
	errorMsg   string
	result     *SetupResult
	serverURLs string
	focused    int
}

func NewModel() Model {
	serverURL := textinput.New()
	serverURL.Placeholder = "http://localhost:4040"
	serverURL.Focus()
	serverURL.CharLimit = 200

	username := textinput.New()
	username.Placeholder = "Enter your username"
	username.CharLimit = 100

	password := textinput.New()
	password.Placeholder = "Enter your password"
	password.EchoMode = textinput.EchoPassword
	password.CharLimit = 100

	return Model{
		step:      StepWelcome,
		serverURL: serverURL,
		username:  username,
		password:  password,
		spinner:   spinner.New(spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205")))),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			switch m.step {
			case StepWelcome:
				return m, tea.Quit
			case StepServerURL:
				m.step = StepWelcome
			case StepUsername:
				m.step = StepServerURL
				m.focused = 1
				m.serverURL.Focus()
			case StepPassword:
				m.step = StepUsername
				m.focused = 2
				m.username.Focus()
			case StepError:
				m.step = StepPassword
				m.focused = 3
				m.password.Focus()
			}
		case tea.KeyEnter:
			mm, cmd := m.handleEnter()
			return mm, cmd
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case connectionTestResult:
		if msg.Success {
			m.step = StepSaving
			return m, m.saveConfig()
		}
		m.step = StepError
		m.errorMsg = msg.Error.Error()
	case saveResult:
		if msg.Success {
			m.step = StepDone
		} else {
			m.step = StepError
			m.errorMsg = "Failed to save configuration: " + msg.Error.Error()
		}
	}

	var cmd tea.Cmd
	switch m.step {
	case StepServerURL:
		m.serverURL, cmd = m.serverURL.Update(msg)
	case StepUsername:
		m.username, cmd = m.username.Update(msg)
	case StepPassword:
		m.password, cmd = m.password.Update(msg)
	}

	return m, cmd
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case StepWelcome:
		m.step = StepServerURL
		m.serverURL.Focus()
case StepServerURL:
		m.step = StepUsername
		m.focused = 2
		m.username.Focus()
	case StepUsername:
		m.step = StepPassword
		m.focused = 3
		m.password.Focus()
	case StepPassword:
		m.step = StepTesting
		return m, m.testConnection()
	case StepError:
		m.step = StepPassword
		m.password.Focus()
	case StepDone:
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) testConnection() tea.Cmd {
	return func() tea.Msg {
		client := subsonic.NewClient(
			m.serverURL.Value(),
			m.username.Value(),
			m.password.Value(),
			"sub-muse",
		)
		err := client.Ping()
		return connectionTestResult{
			Success: err == nil,
			Error:   err,
		}
	}
}

func (m Model) saveConfig() tea.Cmd {
	return func() tea.Msg {
		cfg := &config.Config{
			Configured: true,
			ServerURL:  m.serverURL.Value(),
			Username:   m.username.Value(),
			ClientName: "sub-muse",
		}

		err := config.SaveConfig(cfg)
		if err != nil {
			return saveResult{Success: false, Error: err}
		}

		err = keyring.SavePassword(m.username.Value(), m.password.Value())
		if err != nil {
			return saveResult{Success: false, Error: err}
		}

		return saveResult{Success: true, Error: nil}
	}
}

func (m Model) View() string {
	switch m.step {
	case StepWelcome:
		return m.viewWelcome()
	case StepServerURL:
		return m.viewInput("Server URL", m.serverURL)
	case StepUsername:
		return m.viewInput("Username", m.username)
	case StepPassword:
		return m.viewInput("Password", m.password)
	case StepTesting:
		return m.viewTesting()
	case StepSaving:
		return m.viewSaving()
	case StepDone:
		return m.viewDone()
	case StepError:
		return m.viewError()
	}
	return ""
}

func (m Model) viewWelcome() string {
	title := "Welcome to sub-muse"
	content := "\nThis looks like your first time running sub-muse.\n\nPress Enter to set up your Subsonic server credentials, or Ctrl+C to quit.\n"
	buttons := "\n[ Enter to continue ]"
	return lipgloss.JoinVertical(lipgloss.Center,
		lipgloss.NewStyle().Width(60).Align(lipgloss.Center).Render(title),
		lipgloss.NewStyle().Width(60).Align(lipgloss.Center).Render(content),
		buttons,
	)
}

func (m Model) viewInput(label string, input textinput.Model) string {
	errorMsg := ""
	if m.step == StepError && m.focused == 1 {
		errorMsg = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5370")).Render("✗ "+m.errorMsg)
	}
	if m.step == StepError && m.focused == 2 {
		errorMsg = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5370")).Render("✗ "+m.errorMsg)
	}
	if m.step == StepError && m.focused == 3 {
		errorMsg = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5370")).Render("✗ "+m.errorMsg)
	}

	return lipgloss.JoinVertical(lipgloss.Center,
		lipgloss.NewStyle().Width(60).Align(lipgloss.Center).Render(label),
		input.View(),
		errorMsg,
		"\n[ Enter to continue ]",
	)
}

func (m Model) viewTesting() string {
	return lipgloss.JoinVertical(lipgloss.Center,
		lipgloss.NewStyle().Width(60).Align(lipgloss.Center).Render("Testing connection..."),
		m.spinner.View(),
		"\nPlease wait while we verify your credentials.",
	)
}

func (m Model) viewSaving() string {
	return lipgloss.JoinVertical(lipgloss.Center,
		lipgloss.NewStyle().Width(60).Align(lipgloss.Center).Render("Saving configuration..."),
		m.spinner.View(),
		"\nYour credentials are being securely stored.",
	)
}

func (m Model) viewDone() string {
	title := "Setup Complete!"
	content := "\nYou're all set! sub-muse is now configured and ready to use.\n\nPress Enter to start listening to music."
	return lipgloss.JoinVertical(lipgloss.Center,
		lipgloss.NewStyle().Width(60).Align(lipgloss.Center).Foreground(lipgloss.Color("#00C853")).Render(title),
		lipgloss.NewStyle().Width(60).Align(lipgloss.Center).Render(content),
		"\n[ Enter to continue ]",
	)
}

func (m Model) viewError() string {
	return lipgloss.JoinVertical(lipgloss.Center,
		lipgloss.NewStyle().Width(60).Align(lipgloss.Center).Foreground(lipgloss.Color("#FF5370")).Render("Error"),
		lipgloss.NewStyle().Width(60).Align(lipgloss.Center).Render(m.errorMsg),
		"\n[ Esc to go back ]",
	)
}

type SetupResult struct {
	Configured bool
}
