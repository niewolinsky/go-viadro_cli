package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"viadro_cli/cmd/cli"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var TuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "tui",
	Long:  ``,
	Run:   tuiFunc,
}

func tuiFunc(cli *cobra.Command, args []string) {
	setup()
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Item struct {
	title, desc string
	id          int
}

func (i Item) Title() string            { return i.title }
func (i Item) Description() string      { return i.desc }
func (i Item) DescriptionValue() string { return i.desc }
func (i Item) FilterValue() string      { return i.title }
func (i Item) IdValue() int             { return i.id }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if msg.String() == "p" {
			itemsGet := cli.ListTesting([]string{}, "me")

			items := []list.Item{}

			for _, val := range itemsGet.Documents {
				items = append(items, Item{val.Title, val.Link, val.DocumentID})
			}
			m.list.SetItems(items)
			m.list.Title = "My Documents"

			return m, nil
		}

		if msg.String() == "o" {
			itemsGet := cli.ListTesting([]string{}, "all")

			items := []list.Item{}

			for _, val := range itemsGet.Documents {
				items = append(items, Item{val.Title, val.Link, val.DocumentID})
			}
			m.list.SetItems(items)
			m.list.Title = "All Documents"

			return m, nil
		}

		if msg.String() == "i" {
			itemsGet := cli.ListTesting([]string{}, "exclude")

			items := []list.Item{}

			for _, val := range itemsGet.Documents {
				items = append(items, Item{val.Title, val.Link, val.DocumentID})
			}
			m.list.SetItems(items)
			m.list.Title = "All Documents (excl. my)"

			return m, nil
		}

		if msg.String() == "t" {
			// x := m.list.Index()
			y := m.list.SelectedItem()
			z := y.IdValue()
			zStr := strconv.Itoa(z)

			msg := cli.Toggle([]string{zStr})
			// statusCmd := m.list.NewStatusMessage(fmt.Sprintf("Toggled visibility of document with ID: %d", x))
			statusCmd := m.list.NewStatusMessage(msg)
			itemsGet := cli.ListTesting([]string{}, "all")

			items := []list.Item{}

			for _, val := range itemsGet.Documents {
				items = append(items, Item{val.Title, val.Link, val.DocumentID})
			}
			m.list.SetItems(items)

			return m, tea.Batch(statusCmd)
		}

		if msg.String() == "enter" {
			y := m.list.SelectedItem()
			exec.Command("xdg-open", y.DescriptionValue()).Start()

			return m, nil
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func setup() {
	itemsGet := cli.ListTesting([]string{}, "all")

	items := []list.Item{}

	for _, val := range itemsGet.Documents {
		items = append(items, Item{val.Title, val.Link, val.DocumentID})
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "All documents"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
