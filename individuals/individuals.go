package individuals

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"sync"

	"gofr.dev/pkg/gofr"
	gofrHTTP "gofr.dev/pkg/gofr/http"

	"main/models"
)

const selectIndividualQuery = `
SELECT
    i.id, i.name, i.phone, i.email, i.profession_status, i.img, i.meta_data,
    i.created_at, i.updated_at, i.last_met_at,
    i.h_id, i.m_id, i.r_id, i.address_id,
    p.id, p.name
FROM individuals i
JOIN professions p ON p.id = i.profession_id
WHERE i.deleted_at IS NULL`

// rowScanner is implemented by both *sql.Row and *sql.Rows.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanIndividual(row rowScanner) (*models.IndividualResponse, error) {
	var (
		resp              models.IndividualResponse
		phone, email, img sql.NullString
		lastMetAt         sql.NullTime
		metaData          []byte
	)

	err := row.Scan(
		&resp.Id, &resp.Name, &phone, &email, &resp.ProfessionStatus, &img, &metaData,
		&resp.CreatedAt, &resp.UpdatedAt, &lastMetAt,
		&resp.Halqa.Id, &resp.Masjid.Id, &resp.Road.Id, &resp.Address.Id,
		&resp.Profession.Id, &resp.Profession.Name,
	)
	if err != nil {
		return nil, err
	}

	if metaData != nil {
		resp.MetaData = metaData
	}

	if phone.Valid {
		resp.Phone = &phone.String
	}

	if email.Valid {
		resp.Email = &email.String
	}

	if img.Valid {
		resp.Img = &img.String
	}

	if lastMetAt.Valid {
		resp.LastMetAt = &lastMetAt.Time
	}

	return &resp, nil
}

func getIndividualByID(c *gofr.Context, id string) (*models.IndividualResponse, error) {
	row := c.SQL.QueryRowContext(c, selectIndividualQuery+" AND i.id = $1", id)

	resp, err := scanIndividual(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gofrHTTP.ErrorEntityNotFound{Name: "id", Value: id}
		}

		return nil, err
	}

	if err := enrichIndividual(c, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

var fkColumns = []string{"profession_id"}

func mapSQLError(err error) error {
	msg := strings.ToLower(err.Error())

	switch {
	case strings.Contains(msg, "unique constraint") || strings.Contains(msg, "duplicate"):
		return gofrHTTP.ErrorEntityAlreadyExist{}
	case strings.Contains(msg, "foreign key constraint"):
		for _, col := range fkColumns {
			if strings.Contains(msg, col) {
				return gofrHTTP.ErrorInvalidParam{Params: []string{col}}
			}
		}

		return gofrHTTP.ErrorInvalidParam{Params: []string{"profession_id"}}
	default:
		return err
	}
}

// validateReferenceIDs checks h_id, m_id, r_id and address_id against
// account-service, since individuals only stores the ids.
func validateReferenceIDs(c *gofr.Context, req *models.IndividualRequest) error {
	checks := []struct {
		field string
		fetch func() error
	}{
		{"h_id", func() error { _, err := fetchHalqa(c, req.HId); return err }},
		{"m_id", func() error { _, err := fetchMasjid(c, req.MId); return err }},
		{"r_id", func() error { _, err := fetchRoad(c, req.RId); return err }},
		{"address_id", func() error { _, err := fetchAddress(c, req.AddressId); return err }},
	}

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		invalid []string
		errs    []error
	)

	wg.Add(len(checks))

	for _, chk := range checks {
		go func(field string, fetch func() error) {
			defer wg.Done()

			err := fetch()
			if err == nil {
				return
			}

			mu.Lock()
			defer mu.Unlock()

			if errors.Is(err, errAccountServiceNotFound) {
				invalid = append(invalid, field)
			} else {
				errs = append(errs, err)
			}
		}(chk.field, chk.fetch)
	}

	wg.Wait()

	if len(errs) > 0 {
		return errs[0]
	}

	if len(invalid) > 0 {
		return gofrHTTP.ErrorInvalidParam{Params: invalid}
	}

	return nil
}

func Create(c *gofr.Context) (any, error) {
	var req models.IndividualRequest

	if err := c.Bind(&req); err != nil {
		return nil, err
	}

	if err := validate(c, &req); err != nil {
		return nil, err
	}

	const insertQuery = `
INSERT INTO individuals (name, phone, h_id, m_id, r_id, address_id, email, profession_id, profession_status, meta_data, img)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id`

	var id int

	err := c.SQL.QueryRowContext(c, insertQuery,
		req.Name, req.Phone, req.HId, req.MId, req.RId, req.AddressId, req.Email,
		req.ProfessionId, req.ProfessionStatus, req.MetaData, req.Img,
	).Scan(&id)
	if err != nil {
		return nil, mapSQLError(err)
	}

	return getIndividualByID(c, strconv.Itoa(id))
}

func validate(c *gofr.Context, req *models.IndividualRequest) error {
	if req.Name == "" {
		return gofrHTTP.ErrorMissingParam{Params: []string{"name"}}
	}

	if len(req.Name) < 5 {
		return gofrHTTP.ErrorInvalidParam{Params: []string{"name"}}
	}

	if req.Phone != nil && (len(*req.Phone) < 1 || len(*req.Phone) > 10) {
		return gofrHTTP.ErrorInvalidParam{Params: []string{"phone"}}
	}

	if req.Email != nil && *req.Email == "" {
		return gofrHTTP.ErrorInvalidParam{Params: []string{"email"}}
	}

	return validateReferenceIDs(c, req)
}

// GetAll returns every individual that has not been soft-deleted.
func GetAll(c *gofr.Context) (any, error) {
	rows, err := c.SQL.QueryContext(c, selectIndividualQuery+" ORDER BY i.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	individuals := make([]*models.IndividualResponse, 0)

	for rows.Next() {
		individual, err := scanIndividual(rows)
		if err != nil {
			return nil, err
		}

		if err := enrichIndividual(c, individual); err != nil {
			return nil, err
		}

		individuals = append(individuals, individual)
	}

	return individuals, rows.Err()
}

// Get returns a single individual by id.
func Get(c *gofr.Context) (any, error) {
	id := c.PathParam("id")

	if id == "" {
		return nil, gofrHTTP.ErrorMissingParam{Params: []string{"id"}}
	}

	if idInt, err := strconv.Atoi(id); err != nil || idInt < 0 {
		return nil, gofrHTTP.ErrorInvalidParam{Params: []string{"id"}}
	}

	return getIndividualByID(c, id)
}

// Update replaces an existing individual's details.
func Update(c *gofr.Context) (any, error) {
	id := c.PathParam("id")

	var req models.IndividualRequest

	if err := c.Bind(&req); err != nil {
		return nil, err
	}

	if err := validate(c, &req); err != nil {
		return nil, err
	}

	const updateQuery = `
UPDATE individuals
SET name = $1, phone = $2, h_id = $3, m_id = $4, r_id = $5, address_id = $6, email = $7,
    profession_id = $8, profession_status = $9, meta_data = $10, img = $11, updated_at = now()
WHERE id = $12 AND deleted_at IS NULL`

	result, err := c.SQL.ExecContext(c, updateQuery,
		req.Name, req.Phone, req.HId, req.MId, req.RId, req.AddressId, req.Email,
		req.ProfessionId, req.ProfessionStatus, req.MetaData, req.Img, id,
	)
	if err != nil {
		return nil, mapSQLError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, gofrHTTP.ErrorEntityNotFound{Name: "id", Value: id}
	}

	return getIndividualByID(c, id)
}

// Delete soft-deletes an individual by stamping deleted_at.
func Delete(c *gofr.Context) (any, error) {
	id := c.PathParam("id")

	const deleteQuery = `UPDATE individuals SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`

	result, err := c.SQL.ExecContext(c, deleteQuery, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, gofrHTTP.ErrorEntityNotFound{Name: "id", Value: id}
	}

	return "individual successfully deleted with id: " + id, nil
}

// RegisterRoutes wires up the individuals CRUD endpoints on the given app.
func RegisterRoutes(a *gofr.App) {
	a.POST("/individual", Create)
	a.GET("/individual", GetAll)
	a.GET("/individual/{id}", Get)
	a.PUT("/individual/{id}", Update)
	a.DELETE("/individual/{id}", Delete)
}
