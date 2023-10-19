package utils

func Contains(collection []string, item string) bool {
	for _, element := range collection {
		if element == item {
			return true
		}
	}
	return false
}