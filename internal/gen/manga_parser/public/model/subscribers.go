//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"github.com/google/uuid"
)

type Subscribers struct {
	ID       int32 `sql:"primary_key"`
	PersonID uuid.UUID
	MangaID  uuid.UUID
}
