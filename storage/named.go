package storage

import "git.freeself.one/thegergo02/easyt/basic"

type NamedType struct {
	Name string `json:"name"`
	Type basic.BasicType `json:"type"`
}
