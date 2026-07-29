// This is auto-generated file using 'gofr migrate' tool. DO NOT EDIT.
package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

func All() map[int64]migration.Migrate {
	return map[int64]migration.Migrate{

		20260722210004: _professions_table(),
		20260722214416: _individuals_table(),
	}
}
