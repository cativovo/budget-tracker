package sqlite

import (
	"github.com/huandu/go-sqlbuilder"
)

func setMoreIfNotNil[T any](ub *sqlbuilder.UpdateBuilder, f string, v *T) {
	if v != nil {
		ub.SetMore(ub.Assign(f, v))
	}
}
