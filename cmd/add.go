package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"sysadmin-cli/pkg/models"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Agrega un comando personalizado",
	Long: `Agrega un comando a la base de datos personal del usuario.
Los comandos se guardan en ~/.sysadmin-cli/commands.json y se
fusionan con los comandos incluidos al iniciar la herramienta.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return addCommandInteractive()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func addCommandInteractive() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Agregar nuevo comando personalizado")
	fmt.Println(strings.Repeat("-", 40))

	categoria := prompt(reader, "Categoria: ")
	if categoria == "" {
		return fmt.Errorf("la categoria es obligatoria")
	}

	titulo := prompt(reader, "Titulo: ")
	if titulo == "" {
		return fmt.Errorf("el titulo es obligatorio")
	}

	descripcion := prompt(reader, "Descripcion: ")
	if descripcion == "" {
		return fmt.Errorf("la descripcion es obligatoria")
	}

	comando := prompt(reader, "Comando: ")
	if comando == "" {
		return fmt.Errorf("el comando es obligatorio")
	}

	ejemplo := prompt(reader, "Ejemplo (opcional): ")

	tagsStr := prompt(reader, "Etiquetas (separadas por coma, opcional): ")

	notas := prompt(reader, "Notas (opcional): ")

	var tags []string
	if tagsStr != "" {
		for _, t := range strings.Split(tagsStr, ",") {
			tag := strings.TrimSpace(t)
			if tag != "" {
				tags = append(tags, tag)
			}
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("no se pudo obtener el directorio home: %w", err)
	}

	dir := filepath.Join(home, ".sysadmin-cli")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("no se pudo crear el directorio %s: %w", dir, err)
	}

	path := filepath.Join(dir, "commands.json")

	var db models.Database
	if data, err := os.ReadFile(path); err == nil {
		json.Unmarshal(data, &db)
	}

	found := false
	for i, cat := range db.Categories {
		if strings.EqualFold(cat.Name, categoria) {
			db.Categories[i].Commands = append(db.Categories[i].Commands, models.Command{
				Title:       titulo,
				Description: descripcion,
				Command:     comando,
				Example:     ejemplo,
				Tags:        tags,
				Notes:       notas,
			})
			found = true
			break
		}
	}

	if !found {
		db.Categories = append(db.Categories, models.Category{
			Name:        strings.ToLower(categoria),
			Description: fmt.Sprintf("Comandos personalizados de %s", categoria),
			Commands: []models.Command{
				{
					Title:       titulo,
					Description: descripcion,
					Command:     comando,
					Example:     ejemplo,
					Tags:        tags,
					Notes:       notas,
				},
			},
		})
	}

	output, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return fmt.Errorf("error generando JSON: %w", err)
	}

	if err := os.WriteFile(path, output, 0644); err != nil {
		return fmt.Errorf("error escribiendo %s: %w", path, err)
	}

	fmt.Printf("\nComando guardado en %s\n", path)
	return nil
}

func prompt(reader *bufio.Reader, text string) string {
	fmt.Print(text)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
