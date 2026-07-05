package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sysadmin-cli/pkg/models"
)

func Load() ([]models.Category, error) {
	builtin, err := loadBuiltin()
	if err != nil {
		return nil, fmt.Errorf("error cargando comandos incluidos: %w", err)
	}

	custom, err := loadCustom()
	if err != nil {
		return builtin, nil
	}

	merged := mergeCategories(builtin, custom)
	return merged, nil
}

func loadBuiltin() ([]models.Category, error) {
	var db models.Database
	if err := json.Unmarshal(builtInData, &db); err != nil {
		return nil, err
	}
	return db.Categories, nil
}

func loadCustom() ([]models.Category, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(home, ".sysadmin-cli", "commands.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var db models.Database
	if err := json.Unmarshal(data, &db); err != nil {
		return nil, fmt.Errorf("error en el archivo custom %s: %w", path, err)
	}

	return db.Categories, nil
}

func mergeCategories(builtin, custom []models.Category) []models.Category {
	builtinMap := make(map[string]int)
	for i, cat := range builtin {
		builtinMap[cat.Name] = i
	}

	for _, customCat := range custom {
		if idx, ok := builtinMap[customCat.Name]; ok {
			existing := builtin[idx]
			existing.Commands = append(existing.Commands, customCat.Commands...)
			builtin[idx] = existing
		} else {
			builtin = append(builtin, customCat)
		}
	}

	return builtin
}
