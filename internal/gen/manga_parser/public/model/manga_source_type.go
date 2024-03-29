//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import "errors"

type MangaSourceType string

const (
	MangaSourceType_Mangalib MangaSourceType = "mangalib"
	MangaSourceType_Mangadex MangaSourceType = "mangadex"
)

func (e *MangaSourceType) Scan(value interface{}) error {
	var enumValue string
	switch val := value.(type) {
	case string:
		enumValue = val
	case []byte:
		enumValue = string(val)
	default:
		return errors.New("jet: Invalid scan value for AllTypesEnum enum. Enum value has to be of type string or []byte")
	}

	switch enumValue {
	case "mangalib":
		*e = MangaSourceType_Mangalib
	case "mangadex":
		*e = MangaSourceType_Mangadex
	default:
		return errors.New("jet: Invalid scan value '" + enumValue + "' for MangaSourceType enum")
	}

	return nil
}

func (e MangaSourceType) String() string {
	return string(e)
}
