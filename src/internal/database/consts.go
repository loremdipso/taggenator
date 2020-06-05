package database

const (
	settings_filename = "taggenator_settings.json"
	database_filename = "data.db"
)

const (
	IndexTags             = "tags"
	IndexCreationDate     = "creation_date"
	IndexModificationDate = "modification_date"
	IndexName             = "name"
)

// const (
// EntryID_Key string = "ID_Latest_Entry"
// )
const (
	id_entry_prefix string = "entry_"
	id_tag_prefix   string = "tag_"
	KeyEntries      string = id_entry_prefix + "*"
	KeyTags         string = id_tag_prefix + "*"
)
