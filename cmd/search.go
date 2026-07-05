package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"sysadmin-cli/pkg/data"
	"sysadmin-cli/pkg/models"
)

var searchCmd = &cobra.Command{
	Use:   "search <termino>",
	Short: "Busca comandos en todas las categorias",
	Long: `Busca comandos que coincidan con el termino indicado.
La busqueda se realiza en titulos, descripciones, comandos y etiquetas.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.ToLower(args[0])

		categories, err := data.Load()
		if err != nil {
			return fmt.Errorf("error cargando datos: %w", err)
		}

		var results []struct {
			Category string
			Command  models.Command
		}

		for _, cat := range categories {
			for _, c := range cat.Commands {
				if c.Matches(query) {
					results = append(results, struct {
						Category string
						Command  models.Command
					}{cat.Name, c})
				}
			}
		}

		if len(results) == 0 {
			fmt.Fprintf(os.Stdout, "No se encontraron resultados para: %s\n", args[0])
			return nil
		}

		fmt.Fprintf(os.Stdout, "Resultados para '%s' (%d):\n\n", args[0], len(results))
		for _, r := range results {
			fmt.Fprintf(os.Stdout, "  [%s] %s\n", r.Category, r.Command.Command)
			fmt.Fprintf(os.Stdout, "        %s\n", r.Command.Description)
			fmt.Fprintln(os.Stdout)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
