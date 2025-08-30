package utils

import (
	"encoding/json"
	"errors"
	"image/color"

	"github.com/HubertasVin/findstr/models"
	"github.com/icza/gox/imagex/colorx"
)

// Checks if a string is in json format.
func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

// Parses style json to style model.
func ParseStyle(s string) (models.Style, error) {
	if s == "" {
		s = "{}"
	}
	if !IsJSON(s) {
		return models.Style{}, errors.New("Error: Style value must be in a valid json format")
	}

	styleDto := models.StyleJson{}
	if err := json.Unmarshal([]byte(s), &styleDto); err != nil {
		return models.Style{}, err
	}

	var res models.Style
	var err error
	if res.MatchFg, err = parseColor(styleDto.MatchFg, nil); err != nil {
		return models.Style{}, err
	}
	if res.MatchBg, err = parseColor(styleDto.MatchBg, &color.RGBA{0, 138, 0, 255}); err != nil {
		return models.Style{}, err
	}
	if res.MatchBold = false; styleDto.MatchBold != nil {
		res.MatchBold = *styleDto.MatchBold
	}

	return res, nil
}

func parseColor(c *string, def *color.RGBA) (color.RGBA, error) {
	if c != nil {
		colorRgba, err := colorx.ParseHexColor(*c)
		if err != nil {
			return color.RGBA{}, errors.New("Error: Parsing hex color got \"" + err.Error() + "\"")
		}
		return colorRgba, nil
	} else {
		if def != nil {
			return *def, nil
		}
		return color.RGBA{255, 255, 255, 255}, nil
	}
}

func BuildJson(fileMatches []models.JsonFileMatch) (string, error) {
	var b []byte
	var err error
	if b, err = json.Marshal(fileMatches); err != nil {
		return "", err
	}
	return string(b), nil
}
