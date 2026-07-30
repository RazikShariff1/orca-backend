package individuals

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"gofr.dev/pkg/gofr"

	"main/models"
)

var errAccountServiceNotFound = errors.New("account-service: not found")

// accountServiceEnvelope mirrors the {"data": ..., "error": ...} shape every
// gofr handler in account-service responds with.
type accountServiceEnvelope[T any] struct {
	Data T `json:"data"`
}

func fetchFromAccountService[T any](c *gofr.Context, path string) (*T, error) {
	resp, err := c.GetHTTPService("account-service").Get(c, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errAccountServiceNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("account-service: unexpected status %d for %s", resp.StatusCode, path)
	}

	var envelope accountServiceEnvelope[T]
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}

	return &envelope.Data, nil
}

// accountServiceNamedEntity fits the {id, name, created_at, updated_at, status}
// shape shared by account-service's h, m and road resources; only id/name matter here.
type accountServiceNamedEntity struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// accountServiceAddress mirrors account-service's flat AddressResponse (road_id,
// not a nested road) before it gets resolved into a models.Address below.
type accountServiceAddress struct {
	Id        int      `json:"id"`
	RoadId    int      `json:"road_id"`
	DoorNo    *string  `json:"door_no"`
	Landmark  *string  `json:"landmark"`
	City      string   `json:"city"`
	State     string   `json:"state"`
	Pincode   *string  `json:"pincode"`
	Country   string   `json:"country"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

func fetchHalqa(c *gofr.Context, id int) (*models.Halqa, error) {
	e, err := fetchFromAccountService[accountServiceNamedEntity](c, fmt.Sprintf("h/%d", id))
	if err != nil {
		return nil, err
	}

	return &models.Halqa{Id: e.Id, Name: e.Name}, nil
}

func fetchMasjid(c *gofr.Context, id int) (*models.Masjid, error) {
	e, err := fetchFromAccountService[accountServiceNamedEntity](c, fmt.Sprintf("m/%d", id))
	if err != nil {
		return nil, err
	}

	return &models.Masjid{Id: e.Id, Name: e.Name}, nil
}

func fetchRoad(c *gofr.Context, id int) (*models.Road, error) {
	e, err := fetchFromAccountService[accountServiceNamedEntity](c, fmt.Sprintf("road/%d", id))
	if err != nil {
		return nil, err
	}

	return &models.Road{Id: e.Id, Name: e.Name}, nil
}

func fetchAddress(c *gofr.Context, id int) (*models.Address, error) {
	raw, err := fetchFromAccountService[accountServiceAddress](c, fmt.Sprintf("address/%d", id))
	if err != nil {
		return nil, err
	}

	road, err := fetchRoad(c, raw.RoadId)
	if err != nil {
		return nil, err
	}

	return &models.Address{
		Id:        raw.Id,
		Road:      *road,
		DoorNo:    raw.DoorNo,
		Landmark:  raw.Landmark,
		City:      raw.City,
		State:     raw.State,
		Pincode:   raw.Pincode,
		Country:   raw.Country,
		Latitude:  raw.Latitude,
		Longitude: raw.Longitude,
	}, nil
}

// enrichIndividual resolves the h_id/m_id/r_id/address_id already scanned onto
// resp.{Halqa,Masjid,Road,Address}.Id against account-service concurrently,
// replacing what used to be static local lookups.
func enrichIndividual(c *gofr.Context, resp *models.IndividualResponse) error {
	fetchers := []func() error{
		func() error {
			halqa, err := fetchHalqa(c, resp.Halqa.Id)
			if err != nil {
				return err
			}

			resp.Halqa = *halqa

			return nil
		},
		func() error {
			masjid, err := fetchMasjid(c, resp.Masjid.Id)
			if err != nil {
				return err
			}

			resp.Masjid = *masjid

			return nil
		},
		func() error {
			road, err := fetchRoad(c, resp.Road.Id)
			if err != nil {
				return err
			}

			resp.Road = *road

			return nil
		},
		func() error {
			address, err := fetchAddress(c, resp.Address.Id)
			if err != nil {
				return err
			}

			resp.Address = *address

			return nil
		},
	}

	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)

	wg.Add(len(fetchers))

	for _, fetch := range fetchers {
		go func(fetch func() error) {
			defer wg.Done()

			if err := fetch(); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}(fetch)
	}

	wg.Wait()

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}
