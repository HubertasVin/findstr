package mappers

import (
	"context"

	"github.com/HubertasVin/findstr/models"
)

func MapChanToJsonFile(ctx context.Context, input <-chan models.FileMatch) []models.JsonFileMatch {
	var res []models.JsonFileMatch
	for {
		select {
		case <-ctx.Done():
			return res
		case fm, ok := <-input:
			if !ok {
				return res
			}
			lm := MapFileToLineContents(fm)
			jfm := models.JsonFileMatch{
				FileName:       fm.File,
				MatchedContent: lm,
			}
			res = append(res, jfm)
		}
	}
}

func MapFileToLineContents(intput models.FileMatch) []models.LineContent {
	res := []models.LineContent{}
	for _, ln := range intput.MatchLineNums {
		lm := models.LineContent{
			LineNumber: ln + 1,
			Content:    intput.FileContent[ln],
		}
		res = append(res, lm)
	}
	return res
}
