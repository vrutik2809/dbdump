package utils

import (
	"fmt"
	"strings"
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
	return fmt.Sprintf("{{ green \"dumping %s:\" }} {{ cyan \"%s\" }} {{ bar . \"[\" \"=\" \">\" \"_\" \"]\"}} {{percent .}}", entityType, entity)
}
