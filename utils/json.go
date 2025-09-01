package utils

import (
	"encoding/json"

	"github.com/HubertasVin/findstr/models"
)

func BuildJson(fileMatches []models.JsonFileMatch) (string, error) {
	b, err := json.Marshal(fileMatches)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
