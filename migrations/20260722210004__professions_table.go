package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

const createProfessionsTable = `CREATE TABLE IF NOT EXISTS professions (
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    deleted_at timestamp default null
);`

func _professions_table() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(createProfessionsTable)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
