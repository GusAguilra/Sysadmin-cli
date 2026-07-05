package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sysadmin-cli/pkg/data"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Muestra todas las categorias disponibles",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		categories, err := data.Load()
		if err != nil {
			return fmt.Errorf("error cargando datos: %w", err)
		}

		fmt.Fprintf(os.Stdout, "Categorias disponibles (%d):\n\n", len(categories))
		for _, cat := range categories {
			fmt.Fprintf(os.Stdout, "  %-12s  %s\n", cat.Name, cat.Description)
		}
		fmt.Fprintln(os.Stdout)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
