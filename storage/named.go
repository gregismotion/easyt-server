package storage

import "git.freeself.one/thegergo02/easyt/basic"

type NamedType struct {
	Id   string          `json:"id,omitempty" example:"237e9877-e79b-12d4-a765-321741963000"`
	Name string          `json:"name" example:"height"`
	Type basic.BasicType `json:"type"` // FIXME: can't provide default value in the conventional way, stays 0 for now I guess...
}
