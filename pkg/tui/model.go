package tui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sysadmin-cli/pkg/models"
)

type state int

const (
	categoriesState state = iota
	commandsState
	detailState
)

type Model struct {
	state      state
	categories []models.Category
	detail     *models.Command

	catCursor int
	cmdCursor int

	searchVisible bool
	searchInput   textinput.Model
	searchQuery   string

	filteredCats []models.Category
	filteredCmds []models.Command

	width     int
	height    int
	contentH  int
	viewport  viewport.Model
	styles    Styles
	ready     bool

	copiedTick int
}

func New(categories []models.Category) Model {
	ti := textinput.New()
	s := DefaultStyles()
	ti.PromptStyle = s.SearchPrompt
	ti.TextStyle = s.SearchText
	ti.PlaceholderStyle = s.SearchPlace
	ti.Cursor.Style = s.SearchCursor
	ti.CharLimit = 100
	ti.Width = 40

	return Model{
		categories:   categories,
		filteredCats: categories,
		searchInput:  ti,
		styles:       s,
		viewport:     viewport.New(0, 0),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.contentH = msg.Height - 4
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = m.contentH - 4
		if m.state == detailState && m.detail != nil {
			m.viewport.SetContent(m.renderDetailContent())
		}

	case tea.KeyMsg:
		if m.searchVisible && m.searchInput.Focused() {
			return m.handleSearchKey(msg)
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "c":
			if m.state == detailState && m.detail != nil {
				copyToClipboard(m.detail.Command)
				m.copiedTick = 3
			}
		case "esc":
			switch m.state {
			case detailState:
				m.state = commandsState
				m.detail = nil
			case commandsState:
				m.state = categoriesState
				m.cmdCursor = 0
			}
		case "/":
			m.searchVisible = true
			m.searchInput.Focus()
			m.searchInput.SetValue("")
			m.searchQuery = ""
			cmds = append(cmds, textinput.Blink)
		case "up", "k":
			m.moveUp()
		case "down", "j":
			m.moveDown()
		case "enter":
			m.enterCurrent()
		case "home", "g":
			m.goTop()
		case "end", "G":
			m.goBottom()
		}
	}

	if m.searchVisible && !m.searchInput.Focused() {
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.state == detailState {
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.copiedTick > 0 {
		m.copiedTick--
		if m.copiedTick == 0 {
			cmds = append(cmds, func() tea.Msg { return nil })
		}
	}

	return m, tea.Batch(cmds...)
}

func copyToClipboard(text string) {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	cmd.Run()
}

func (m *Model) handleSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.searchInput.Blur()
	case "esc":
		m.searchVisible = false
		m.searchInput.Blur()
		m.searchInput.SetValue("")
		m.searchQuery = ""
		m.applyFilter()
	default:
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		m.searchQuery = strings.ToLower(m.searchInput.Value())
		m.applyFilter()
		return m, cmd
	}
	return m, nil
}

func (m *Model) moveUp() {
	switch m.state {
	case categoriesState:
		if m.catCursor > 0 {
			m.catCursor--
		}
	case commandsState:
		if m.cmdCursor > 0 {
			m.cmdCursor--
		}
	}
}

func (m *Model) moveDown() {
	switch m.state {
	case categoriesState:
		if m.catCursor < len(m.filteredCats)-1 {
			m.catCursor++
		}
	case commandsState:
		if m.cmdCursor < len(m.filteredCmds)-1 {
			m.cmdCursor++
		}
	}
}

func (m *Model) goTop() {
	switch m.state {
	case categoriesState:
		m.catCursor = 0
	case commandsState:
		m.cmdCursor = 0
	}
}

func (m *Model) goBottom() {
	switch m.state {
	case categoriesState:
		m.catCursor = len(m.filteredCats) - 1
	case commandsState:
		m.cmdCursor = len(m.filteredCmds) - 1
	}
}

func (m *Model) enterCurrent() {
	switch m.state {
	case categoriesState:
		if len(m.filteredCats) > 0 && m.catCursor < len(m.filteredCats) {
			cat := m.filteredCats[m.catCursor]
			m.filteredCmds = cat.Commands
			m.cmdCursor = 0
			m.state = commandsState
		}
	case commandsState:
		if len(m.filteredCmds) > 0 && m.cmdCursor < len(m.filteredCmds) {
			cmd := m.filteredCmds[m.cmdCursor]
			m.detail = &cmd
			m.state = detailState
			m.viewport.SetContent(m.renderDetailContent())
			m.viewport.GotoTop()
		}
	}
}

func (m *Model) applyFilter() {
	if m.searchQuery == "" {
		switch m.state {
		case categoriesState:
			m.filteredCats = m.categories
		case commandsState:
			if len(m.categories) > 0 && m.catCursor < len(m.categories) {
				m.filteredCmds = m.categories[m.catCursor].Commands
			}
		}
		return
	}

	switch m.state {
	case categoriesState:
		var filtered []models.Category
		for _, cat := range m.categories {
			if cat.Matches(m.searchQuery) {
				filtered = append(filtered, cat)
			}
		}
		m.filteredCats = filtered
		if m.catCursor >= len(m.filteredCats) {
			m.catCursor = max(0, len(m.filteredCats)-1)
		}

	case commandsState:
		if m.catCursor < len(m.categories) {
			cat := m.categories[m.catCursor]
			var filtered []models.Command
			for _, cmd := range cat.Commands {
				if cmd.Matches(m.searchQuery) {
					filtered = append(filtered, cmd)
				}
			}
			m.filteredCmds = filtered
			if m.cmdCursor >= len(m.filteredCmds) {
				m.cmdCursor = max(0, len(m.filteredCmds)-1)
			}
		}
	}
}

func (m Model) View() string {
	if !m.ready {
		return "Cargando..."
	}

	var breadcrumb string
	var body string

	switch m.state {
	case categoriesState:
		body = m.renderCategories()
	case commandsState:
		catName := ""
		if len(m.filteredCats) > 0 && m.catCursor < len(m.filteredCats) {
			catName = m.filteredCats[m.catCursor].Name
		}
		breadcrumb = catName
		body = m.renderCommands()
	case detailState:
		catName := ""
		if len(m.categories) > 0 && m.catCursor < len(m.categories) {
			catName = m.categories[m.catCursor].Name
		}
		cmdName := ""
		if m.detail != nil {
			cmdName = m.detail.Command
		}
		breadcrumb = catName + " / " + cmdName
		body = m.renderDetail()
	}

	header := m.renderHeader(breadcrumb)
	footer := m.renderFooter()

	full := lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
	lines := strings.Split(full, "\n")
	if len(lines) > m.height {
		lines = lines[:m.height]
	}
	return strings.Join(lines, "\n")
}

func (m Model) renderHeader(breadcrumb string) string {
	title := m.styles.Breadcrumb.Render("sysadmin-cli")
	if breadcrumb != "" {
		title += "  " + m.styles.Dimmed.Render(">") + "  " + breadcrumb
	}

	right := m.styles.Dimmed.Render("[q] Salir")
	avail := m.width - lipgloss.Width(title) - lipgloss.Width(right) - 2
	if avail < 1 {
		avail = 1
	}
	pad := strings.Repeat(" ", avail)

	return m.styles.Header.Width(m.width).Render(title + pad + right)
}

func (m Model) renderCategories() string {
	if len(m.filteredCats) == 0 {
		return m.styles.Empty.Width(m.width - 2).Render("  Sin resultados. Intenta con otra busqueda.")
	}

	maxNameLen := 0
	for _, cat := range m.filteredCats {
		if l := len(cat.Name); l > maxNameLen {
			maxNameLen = l
		}
	}
	nameWidth := maxNameLen + 1
	if nameWidth < 10 {
		nameWidth = 10
	}

	avail := m.width - nameWidth - 10
	if avail < 15 {
		avail = 15
	}

	var items []string
	total := len(m.filteredCats)
	maxVisible := m.contentH - 4
	start := 0
	if m.catCursor >= maxVisible {
		start = m.catCursor - maxVisible + 1
	}

	for i, cat := range m.filteredCats {
		if i < start {
			continue
		}
		if len(items) >= maxVisible {
			break
		}

		num := fmt.Sprintf("%2d.", i+1)
		name := fmt.Sprintf("%-*s", nameWidth, cat.Name)
		desc := cat.Description
		if idx := strings.Index(desc, "."); idx > 0 {
			desc = desc[:idx+1]
		}
		desc += fmt.Sprintf(" (%d)", len(cat.Commands))
		if len(desc) > avail {
			desc = desc[:avail-3] + "..."
		}

		line := fmt.Sprintf("  %s  %s  %s", num, name, desc)

		if i == m.catCursor {
			items = append(items, m.styles.CatSelected.Width(m.width-2).Render(line))
		} else {
			items = append(items, m.styles.CatItem.Width(m.width-2).Render(line))
		}
	}

	if start > 0 || len(m.filteredCats) > maxVisible {
		scroll := fmt.Sprintf("  Mostrando %d de %d", start+1, total)
		items = append(items, m.styles.ScrollBar.Render(scroll))
	}

	return strings.Join(items, "\n")
}

func (m Model) renderCommands() string {
	if len(m.filteredCmds) == 0 {
		return m.styles.Empty.Width(m.width - 2).Render("  Sin resultados. Intenta con otra busqueda.")
	}

	maxCmdLen := 0
	for _, cmd := range m.filteredCmds {
		if l := len(cmd.Command); l > maxCmdLen {
			maxCmdLen = l
		}
	}
	cmdWidth := maxCmdLen + 2
	if cmdWidth < 10 {
		cmdWidth = 10
	}

	header := m.styles.CmdHeader.Render(fmt.Sprintf("  %-*s  %s", cmdWidth, "COMANDO", "DESCRIPCION"))
	div := m.styles.Divider.Render(strings.Repeat("-", m.width-4))

	var items []string
	items = append(items, "", header, div)

	total := len(m.filteredCmds)
	maxVisible := m.contentH - 6
	start := 0
	if m.cmdCursor >= maxVisible {
		start = m.cmdCursor - maxVisible + 1
	}

	for i, cmd := range m.filteredCmds {
		if i < start {
			continue
		}
		if len(items)-2 >= maxVisible {
			break
		}

		short := cmd.Title
		if short == "" {
			short = cmd.Description
		}

		if i == m.cmdCursor {
			nameColored := m.styles.CmdName.Render(fmt.Sprintf("%-*s", cmdWidth, cmd.Command))
			rendered := m.styles.CmdSelected.Width(m.width - 2).Render(fmt.Sprintf("  %s  %s", nameColored, short))
			items = append(items, rendered)
		} else {
			name := fmt.Sprintf("%-*s", cmdWidth, cmd.Command)
			line := fmt.Sprintf("  %s  %s", name, short)
			rendered := m.styles.CmdItem.Width(m.width - 2).Render(line)
			items = append(items, rendered)
		}
	}

	if start > 0 || total > maxVisible {
		scroll := fmt.Sprintf("  Mostrando %d de %d", start+1, total)
		items = append(items, m.styles.ScrollBar.Render(scroll))
	}

	return strings.Join(items, "\n")
}

func (m Model) renderDetailContent() string {
	if m.detail == nil {
		return ""
	}

	cmd := m.detail
	w := m.width - 10

	title := m.styles.DetailTitle.Render(cmd.Command)
	desc := m.styles.DetailValue.Render(cmd.Description)

	cmdLabel := m.styles.DetailLabel.Render("Comando:")
	cmdVal := m.styles.DetailCode.Render("  $ " + cmd.Command)

	exLabel := m.styles.DetailLabel.Render("Ejemplo:")
	exVal := m.styles.DetailCode.Render("  $ " + cmd.Example)

	tagsLabel := m.styles.DetailLabel.Render("Etiquetas:")
	tagsVal := m.styles.DetailValue.Render("  " + strings.Join(cmd.Tags, ", "))

	notesLabel := m.styles.DetailLabel.Render("Notas:")
	notesVal := m.styles.DetailNotes.Render("  " + cmd.Notes)

	content := strings.Join([]string{
		"",
		title,
		"",
		desc,
		"",
		cmdLabel,
		cmdVal,
		"",
		exLabel,
		exVal,
		"",
		tagsLabel,
		tagsVal,
		"",
		notesLabel,
		notesVal,
		"",
	}, "\n")

	return m.styles.DetailBox.Width(w).Render(content)
}

func (m Model) renderDetail() string {
	return m.viewport.View()
}

func (m Model) renderFooter() string {
	if m.searchVisible && m.searchInput.Focused() {
		prompt := m.styles.SearchPrompt.Render("> ")
		input := m.searchInput.View()
		return m.styles.Footer.Width(m.width).Render(prompt + input)
	}

	var binds []string
	switch m.state {
	case categoriesState:
		binds = []string{
			key("/", "Buscar"),
			key("up/down", "Navegar"),
			key("enter", "Seleccionar"),
			key("q", "Salir"),
		}
	case commandsState:
		binds = []string{
			key("/", "Buscar"),
			key("up/down", "Navegar"),
			key("enter", "Detalle"),
			key("esc", "Atras"),
		}
	case detailState:
		extra := ""
		if m.copiedTick > 0 {
			extra = m.styles.Highlight.Render("  Copiado!")
		}
		binds = []string{
			key("c", "Copiar comando"),
			key("esc", "Atras"),
		}
		if extra != "" {
			binds = append(binds, extra)
		}
	}

	return m.styles.Footer.Width(m.width).Render(strings.Join(binds, "  "))
}

func key(k, desc string) string {
	s := DefaultStyles()
	return s.KeyBind.Render(k) + " " + s.KeyDesc.Render(desc)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
