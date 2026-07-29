package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

const createIndividualsTable = `CREATE TABLE IF NOT EXISTS individuals (
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL,
    phone varchar(15) default null,
    h_id int NOT NULL ,
    m_id int NOT NULL ,
    r_id int NOT NULL ,
    address_id int NOT NULL REFERENCES addresses(id),
    email varchar(100) unique default null,
    profession_id int NOT NULL REFERENCES professions(id),
    profession_status int NOT NULL,
    meta_data jsonb default null,
    img text default null,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    last_met_at timestamp default null,
    deleted_at timestamp default null
);`

func _individuals_table() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(createIndividualsTable)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
