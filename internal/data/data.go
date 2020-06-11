package data

import "time"

/*
	Structure:
		EntryID: int => {
			name string
			tags []

			type VIDEO | PICTURE | COMIC | OTHER

			size int
			length int
			quality ???

			// optional
			tags_meta {
				tag: {
					auto bool
					times int[] // time(s) in VLC when we added this tag
				}
			}

			times_opened int
			add_date date_time
			creation_date date_time
			access_date date_time
			access_date date_time

			location string // temporary? When do we reset this?
		}

		// INDEX
		tags: int => { EntryID }

		// aliases?
*/

//#region Data
// type keytype string

type Entry struct {
	Name string
	Size int64

	// TODO: video or other filetype-specific info
	Length int64 // TODO: this
	// quality

	Date_Added         time.Time
	Date_Created       time.Time
	Date_Last_Accessed time.Time

	Location     string // TODO: should this be temporary instead?
	Tags         []string
	Times_Opened int64

	HaveManuallyTouched bool
	// Tags_Meta map[string]TagMeta

	// NOTE: now we're going to use the filename as a unique name
	// UID string // NOTE: bad design. This is for internal use to store this object's id
}

type HistoricDB map[string]*HistoricDBEntry
type HistoricDBEntry struct {
	Name     string
	Location string
	Length   interface{}

	Add_Date      string
	Access_Date   int64 // uggg, bad date >_<
	Creation_Date string
	Tags          []string
	Times_Opened  int64

	// TODO: video or other filetype-specific info
	// Length       int
	// quality

	// Date_Added         time.Time
	// Date_Created       time.Time
	// Date_Last_Accessed time.Time

	// Location  string // TODO: should this be temporary instead?
	// Tags      []string
	// Tags_Meta map[string]TagMeta

	// NOTE: now we're going to use the filename as a unique name
	// UID string // NOTE: bad design. This is for internal use to store this object's id
}

// type TagMeta struct {
// 	auto  bool
// 	times []int
// }
type Entries []*Entry
type SimpleFilterFunc func(*Entry) bool
type SimpleSortFunc func(Entries, int, int) bool
type CustomFilterFunc func(Entries) Entries
type CustomSortFunc func(Entries) Entries
type SearchAndFilter struct {
	SimpleFilter SimpleFilterFunc
	SimpleSort   SimpleSortFunc
	CustomFilter CustomFilterFunc
	CustomSort   CustomSortFunc

	InitialIndex string
}

type OpenerConfig map[string]*OpenerFileTypeConfig

type OpenerFileTypeConfig struct {
	Open   string
	Close  string
	Update string
}
