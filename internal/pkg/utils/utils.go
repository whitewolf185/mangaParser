package utils

// StringMatching - проверяет совпадение строки из списка
func StringMatching(targetStr string, check ...string) bool {
	for _, str := range check {
		if str == targetStr {
			return true
		}
	}
	return false
}
