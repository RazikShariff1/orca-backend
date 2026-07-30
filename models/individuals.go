package models

import (
	"encoding/json"
	"time"
)

// IndividualRequest is the payload accepted by the create and update endpoints.
type IndividualRequest struct {
	Name             string          `json:"name"`
	Phone            *string         `json:"phone"`
	HId              int             `json:"h_id"`
	MId              int             `json:"m_id"`
	RId              int             `json:"r_id"`
	AddressId        int             `json:"address_id"`
	Email            *string         `json:"email"`
	ProfessionId     int             `json:"profession_id"`
	ProfessionStatus int             `json:"profession_status"`
	MetaData         json.RawMessage `json:"meta_data"`
	Img              *string         `json:"img"`
}

// IndividualResponse is the enriched view returned by the get/list endpoints,
// with foreign keys resolved to their referenced records.
type IndividualResponse struct {
	Id               int             `json:"id"`
	Name             string          `json:"name"`
	Phone            *string         `json:"phone"`
	Email            *string         `json:"email"`
	ProfessionStatus int             `json:"profession_status"`
	Img              *string         `json:"img"`
	MetaData         json.RawMessage `json:"meta_data"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	LastMetAt        *time.Time      `json:"last_met_at"`
	Halqa            Halqa           `json:"halqa"`
	Masjid           Masjid          `json:"masjid"`
	Road             Road            `json:"road"`
	Profession       Profession      `json:"profession"`
	Address          Address         `json:"address"`
}

type Halqa struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Masjid struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Road struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Profession struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Address struct {
	Id        int      `json:"id"`
	Road      Road     `json:"road"`
	DoorNo    *string  `json:"door_no"`
	Landmark  *string  `json:"landmark"`
	City      string   `json:"city"`
	State     string   `json:"state"`
	Pincode   *string  `json:"pincode"`
	Country   string   `json:"country"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

// MasjidRequest is the payload accepted by the create and update endpoints.
type MasjidRequest struct {
	Name string `json:"name"`
	HId  int    `json:"h_id"`
}

// MasjidResponse is the view returned by the get/list endpoints.
type MasjidResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Halqa     Halqa     `json:"halqa"`
}
