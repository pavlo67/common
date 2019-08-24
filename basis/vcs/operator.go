package vcs

import (
	"time"

	"github.com/pavlo67/workshop/basis/common"
	"github.com/pavlo67/workshop/crud"
	"github.com/pavlo67/workshop/dataspace"
)

const TemporaryDomain = "..." // to create "temporary commit" that can be used in GetDiffVCS

type Version [20]byte

type Label struct {
	To      dataspace.ID
	Version Version
	Author  common.ID
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
	Commit(from common.ID, to dataspace.ID, data crud.StringMap, previous Version) (Version, error)

	// Get gets a commit from the repo
	Get(from common.ID, to dataspace.ID, versio Version) (Commit, error)

	// Log gets series of commit labels from the repo
	Log(from common.ID, to dataspace.ID, before *Version, limit uint16) ([]Label, error)

	// Diff gets series of commits titles from the repo
	Diff(from common.ID, to, toCompare dataspace.ID) ([]Diff, error)
}
