package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"sysadmin-cli/pkg/data"
	"sysadmin-cli/pkg/models"
	"sysadmin-cli/pkg/tui"
)

var (
	cliMode bool
	jsonOut bool
)

var rootCmd = &cobra.Command{
	Use:     "sysadmin-cli [categoria]",
	Short:   "Manual rapido de comandos para Sysadmin",
	Long: `sysadmin-cli es una herramienta interactiva para consultar
comandos utiles de diagnostico y administracion de sistemas.

Usa el modo interactivo (TUI) por defecto o los flags --cli y --json
para salida de texto plano o JSON respectivamente.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		categories, err := data.Load()
		if err != nil {
			return fmt.Errorf("error cargando datos: %w", err)
		}

		var filterCat string
		if len(args) > 0 {
			filterCat = strings.ToLower(args[0])
		}

		if jsonOut {
			return renderJSON(categories, filterCat)
		}

		if cliMode || !isTerminal() {
			return renderCLI(categories, filterCat)
		}

		return runTUI(categories, filterCat)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&cliMode, "cli", false, "Modo terminal (sin TUI)")
	rootCmd.Flags().BoolVar(&jsonOut, "json", false, "Salida en formato JSON")
}

func isTerminal() bool {
	stat, _ := os.Stdout.Stat()
	return (stat.Mode() & os.ModeCharDevice) != 0
}

func runTUI(categories []models.Category, filter string) error {
	filtered := categories
	if filter != "" {
		var cats []models.Category
		for _, cat := range categories {
			if strings.ToLower(cat.Name) == filter {
				cats = append(cats, cat)
				break
			}
		}
		if len(cats) == 0 {
			for _, cat := range categories {
				if strings.Contains(strings.ToLower(cat.Name), filter) {
					cats = append(cats, cat)
				}
			}
		}
		if len(cats) > 0 {
			filtered = cats
		}
	}

	m := tui.New(filtered)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func renderJSON(categories []models.Category, filter string) error {
	filtered := filterCategories(categories, filter)
	if len(filtered) == 0 {
		fmt.Println("[]")
		return nil
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	db := models.Database{Categories: filtered}
	return encoder.Encode(db)
}

func renderCLI(categories []models.Category, filter string) error {
	filtered := filterCategories(categories, filter)
	if len(filtered) == 0 {
		fmt.Println("No se encontraron categorias.")
		return nil
	}

	for i, cat := range filtered {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("=== %s ===\n", cat.Name)
		fmt.Println(cat.Description)
		fmt.Println()

		for _, cmd := range cat.Commands {
			fmt.Printf("  %s\n", cmd.Command)
			fmt.Printf("    %s\n", cmd.Description)
			fmt.Printf("    Ejemplo: $ %s\n", cmd.Example)
			fmt.Println()
		}
	}

	return nil
}

func filterCategories(categories []models.Category, filter string) []models.Category {
	if filter == "" {
		return categories
	}

	filter = strings.ToLower(filter)
	var result []models.Category
	for _, cat := range categories {
		if strings.Contains(strings.ToLower(cat.Name), filter) ||
			strings.Contains(strings.ToLower(cat.Description), filter) {
			var cmds []models.Command
			for _, cmd := range cat.Commands {
				if cmd.Matches(filter) {
					cmds = append(cmds, cmd)
				}
			}
			cat.Commands = cmds
			result = append(result, cat)
		}
	}

	if len(result) == 0 {
		for i := range categories {
			cat := categories[i]
			var cmds []models.Command
			for _, cmd := range cat.Commands {
				if cmd.Matches(filter) {
					cmds = append(cmds, cmd)
				}
			}
			if len(cmds) > 0 {
				cat.Commands = cmds
				result = append(result, cat)
			}
		}
	}

	return result
}
