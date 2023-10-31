package utils

import (
	"fmt"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

func Contains(collection []string, item string) bool {
	for _, element := range collection {
		if element == item {
			return true
		}
	}
	return false
}

func GetName(filename string) string {
	return strings.Split(filename, ".")[0]
}

func GetBarTemplate(entityType string, entity string) string {
	return fmt.Sprintf("{{ \"dumping %s:\" }} {{ \"%s\" }} {{ bar . \"[\" \"=\" \">\" \"_\" \"]\"}} {{percent .}}", entityType, entity)
}

func GetBars(arr []string, entityType string) []*pb.ProgressBar {
	bars := []*pb.ProgressBar{}
	for _, entity := range arr {
		bar := pb.New(0).SetTemplateString(GetBarTemplate(entityType, entity)).SetMaxWidth(80)
		bars = append(bars, bar)
	}
	return bars
}
