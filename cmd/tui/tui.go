package tui

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"
	"viadro_cli/cmd/cli"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var TuiCmd = &cobra.Command{
	Use:     "interactive",
	Aliases: []string{"i", "tui", "int"},
	Example: "viadro interactive",
	Short:   "Launch app in interactive (TUI) mode",
	Long:    "Launch app in interactive (TUI) mode",
	Run:     tuiFunc,
}

func tuiFunc(cli *cobra.Command, args []string) {
	setup()
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Item struct {
	id         int       //id
	title      string    //title
	desc       string    //link
	tags       []string  //tags
	uploadedat time.Time //uploaded
}

func (i Item) Title() string { return i.title }
func (i Item) Description() string {
	timeFormatted := i.uploadedat.Format(time.RFC822)
	// isHidden := ""
	// if i.hidden == true {
	// 	isHidden = "hidden"
	// } else {
	// 	isHidden = "visible"
	// }
	return fmt.Sprintf("ID: %d • LINK: %s • TAGS: %v • UPLOADED: %s", i.id, i.desc, i.tags, timeFormatted)
}
func (i Item) DescriptionValue() string { return i.desc }
func (i Item) FilterValue() string      { return i.title }
func (i Item) IdValue() int             { return i.id }

type model struct {
	list        list.Model
	CurrentMode string
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

		if msg.String() == "2" {
			m.CurrentMode = "me"
			itemsGet := cli.List([]string{}, "me")

			items := []list.Item{}

			for _, val := range itemsGet.Documents {
				items = append(items, Item{val.DocumentID, val.Title, val.Link, val.Tags, val.CreatedAt})
			}
			m.list.SetItems(items)
			m.list.Title = "My Documents"

			return m, nil
		}

		if msg.String() == "1" {
			m.CurrentMode = "all"
			itemsGet := cli.List([]string{}, "all")

			items := []list.Item{}

			for _, val := range itemsGet.Documents {
				items = append(items, Item{val.DocumentID, val.Title, val.Link, val.Tags, val.CreatedAt})
			}
			m.list.SetItems(items)
			m.list.Title = "All Documents"

			return m, nil
		}

		if msg.String() == "3" {
			m.CurrentMode = "exclude"
			itemsGet := cli.List([]string{}, "exclude")

			items := []list.Item{}

			for _, val := range itemsGet.Documents {
				items = append(items, Item{val.DocumentID, val.Title, val.Link, val.Tags, val.CreatedAt})
			}
			m.list.SetItems(items)
			m.list.Title = "All Documents (excl. my)"

			return m, nil
		}

		if msg.String() == "d" {
			selectedItem := m.list.SelectedItem()
			selectedItemId := selectedItem.IdValue()
			selectedItemIdStr := strconv.Itoa(selectedItemId)

			msg := cli.Delete([]string{selectedItemIdStr})
			statusCmd := m.list.NewStatusMessage(msg)
			itemsGet := cli.List([]string{}, m.CurrentMode)

			items := []list.Item{}

			for _, val := range itemsGet.Documents {
				items = append(items, Item{val.DocumentID, val.Title, val.Link, val.Tags, val.CreatedAt})
			}
			m.list.SetItems(items)

			return m, tea.Batch(statusCmd)
		}

		if msg.String() == "t" {
			selectedItem := m.list.SelectedItem()
			selectedItemId := selectedItem.IdValue()
			selectedItemIdStr := strconv.Itoa(selectedItemId)

			msg := cli.Toggle([]string{selectedItemIdStr})
			statusCmd := m.list.NewStatusMessage(msg)
			itemsGet := cli.List([]string{}, m.CurrentMode)

			items := []list.Item{}

			for _, val := range itemsGet.Documents {
				items = append(items, Item{val.DocumentID, val.Title, val.Link, val.Tags, val.CreatedAt})
			}
			m.list.SetItems(items)

			return m, tea.Batch(statusCmd)
		}

		if msg.String() == "enter" {
			selectedItem := m.list.SelectedItem()
			exec.Command("xdg-open", selectedItem.DescriptionValue()).Start()

			statusCmd := m.list.NewStatusMessage("Link opened in browser")

			return m, tea.Batch(statusCmd)
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
	documents := cli.List([]string{}, "all")
	items := []list.Item{}

	for _, val := range documents.Documents {
		items = append(items, Item{val.DocumentID, val.Title, val.Link, val.Tags, val.CreatedAt})
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0), CurrentMode: "all"}
	m.list.Title = "All documents"

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		cli.Logger.Fatal("app error")
	}
}
