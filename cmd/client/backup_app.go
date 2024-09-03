package main

//
//import (
//	"log"
//	"net/mail"
//	"os"
//	"strings"
//
//	"github.com/PaBah/GophKeeper/internal/client"
//	"github.com/charmbracelet/bubbles/spinner"
//	tea "github.com/charmbracelet/bubbletea"
//	"github.com/charmbracelet/huh"
//	"github.com/charmbracelet/lipgloss"
//)
//
//type State int
//
//type Styles struct {
//	Base,
//	HeaderText,
//	Status,
//	StatusHeader,
//	Highlight,
//	ErrorHeaderText,
//	Help lipgloss.Style
//}
//
//var (
//	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
//	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
//	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
//)
//
//const (
//	Initial         State = iota
//	SignIn          State = iota
//	SignUp          State = iota
//	CredentialsForm State = iota
//	CardForm        State = iota
//
//	FileLoad  State = iota
//	Dashboard State = iota
//)
//
//type Model struct {
//	client      *client.ClientService
//	clientError error
//	state       State
//	spinner     spinner.Model
//	lg          *lipgloss.Renderer
//	forms       []*huh.Form
//	styles      *Styles
//	width       int
//	authMethod  string
//	changedForm bool
//}
//
//func NewModel(state State) Model {
//	s := spinner.New()
//	s.Spinner = spinner.Dot
//	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
//
//	m := Model{spinner: s, state: state}
//	m.client = client.NewClientService(":3200")
//	m.client.TryToConnect()
//	m.lg = lipgloss.DefaultRenderer()
//	m.styles = NewStyles(m.lg)
//	m.changedForm = false
//
//	m.forms = make([]*huh.Form, CardForm+1)
//	m.forms[Initial] = huh.NewForm(
//		huh.NewGroup(
//			huh.NewSelect[string]().
//				Key("authMethod").
//				Options(huh.NewOptions("SignIn", "SignUp")...).
//				Title("Hi! Welcome to GophKeeper!").
//				Description("Please, choose authorisation method before continue."),
//		)).WithWidth(80).WithShowErrors(false)
//	m.forms[SignIn] = huh.NewForm(
//		huh.NewGroup(
//			huh.NewNote().Title("Enter yours auth credentials"),
//			huh.NewInput().
//				Key("email").
//				Validate(func(email string) error {
//					_, err := mail.ParseAddress(email)
//					if err != nil {
//						return err
//					}
//					return m.clientError
//				}).
//				Placeholder("Email"),
//			huh.NewInput().
//				Key("password").
//				EchoMode(huh.EchoModePassword).
//				Placeholder(
//					"Password"),
//		),
//	).WithWidth(123).WithShowErrors(false)
//	m.forms[SignUp] = huh.NewForm(
//		huh.NewGroup(
//			huh.NewNote().Title("Enter yours auth credentials"),
//			huh.NewInput().
//				Key("email").
//				Validate(func(email string) error {
//					_, err := mail.ParseAddress(email)
//					if err != nil {
//						return err
//					}
//					return m.clientError
//				}).
//				Placeholder("Email"),
//			huh.NewInput().
//				Key("password").
//				EchoMode(huh.EchoModePassword).
//				Placeholder(
//					"Password"),
//		),
//	).WithWidth(80).WithShowErrors(false)
//	m.forms[CredentialsForm] = huh.NewForm(
//		huh.NewGroup(
//			huh.NewSelect[string]().
//				Key("authMethod").
//				Options(huh.NewOptions("SignIn", "SignUp")...).
//				Title("Hi! Welcome to GophKeeper!").
//				Description("Please, choose authorisation method before continue."),
//		)).WithWidth(80).WithShowErrors(false)
//	m.forms[CardForm] = huh.NewForm(
//		huh.NewGroup(
//			huh.NewSelect[string]().
//				Key("authMethod").
//				Options(huh.NewOptions("SignIn", "SignUp")...).
//				Title("Hi! Welcome to GophKeeper!").
//				Description("Please, choose authorisation method before continue."),
//		)).WithWidth(80).WithShowErrors(false)
//	return m
//}
//
//func NewStyles(lg *lipgloss.Renderer) *Styles {
//	s := Styles{}
//	s.Base = lg.NewStyle().
//		Padding(1, 4, 0, 1)
//	s.HeaderText = lg.NewStyle().
//		Foreground(indigo).
//		Bold(true).
//		Padding(0, 1, 0, 2)
//	s.Status = lg.NewStyle().
//		Border(lipgloss.RoundedBorder()).
//		BorderForeground(indigo).
//		PaddingLeft(1).
//		MarginTop(1)
//	s.StatusHeader = lg.NewStyle().
//		Foreground(green).
//		Bold(true)
//	s.Highlight = lg.NewStyle().
//		Foreground(lipgloss.Color("212"))
//	s.ErrorHeaderText = s.HeaderText.Copy().
//		Foreground(red)
//	s.Help = lg.NewStyle().
//		Foreground(lipgloss.Color("240"))
//	return &s
//}
//
//func (m Model) Init() tea.Cmd {
//	return tea.Batch(m.forms[0].Init(), m.spinner.Tick)
//}
//
//func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
//	var cmds []tea.Cmd
//	switch msg := message.(type) {
//	case tea.WindowSizeMsg:
//		m.width = min(msg.Width, 120) - m.styles.Base.GetHorizontalFrameSize()
//	case tea.KeyMsg:
//		switch msg.String() {
//		case "esc", "ctrl+c":
//			return m, tea.Quit
//		}
//
//	default:
//		var cmd tea.Cmd
//		m.spinner, cmd = m.spinner.Update(msg)
//		cmds = append(cmds, cmd)
//	}
//
//	// Process the form
//	form, cmd := m.forms[m.state].Update(message)
//	if f, ok := form.(*huh.Form); ok {
//		m.forms[m.state] = f
//		cmds = append(cmds, cmd)
//	}
//	return m, tea.Batch(cmds...)
//}
//
//func (m Model) View() string {
//	var body string
//	if m.forms[m.state].GetString("authMethod") != "" && m.state == Initial {
//		switch m.forms[m.state].GetString("authMethod") {
//		case "SignIn":
//			m.changedForm = true
//			m.state = SignIn
//			m.forms[SignIn].Init()()
//		case "SignUp":
//			m.state = SignUp
//		}
//	}
//	if m.forms[m.state].GetString("email") != "" && m.forms[m.state].GetString("password") != "" && m.state == SignIn {
//		err := m.client.SignIn(m.forms[m.state].GetString("email"), m.forms[m.state].GetString("password"))
//		if err != nil {
//			m.clientError = err
//		} else {
//			m.clientError = nil
//			m.state = Dashboard
//		}
//	}
//	errs := m.forms[m.state].Errors()
//	if m.clientError != nil {
//		errs = append(errs, m.clientError)
//	}
//
//	header := m.appBoundaryView("GophKeeper App")
//	if len(errs) > 0 {
//		header = m.appErrorBoundaryView(m.errorView(errs))
//	}
//
//	switch m.state {
//	case Initial, SignIn:
//		v := strings.TrimSuffix(m.forms[m.state].View(), "\n\n")
//		form := m.lg.NewStyle().Margin(1, 0).Render(v)
//		body = m.styles.Base.Render(lipgloss.JoinHorizontal(lipgloss.Top, form))
//	default:
//		body = m.styles.Base.Render(m.spinner.View(), "Oh-oh, something crashed... press esc to quit")
//	}
//	return header + "\n" + body
//}
//
//func (m Model) appBoundaryView(text string) string {
//	return lipgloss.PlaceHorizontal(
//		m.width,
//		lipgloss.Left,
//		m.styles.HeaderText.Render(text),
//		lipgloss.WithWhitespaceChars("/"),
//		lipgloss.WithWhitespaceForeground(indigo),
//	)
//}
//
//func (m Model) appErrorBoundaryView(text string) string {
//	return lipgloss.PlaceHorizontal(
//		m.width,
//		lipgloss.Left,
//		m.styles.ErrorHeaderText.Render(text),
//		lipgloss.WithWhitespaceChars("/"),
//		lipgloss.WithWhitespaceForeground(red),
//	)
//}
//
//func (m Model) errorView(errs []error) string {
//	var s string
//	for _, err := range errs {
//		s += err.Error()
//	}
//	return s
//}
//
//func main() {
//	file, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer file.Close()
//	log.SetOutput(file)
//
//	p := tea.NewProgram(NewModel(0), tea.WithAltScreen())
//	if _, err := p.Run(); err != nil {
//		log.Fatal(err)
//	}
//}
