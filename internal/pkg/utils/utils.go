package utils

import (
	"os"

	"github.com/pkg/errors"
)

// StringMatching - проверяет совпадение строки из списка
func StringMatching(targetStr string, check ...string) bool {
	for _, str := range check {
		if str == targetStr {
			return true
		}
	}
	return false
}

func FolderController(folderPathToSave string) error {
	err := os.MkdirAll(folderPathToSave, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "unexpected error in folderController")
	}
	return nil
}
