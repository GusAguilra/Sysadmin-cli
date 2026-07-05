package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var cachedStyles *Styles

type Styles struct {
	App      lipgloss.Style
	Header   lipgloss.Style
	Footer   lipgloss.Style
	Content  lipgloss.Style

	Breadcrumb lipgloss.Style

	CatItem      lipgloss.Style
	CatSelected  lipgloss.Style
	CatCount     lipgloss.Style
	CatNum       lipgloss.Style

	CmdItem      lipgloss.Style
	CmdSelected  lipgloss.Style
	CmdDesc      lipgloss.Style
	CmdName      lipgloss.Style
	CmdHeader    lipgloss.Style

	DetailBox    lipgloss.Style
	DetailTitle  lipgloss.Style
	DetailLabel  lipgloss.Style
	DetailValue  lipgloss.Style
	DetailCode   lipgloss.Style
	DetailNotes  lipgloss.Style

	KeyBind      lipgloss.Style
	KeyDesc      lipgloss.Style
	KeySep       lipgloss.Style

	SearchPrompt lipgloss.Style
	SearchText   lipgloss.Style
	SearchPlace  lipgloss.Style
	SearchCursor lipgloss.Style

	Dimmed       lipgloss.Style
	Divider      lipgloss.Style
	Empty        lipgloss.Style
	Error        lipgloss.Style
	Highlight    lipgloss.Style
	ScrollBar    lipgloss.Style
}

func DefaultStyles() Styles {
	if cachedStyles != nil {
		return *cachedStyles
	}

	primary := lipgloss.Color("#00BFFF")
	secondary := lipgloss.Color("#87CEEB")
	text := lipgloss.Color("#E0E0E0")
	dim := lipgloss.Color("#888888")
	border := lipgloss.Color("#333355")
	bg := lipgloss.Color("#1A1A2E")
	sel := lipgloss.Color("#005F87")
	cmdSel := lipgloss.Color("#004466")
	white := lipgloss.Color("#FFFFFF")
	gold := lipgloss.Color("#FFD700")
	green := lipgloss.Color("#00FF87")
	orange := lipgloss.Color("#FFA500")
	red := lipgloss.Color("#FF4444")
	darkBg := lipgloss.Color("#2D2D2D")
	appBg := lipgloss.Color("#16162A")

	s := Styles{
		App: lipgloss.NewStyle().Background(appBg),

		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			Background(bg).
			Padding(0, 1),

		Footer: lipgloss.NewStyle().
			Foreground(dim).
			Background(bg).
			Padding(0, 1),

		Content: lipgloss.NewStyle().Background(bg).Padding(1, 2),

		Breadcrumb: lipgloss.NewStyle().Foreground(primary).Bold(true),

		CatItem: lipgloss.NewStyle().
			Foreground(text).Background(bg).Padding(0, 1),

		CatSelected: lipgloss.NewStyle().
			Foreground(white).Background(sel).Padding(0, 1),

		CatCount: lipgloss.NewStyle().
			Foreground(secondary).Background(bg),

		CatNum: lipgloss.NewStyle().
			Foreground(dim).Background(bg),

		CmdItem: lipgloss.NewStyle().
			Foreground(text).Background(bg).Padding(0, 1),

		CmdSelected: lipgloss.NewStyle().
			Foreground(white).Background(cmdSel).Padding(0, 1),

		CmdDesc: lipgloss.NewStyle().
			Foreground(dim).Background(bg),

		CmdName: lipgloss.NewStyle().
			Foreground(green).Bold(true),

		CmdHeader: lipgloss.NewStyle().
			Foreground(dim),

		DetailBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(border).
			Padding(1, 2).Background(bg),

		DetailTitle: lipgloss.NewStyle().
			Bold(true).Foreground(primary).Background(bg),

		DetailLabel: lipgloss.NewStyle().
			Bold(true).Foreground(secondary).Background(bg),

		DetailValue: lipgloss.NewStyle().
			Foreground(text).Background(bg),

		DetailCode: lipgloss.NewStyle().
			Foreground(gold).Background(darkBg).Padding(0, 1),

		DetailNotes: lipgloss.NewStyle().
			Foreground(orange).Background(bg),

		KeyBind: lipgloss.NewStyle().
			Foreground(gold).Background(bg),

		KeyDesc: lipgloss.NewStyle().
			Foreground(dim).Background(bg),

		KeySep: lipgloss.NewStyle().
			Foreground(border).Background(bg),

		SearchPrompt: lipgloss.NewStyle().
			Foreground(primary).Bold(true).Background(bg),

		SearchText: lipgloss.NewStyle().
			Foreground(text).Background(darkBg),

		SearchPlace: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).Background(darkBg),

		SearchCursor: lipgloss.NewStyle().
			Foreground(primary).Background(darkBg),

		Dimmed: lipgloss.NewStyle().
			Foreground(dim).Background(bg),

		Divider: lipgloss.NewStyle().
			Foreground(border).Background(bg),

		Empty: lipgloss.NewStyle().
			Foreground(dim).Background(bg).Padding(2, 0),

		Error: lipgloss.NewStyle().
			Foreground(red).Background(bg),

		Highlight: lipgloss.NewStyle().
			Foreground(gold),

		ScrollBar: lipgloss.NewStyle().
			Foreground(border).Background(bg),
	}

	cachedStyles = &s
	return s
}
