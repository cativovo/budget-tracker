package sqlite

import (
	"github.com/huandu/go-sqlbuilder"
)

var builder = sqlbuilder.SQLite

func setMoreIfNotNil[T any](ub *sqlbuilder.UpdateBuilder, f string, v *T) {
	if v != nil {
		ub.SetMore(ub.Assign(f, v))
	}
}

func setUpdatedAt(ub *sqlbuilder.UpdateBuilder) {
	ub.SetMore(ub.Assign("updated_at", sqlbuilder.Raw("CURRENT_TIMESTAMP")))
}
