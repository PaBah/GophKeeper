package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/PaBah/GophKeeper/internal/models"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultMessage = "Welcome to the dashboard. Select an option from the menu on the left!"
)

var (
	headerStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("63")).
			Foreground(lipgloss.Color("230")).
			Bold(true).
			Padding(0, 2)

	cellStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Padding(0, 2)

	selectedCellStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("205")).
				Foreground(lipgloss.Color("230")).
				Padding(0, 2)

	menuWidth    = 20
	menuStyle    = lipgloss.NewStyle().Padding(1, 2).Width(menuWidth)
	activeMenu   = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	inactiveMenu = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	borderStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63"))
	contentStyle = lipgloss.NewStyle().Padding(1, 2)
)

const (
	credentials = iota
	cards
	files
	exit
)

type DashboardScreen struct {
	cursor           int
	tableCursor      int
	cardsState       []models.Card
	credentialsState []models.Credentials
	filesState       []models.File
	menu             []string
	content          string
	updateMsg        string
	tableNavigation  bool
}

func NewDashboardScreen() *DashboardScreen {
	return &DashboardScreen{
		menu:            []string{"Credentials", "Cards", "Files", "Exit"},
		content:         defaultMessage,
		tableNavigation: false,
	}
}

func (ds *DashboardScreen) renderRow(index int, cols ...string) string {
	style := cellStyle
	if index == ds.tableCursor {
		style = selectedCellStyle
	}

	var row []string
	for _, col := range cols {
		row = append(row, style.Render(col))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, row...)
}

func (ds *DashboardScreen) drawCredentials(m *Model) string {
	headers := lipgloss.JoinHorizontal(
		lipgloss.Top,
		headerStyle.Render("ID"),
		headerStyle.Render("ServiceName"),
		headerStyle.Render("UploadedAt"),
	)
	tableData := []string{borderStyle.Render(headers)}
	for index, credential := range ds.credentialsState {
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			ds.renderRow(index, credential.ID, credential.ServiceName, credential.UploadedAt.Format(time.RFC3339)),
		)
		tableData = append(tableData, borderStyle.Render(row))
	}

	table := lipgloss.JoinVertical(
		lipgloss.Left,
		tableData...,
	)

	return table
}

func (ds *DashboardScreen) drawCards(m *Model) string {
	headers := lipgloss.JoinHorizontal(
		lipgloss.Top,
		headerStyle.Render("LastDigits"),
		headerStyle.Render("ExpirationDate"),
		headerStyle.Render("UploadedAt"),
	)
	tableData := []string{borderStyle.Render(headers)}
	for index, card := range ds.cardsState {
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			ds.renderRow(index, "*"+card.Number[12:], card.ExpirationDate, card.UploadedAt.Format(time.RFC3339)),
		)
		tableData = append(tableData, borderStyle.Render(row))
	}

	table := lipgloss.JoinVertical(
		lipgloss.Left,
		tableData...,
	)

	return table
}

func (ds *DashboardScreen) drawFiles(m *Model) string {
	headers := lipgloss.JoinHorizontal(
		lipgloss.Top,
		headerStyle.Render("Name"),
		headerStyle.Render("UploadedAt"),
		headerStyle.Render("Size"),
	)
	tableData := []string{borderStyle.Render(headers)}
	for index, file := range ds.filesState {
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			ds.renderRow(index, file.Name, file.UploadedAt.Format(time.RFC3339), file.Size),
		)
		tableData = append(tableData, borderStyle.Render(row))
	}

	table := lipgloss.JoinVertical(
		lipgloss.Left,
		tableData...,
	)

	return table
}

func (ds *DashboardScreen) drawContent(m *Model) string {
	switch ds.cursor {
	case credentials:
		return ds.drawCredentials(m)
	case cards:
		return ds.drawCards(m)
	case files:
		return ds.drawFiles(m)
	default:
		return ""
	}
}

func (ds *DashboardScreen) getListAmount(m *Model) int {
	switch ds.cursor {
	case credentials:
		return len(ds.credentialsState)
	case cards:
		return len(ds.cardsState)
	case files:
		return len(ds.filesState)
	default:
		return 0
	}
}

func formatCardNumber(cardNumber string) string {
	var result []string

	for i := 0; i < len(cardNumber); i += 4 {
		end := i + 4
		if end > len(cardNumber) {
			end = len(cardNumber)
		}
		result = append(result, cardNumber[i:end])
	}

	return strings.Join(result, " ")
}

func (ds *DashboardScreen) _Update(m *Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	if ds.updateMsg != "" {
		m.err = errors.New(ds.updateMsg)
	} else {
		m.err = nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if !ds.tableNavigation {
				if ds.cursor > 0 {
					ds.cursor--
				}
			} else {
				if ds.tableCursor > 0 {
					ds.tableCursor--
					ds.content = ds.drawContent(m)
				}
			}
		case "down":
			if !ds.tableNavigation {
				if ds.cursor < len(ds.menu)-1 {
					ds.cursor++
				}
			} else {
				if ds.tableCursor < ds.getListAmount(m)-1 {
					ds.tableCursor++
					ds.content = ds.drawContent(m)
				}
			}
		case "left":
			if ds.tableNavigation {
				ds.content = defaultMessage
				ds.tableNavigation = false
			}
		case "f1":
			switch ds.cursor {
			case credentials:
				m.credentialsScreen.createMode = true
				m.credentialsScreen.inputs[serviceName].SetValue("")
				m.credentialsScreen.inputs[identity].SetValue("")
				m.credentialsScreen.inputs[password].SetValue("")
				m.credentialsScreen.updateID = ""
				m.state = CredentialsForm
			case cards:
				m.cardsScreen.createMode = true
				m.cardsScreen.inputs[cardNumber].SetValue("")
				m.cardsScreen.inputs[expiryDate].SetValue("")
				m.cardsScreen.inputs[cardHolder].SetValue("")
				m.cardsScreen.inputs[cvv].SetValue("")
				m.cardsScreen.updateID = ""
				m.state = CardForm
			case files:
				m.state = FileLoad
				return m, m.filesScreen.filepicker.Init()
			default:
				return m, nil
			}
		case "f2":
			switch ds.cursor {
			case credentials:
				m.credentialsScreen.createMode = false
				m.credentialsScreen.inputs[serviceName].SetValue(ds.credentialsState[ds.tableCursor].ServiceName)
				m.credentialsScreen.inputs[identity].SetValue(ds.credentialsState[ds.tableCursor].Identity)
				m.credentialsScreen.inputs[password].SetValue(ds.credentialsState[ds.tableCursor].Password)
				m.credentialsScreen.updateID = ds.credentialsState[ds.tableCursor].ID
				m.state = CredentialsForm
			case cards:
				m.cardsScreen.createMode = false
				m.cardsScreen.inputs[cardNumber].SetValue(formatCardNumber(ds.cardsState[ds.tableCursor].Number))
				m.cardsScreen.inputs[expiryDate].SetValue(ds.cardsState[ds.tableCursor].ExpirationDate)
				m.cardsScreen.inputs[cardHolder].SetValue(ds.cardsState[ds.tableCursor].HolderName)
				m.cardsScreen.inputs[cvv].SetValue(ds.cardsState[ds.tableCursor].CVV)
				m.cardsScreen.updateID = ds.cardsState[ds.tableCursor].ID
				m.state = CardForm
			case files:
				m.clientService.DownloadsFile(context.Background(), ds.filesState[ds.tableCursor].Name)
			default:
				return m, nil
			}
		case "f3":
			switch ds.cursor {
			case credentials:
				_ = m.clientService.DeleteCredentials(context.Background(), ds.credentialsState[ds.tableCursor].ID)
				ds.tableCursor = max(ds.tableCursor-1, 0)
				ds.loadActual(m)
				ds.content = ds.drawContent(m)
			case cards:
				_ = m.clientService.DeleteCard(context.Background(), ds.cardsState[ds.tableCursor].ID)
				ds.tableCursor = max(ds.tableCursor-1, 0)
				ds.loadActual(m)
				ds.content = ds.drawContent(m)
			case files:
				_ = m.clientService.DeleteFile(context.Background(), ds.filesState[ds.tableCursor].Name)
				ds.tableCursor = max(ds.tableCursor-1, 0)
				ds.loadActual(m)
				ds.content = ds.drawContent(m)
			default:
				return m, nil
			}
		case "f4":
			switch ds.cursor {
			case credentials:
				_ = clipboard.WriteAll(ds.credentialsState[ds.tableCursor].Identity)
			case cards:
				_ = clipboard.WriteAll(ds.cardsState[ds.tableCursor].Number)
			default:
				return m, nil
			}
		case "f5":
			switch ds.cursor {
			case credentials:
				_ = clipboard.WriteAll(ds.credentialsState[ds.tableCursor].Password)
			case cards:
				_ = clipboard.WriteAll(ds.cardsState[ds.tableCursor].ExpirationDate)
			default:
				return m, nil
			}
		case "f6":
			switch ds.cursor {
			case cards:
				_ = clipboard.WriteAll(ds.cardsState[ds.tableCursor].HolderName)
			default:
				return m, nil
			}
		case "f7":
			switch ds.cursor {
			case cards:
				_ = clipboard.WriteAll(ds.cardsState[ds.tableCursor].CVV)
			default:
				return m, nil
			}
		case "shift+right":
			if m.err != nil {
				ds.loadActual(m)
				ds.content = ds.drawContent(m)
			}
		case "enter":
			if !ds.tableNavigation {
				ds.tableNavigation = true
				ds.tableCursor = 0
				ds.loadActual(m)
				ds.content = ds.drawContent(m)
				if ds.cursor == exit {
					return m, tea.Quit
				}
			}
		}
	}
	return m, nil
}

func (ds *DashboardScreen) Update(m *Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	ds.updateErrorMessage(m)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return ds.handleKeyMsg(m, msg)
	}

	return m, nil
}

func (ds *DashboardScreen) updateErrorMessage(m *Model) {
	if ds.updateMsg != "" {
		m.err = errors.New(ds.updateMsg)
	} else {
		m.err = nil
	}
}

func (ds *DashboardScreen) handleKeyMsg(m *Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up":
		ds.handleUpKey(m)
	case "down":
		ds.handleDownKey(m)
	case "left":
		ds.handleLeftKey()
	case "f1":
		return ds.handleF1Key(m)
	case "f2":
		return ds.handleF2Key(m)
	case "f3":
		return ds.handleF3Key(m)
	case "f4":
		return ds.handleF4Key(m)
	case "f5":
		return ds.handleF5Key(m)
	case "f6":
		return ds.handleF6Key(m)
	case "f7":
		return ds.handleF7Key(m)
	case "shift+right":
		ds.handleShiftRightKey(m)
	case "enter":
		return ds.handleEnterKey(m)
	}

	return m, nil
}

func (ds *DashboardScreen) handleUpKey(m *Model) {
	if !ds.tableNavigation {
		if ds.cursor > 0 {
			ds.cursor--
		}
	} else {
		if ds.tableCursor > 0 {
			ds.tableCursor--
			ds.content = ds.drawContent(m)
		}
	}
}

func (ds *DashboardScreen) handleDownKey(m *Model) {
	if !ds.tableNavigation {
		if ds.cursor < len(ds.menu)-1 {
			ds.cursor++
		}
	} else {
		if ds.tableCursor < ds.getListAmount(m)-1 {
			ds.tableCursor++
			ds.content = ds.drawContent(m)
		}
	}
}

func (ds *DashboardScreen) handleLeftKey() {
	if ds.tableNavigation {
		ds.content = defaultMessage
		ds.tableNavigation = false
	}
}

func (ds *DashboardScreen) handleCredentialsCreate(m *Model) {
	m.credentialsScreen.createMode = true
	m.credentialsScreen.inputs[serviceName].SetValue("")
	m.credentialsScreen.inputs[identity].SetValue("")
	m.credentialsScreen.inputs[password].SetValue("")
	m.credentialsScreen.updateID = ""
	m.state = CredentialsForm
}

func (ds *DashboardScreen) handleCardsCreate(m *Model) {
	m.cardsScreen.createMode = false
	m.cardsScreen.inputs[cardNumber].SetValue(formatCardNumber(ds.cardsState[ds.tableCursor].Number))
	m.cardsScreen.inputs[expiryDate].SetValue(ds.cardsState[ds.tableCursor].ExpirationDate)
	m.cardsScreen.inputs[cardHolder].SetValue(ds.cardsState[ds.tableCursor].HolderName)
	m.cardsScreen.inputs[cvv].SetValue(ds.cardsState[ds.tableCursor].CVV)
	m.cardsScreen.updateID = ds.cardsState[ds.tableCursor].ID
	m.state = CardForm
}

func (ds *DashboardScreen) handleF1Key(m *Model) (tea.Model, tea.Cmd) {
	switch ds.cursor {
	case credentials:
		ds.handleCredentialsCreate(m)
	case cards:
		ds.handleCardsCreate(m)
	case files:
		m.state = FileLoad
		return m, m.filesScreen.filepicker.Init()
	default:
		return m, nil
	}
	return m, nil
}

func (ds *DashboardScreen) handleCredentialsEdit(m *Model) {
	m.credentialsScreen.createMode = false
	m.credentialsScreen.inputs[serviceName].SetValue(ds.credentialsState[ds.tableCursor].ServiceName)
	m.credentialsScreen.inputs[identity].SetValue(ds.credentialsState[ds.tableCursor].Identity)
	m.credentialsScreen.inputs[password].SetValue(ds.credentialsState[ds.tableCursor].Password)
	m.credentialsScreen.updateID = ds.credentialsState[ds.tableCursor].ID
	m.state = CredentialsForm
}

func (ds *DashboardScreen) handleCardsEdit(m *Model) {
	m.cardsScreen.createMode = false
	m.cardsScreen.inputs[cardNumber].SetValue(formatCardNumber(ds.cardsState[ds.tableCursor].Number))
	m.cardsScreen.inputs[expiryDate].SetValue(ds.cardsState[ds.tableCursor].ExpirationDate)
	m.cardsScreen.inputs[cardHolder].SetValue(ds.cardsState[ds.tableCursor].HolderName)
	m.cardsScreen.inputs[cvv].SetValue(ds.cardsState[ds.tableCursor].CVV)
	m.cardsScreen.updateID = ds.cardsState[ds.tableCursor].ID
	m.state = CardForm
}

func (ds *DashboardScreen) handleF2Key(m *Model) (tea.Model, tea.Cmd) {
	switch ds.cursor {
	case credentials:
		ds.handleCredentialsEdit(m)
	case cards:
		ds.handleCardsEdit(m)
	case files:
		m.clientService.DownloadsFile(context.Background(), ds.filesState[ds.tableCursor].Name)
	default:
		return m, nil
	}
	return m, nil
}

func (ds *DashboardScreen) deleteCredentials(m *Model) {
	_ = m.clientService.DeleteCredentials(context.Background(), ds.credentialsState[ds.tableCursor].ID)
	ds.tableCursor = max(ds.tableCursor-1, 0)
	ds.loadActual(m)
	ds.content = ds.drawContent(m)
}

func (ds *DashboardScreen) deleteCard(m *Model) {
	_ = m.clientService.DeleteCard(context.Background(), ds.cardsState[ds.tableCursor].ID)
	ds.tableCursor = max(ds.tableCursor-1, 0)
	ds.loadActual(m)
	ds.content = ds.drawContent(m)
}

func (ds *DashboardScreen) deleteFile(m *Model) {
	_ = m.clientService.DeleteFile(context.Background(), ds.filesState[ds.tableCursor].Name)
	ds.tableCursor = max(ds.tableCursor-1, 0)
	ds.loadActual(m)
	ds.content = ds.drawContent(m)
}

func (ds *DashboardScreen) handleF3Key(m *Model) (tea.Model, tea.Cmd) {
	switch ds.cursor {
	case credentials:
		ds.deleteCredentials(m)
	case cards:
		ds.deleteCard(m)
	case files:
		ds.deleteFile(m)
	default:
		return m, nil
	}
	return m, nil
}

func (ds *DashboardScreen) handleF4Key(m *Model) (tea.Model, tea.Cmd) {
	if ds.cursor == credentials {
		_ = clipboard.WriteAll(ds.credentialsState[ds.tableCursor].Identity)
	} else if ds.cursor == cards {
		_ = clipboard.WriteAll(ds.cardsState[ds.tableCursor].Number)
	}
	return m, nil
}

func (ds *DashboardScreen) handleF5Key(m *Model) (tea.Model, tea.Cmd) {
	if ds.cursor == credentials {
		_ = clipboard.WriteAll(ds.credentialsState[ds.tableCursor].Password)
	} else if ds.cursor == cards {
		_ = clipboard.WriteAll(ds.cardsState[ds.tableCursor].ExpirationDate)
	}
	return m, nil
}

func (ds *DashboardScreen) handleF6Key(m *Model) (tea.Model, tea.Cmd) {
	if ds.cursor == cards {
		_ = clipboard.WriteAll(ds.cardsState[ds.tableCursor].HolderName)
	}
	return m, nil
}

func (ds *DashboardScreen) handleF7Key(m *Model) (tea.Model, tea.Cmd) {
	if ds.cursor == cards {
		_ = clipboard.WriteAll(ds.cardsState[ds.tableCursor].CVV)
	}
	return m, nil
}

func (ds *DashboardScreen) handleShiftRightKey(m *Model) {
	if m.err != nil {
		ds.loadActual(m)
		ds.content = ds.drawContent(m)
	}
}

func (ds *DashboardScreen) handleEnterKey(m *Model) (tea.Model, tea.Cmd) {
	if !ds.tableNavigation {
		ds.tableNavigation = true
		ds.tableCursor = 0
		ds.loadActual(m)
		ds.content = ds.drawContent(m)
		if ds.cursor == exit {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (ds *DashboardScreen) loadActual(m *Model) {
	ds.updateMsg = ""
	switch ds.cursor {
	case credentials:
		ds.credentialsState, _ = m.clientService.GetCredentials(context.Background())
	case cards:
		ds.cardsState, _ = m.clientService.GetCards(context.Background())
	case files:
		ds.filesState, _ = m.clientService.GetFiles(context.Background())
	default:
		ds.updateMsg = ""
	}
}

func (ds *DashboardScreen) View(m *Model) string {
	var menu string
	for i, choice := range ds.menu {
		if i == ds.cursor {
			menu += activeMenu.Render("> "+choice) + "\n"
		} else {
			menu += inactiveMenu.Render(choice) + "\n"
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		borderStyle.Render(menuStyle.Render(menu)),
		contentStyle.Render(ds.content),
	)
}
