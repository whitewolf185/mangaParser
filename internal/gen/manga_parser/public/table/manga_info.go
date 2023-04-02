//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var MangaInfo = newMangaInfoTable("public", "manga_info", "")

type mangaInfoTable struct {
	postgres.Table

	//Columns
	ID                postgres.ColumnString
	MangaName         postgres.ColumnString
	URL               postgres.ColumnString
	SourceType        postgres.ColumnString
	LastUpdatedNumber postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type MangaInfoTable struct {
	mangaInfoTable

	EXCLUDED mangaInfoTable
}

// AS creates new MangaInfoTable with assigned alias
func (a MangaInfoTable) AS(alias string) *MangaInfoTable {
	return newMangaInfoTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new MangaInfoTable with assigned schema name
func (a MangaInfoTable) FromSchema(schemaName string) *MangaInfoTable {
	return newMangaInfoTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new MangaInfoTable with assigned table prefix
func (a MangaInfoTable) WithPrefix(prefix string) *MangaInfoTable {
	return newMangaInfoTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new MangaInfoTable with assigned table suffix
func (a MangaInfoTable) WithSuffix(suffix string) *MangaInfoTable {
	return newMangaInfoTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newMangaInfoTable(schemaName, tableName, alias string) *MangaInfoTable {
	return &MangaInfoTable{
		mangaInfoTable: newMangaInfoTableImpl(schemaName, tableName, alias),
		EXCLUDED:       newMangaInfoTableImpl("", "excluded", ""),
	}
}

func newMangaInfoTableImpl(schemaName, tableName, alias string) mangaInfoTable {
	var (
		IDColumn                = postgres.StringColumn("id")
		MangaNameColumn         = postgres.StringColumn("manga_name")
		URLColumn               = postgres.StringColumn("url")
		SourceTypeColumn        = postgres.StringColumn("source_type")
		LastUpdatedNumberColumn = postgres.IntegerColumn("last_updated_number")
		allColumns              = postgres.ColumnList{IDColumn, MangaNameColumn, URLColumn, SourceTypeColumn, LastUpdatedNumberColumn}
		mutableColumns          = postgres.ColumnList{MangaNameColumn, URLColumn, SourceTypeColumn, LastUpdatedNumberColumn}
	)

	return mangaInfoTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:                IDColumn,
		MangaName:         MangaNameColumn,
		URL:               URLColumn,
		SourceType:        SourceTypeColumn,
		LastUpdatedNumber: LastUpdatedNumberColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}