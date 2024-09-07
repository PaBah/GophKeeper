package main

import (
	"context"
	"strings"
	"time"

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
	cursor          int
	tableCursor     int
	menu            []string
	content         string
	tableNavigation bool
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
	credentials, err := m.clientService.GetCredentials(context.Background())
	if err != nil {
		m.err = err
		return ""
	}
	headers := lipgloss.JoinHorizontal(
		lipgloss.Top,
		headerStyle.Render("ID"),
		headerStyle.Render("ServiceName"),
		headerStyle.Render("UploadedAt"),
	)
	tableData := []string{borderStyle.Render(headers)}
	for index, credential := range credentials {
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
	cards, err := m.clientService.GetCards(context.Background())
	if err != nil {
		m.err = err
		return ""
	}
	headers := lipgloss.JoinHorizontal(
		lipgloss.Top,
		headerStyle.Render("LastDigits"),
		headerStyle.Render("ExpirationDate"),
		headerStyle.Render("UploadedAt"),
	)
	tableData := []string{borderStyle.Render(headers)}
	for index, card := range cards {
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

func (ds *DashboardScreen) drawContent(m *Model) string {
	switch ds.cursor {
	case credentials:
		return ds.drawCredentials(m)
	case cards:
		return ds.drawCards(m)
	default:
		return ""
	}
}

func (ds *DashboardScreen) getListAmount(m *Model) int {
	switch ds.cursor {
	case credentials:
		list, _ := m.clientService.GetCredentials(context.Background())
		return len(list)
	case cards:
		list, _ := m.clientService.GetCards(context.Background())
		return len(list)
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

func (ds *DashboardScreen) Update(m *Model, msg tea.Msg) (tea.Model, tea.Cmd) {
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
			default:
				return m, nil
			}
		case "f2":
			switch ds.cursor {
			case credentials:
				m.credentialsScreen.createMode = false
				list, _ := m.clientService.GetCredentials(context.Background())
				m.credentialsScreen.inputs[serviceName].SetValue(list[ds.tableCursor].ServiceName)
				m.credentialsScreen.inputs[identity].SetValue(list[ds.tableCursor].Identity)
				m.credentialsScreen.inputs[password].SetValue(list[ds.tableCursor].Password)
				m.credentialsScreen.updateID = list[ds.tableCursor].ID
				m.state = CredentialsForm
			case cards:
				m.cardsScreen.createMode = false
				list, _ := m.clientService.GetCards(context.Background())
				m.cardsScreen.inputs[cardNumber].SetValue(formatCardNumber(list[ds.tableCursor].Number))
				m.cardsScreen.inputs[expiryDate].SetValue(list[ds.tableCursor].ExpirationDate)
				m.cardsScreen.inputs[cardHolder].SetValue(list[ds.tableCursor].HolderName)
				m.cardsScreen.inputs[cvv].SetValue(list[ds.tableCursor].CVV)
				m.cardsScreen.updateID = list[ds.tableCursor].ID
				m.state = CardForm
			default:
				return m, nil
			}
		case "f3":
			switch ds.cursor {
			case credentials:
				list, _ := m.clientService.GetCredentials(context.Background())
				_ = m.clientService.DeleteCredentials(context.Background(), list[ds.tableCursor].ID)
				ds.tableCursor = max(ds.tableCursor-1, 0)
				ds.content = ds.drawContent(m)
			case cards:
				list, _ := m.clientService.GetCards(context.Background())
				_ = m.clientService.DeleteCard(context.Background(), list[ds.tableCursor].ID)
				ds.tableCursor = max(ds.tableCursor-1, 0)
				ds.content = ds.drawContent(m)
			default:
				return m, nil
			}
		case "f4":
			switch ds.cursor {
			case credentials:
				list, _ := m.clientService.GetCredentials(context.Background())
				_ = clipboard.WriteAll(list[ds.tableCursor].Identity)
			case cards:
				list, _ := m.clientService.GetCards(context.Background())
				_ = clipboard.WriteAll(list[ds.tableCursor].Number)
			default:
				return m, nil
			}
		case "f5":
			switch ds.cursor {
			case credentials:
				list, _ := m.clientService.GetCredentials(context.Background())
				_ = clipboard.WriteAll(list[ds.tableCursor].Password)
			case cards:
				list, _ := m.clientService.GetCards(context.Background())
				_ = clipboard.WriteAll(list[ds.tableCursor].ExpirationDate)
			default:
				return m, nil
			}
		case "f6":
			switch ds.cursor {
			case cards:
				list, _ := m.clientService.GetCards(context.Background())
				_ = clipboard.WriteAll(list[ds.tableCursor].HolderName)
			default:
				return m, nil
			}
		case "f7":
			switch ds.cursor {
			case cards:
				list, _ := m.clientService.GetCards(context.Background())
				_ = clipboard.WriteAll(list[ds.tableCursor].CVV)
			default:
				return m, nil
			}
		case "enter":
			if !ds.tableNavigation {
				ds.tableNavigation = true
				ds.tableCursor = 0
				ds.content = ds.drawContent(m)
				if ds.menu[ds.cursor] == "Exit" {
					return m, tea.Quit
				}
			}
		}
	}

	return m, nil
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
