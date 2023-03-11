package main

import (
	"fmt"

	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/* MODEL MANAGEMENT */
const divisor = 4

type status int

const (
	todo status = iota
	inProgress
	done
)

const (
	mainModel status = iota
	formModel
)

var models []tea.Model

/* STYLING */

var (
	columnStyle = lipgloss.NewStyle().
			Padding(1, 2)
	focusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

/* CUSTOM ITEM */

type Task struct {
	status      status
	title       string
	description string
}

func NewTask(status status, title, description string) Task {
	return Task{status: status, title: title, description: description}
}

func (t *Task) increaseTaskStatus() {
	if t.status == done {
		t.status = todo
	} else {
		t.status++
	}
}

// implement the list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

/* MAIN MODEL */

type Model struct {
	loaded   bool
	focused  status
	lists    []list.Model
	err      error
	quitting bool
}

func New() *Model {
	return &Model{}
}

func (m *Model) moveTaskToNext() {
	selectedTask := m.lists[m.focused].SelectedItem().(Task)
	m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
	selectedTask.increaseTaskStatus()
	m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, list.Item(selectedTask))

}

func (m *Model) focusOnNextList() {
	if m.focused == done {
		m.focused = todo
	} else {
		m.focused++
	}
}
func (m *Model) focusOnPrevList() {
	if m.focused == todo {
		m.focused = done
	} else {
		m.focused--
	}
}

func (m *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height/2)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	// init To Do
	m.lists[todo].Title = "To Do"
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "estudar EST", description: "materia do prof rabanete"},
		Task{status: todo, title: "estudar ELE", description: "materia do governador da bahia"},
	})

	// init in progress
	m.lists[inProgress].Title = "In Progress"
	m.lists[inProgress].SetItems([]list.Item{
		Task{status: inProgress, title: "study go", description: "learn go syntax"},
		Task{status: inProgress, title: "study bubbletea", description: "learning go's bubbletea lib"},
	})

	// init done
	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		Task{status: done, title: "leetcode grind", description: "get ready for the tech interview"},
		Task{status: done, title: "code more", description: "be a more productive programmer"},
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// message that gives dimensions of the terminal
	case tea.WindowSizeMsg:
		if !m.loaded {
			columnStyle.Width(msg.Width / divisor)
			columnStyle.Height(msg.Height / 2)
			focusedStyle.Width(msg.Width / divisor)
			focusedStyle.Height(msg.Height / 2)
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "right":
			m.focusOnNextList()
		case "left":
			m.focusOnPrevList()
		case "enter":
			m.moveTaskToNext()
		case "n":
			models[mainModel] = m // save state of current model
			models[formModel] = NewForm(m.focused)
			return models[formModel].Update(nil)
		}
	case Task:
		task := msg
		return m, m.lists[task.status].InsertItem(len(m.lists[task.status].Items()), list.Item(task))
	}
	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	if m.loaded {
		todoView := m.lists[todo].View()
		inProgressView := m.lists[inProgress].View()
		doneView := m.lists[done].View()

		switch m.focused {
		case inProgress:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				focusedStyle.Render(inProgressView),
				columnStyle.Render(doneView),
			)
		case done:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				columnStyle.Render(inProgressView),
				focusedStyle.Render(doneView),
			)
		default:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				focusedStyle.Render(todoView),
				columnStyle.Render(inProgressView),
				columnStyle.Render(doneView),
			)
		}
	} else {
		return "loading..."
	}
}

/* FORM MODEL */
type Form struct {
	focused     status
	title       textinput.Model
	description textarea.Model
}

func NewForm(focused status) *Form {
	form := &Form{}
	form.focused = focused
	form.title = textinput.New()
	form.title.Focus() // set the focus state on the model so it can receive input
	form.description = textarea.New()
	return form
}

func (m Form) CreateTask() tea.Msg {
	return NewTask(m.focused, m.title.Value(), m.description.Value())
}

func (m Form) Init() tea.Cmd {
	return nil
}

func (m Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.title.Focused() {
				m.title.Blur() // remove the focus from title
				m.description.Focus()
				return m, textarea.Blink
			} else {
				models[formModel] = m
				return models[mainModel], m.CreateTask
			}
		}
	}
	if m.title.Focused() {
		m.title, cmd = m.title.Update(msg)
	} else {
		m.description, cmd = m.description.Update(msg)
	}

	return m, cmd
}

func (m Form) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.title.View(),
		m.description.View(),
	)
}

func main() {
	models = []tea.Model{New(), NewForm(todo)}
	m := models[mainModel]
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
