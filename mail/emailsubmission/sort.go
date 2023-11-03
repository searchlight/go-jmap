package emailsubmission

import "git.sr.ht/~rockorager/go-jmap"

type SortComparator struct {
	Property string `json:"property,omitempty"`

	IsAscending bool `json:"isAscending,omitempty"`

	Collation jmap.CollationAlgo `json:"collation,omitempty"`
}
