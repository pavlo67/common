package vcs

import (
	"time"

	"github.com/pavlo67/associatio/auth"
	"github.com/pavlo67/associatio/crud"
	"github.com/pavlo67/associatio/dataspace"
)

const TemporaryDomain = "..." // to create "temporary commit" that can be used in GetDiffVCS

type Version [20]byte

type Label struct {
	To      dataspace.ID
	Version Version
	Author  auth.ID
	Time    time.Time
}

type Commit struct {
	Label Label
	Data  crud.StringMap
}

type Diff struct {
	Place    string
	Position uint64
	Deletion string
	Addition string
}

//type Qualifier string
//
//const Eq Qualifier = ""
//const Lt Qualifier = "<"
//const Gt Qualifier = ">"
//
//type VersionDefinition struct {
//	Version
//	Qualifier
//}

type Operator interface {

	// Commit tries to put the commit into the repo
	Commit(from auth.ID, to dataspace.ID, data crud.StringMap, previous Version) (Version, error)

	// Get gets a commit from the repo
	Get(from auth.ID, to dataspace.ID, versio Version) (Commit, error)

	// Log gets series of commit labels from the repo
	Log(from auth.ID, to dataspace.ID, before *Version, limit uint16) ([]Label, error)

	// Diff gets series of commits titles from the repo
	Diff(from auth.ID, to, toCompare dataspace.ID) ([]Diff, error)
}
