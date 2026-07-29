package models

import "time"

// ProfessionRequest is the payload accepted by the create and update endpoints.
type ProfessionRequest struct {
	Name string `json:"name"`
}

// ProfessionResponse is the view returned by the get/list endpoints.
type ProfessionResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HalqaRequest is the payload accepted by the create and update endpoints.
type HalqaRequest struct {
	Name string `json:"name"`
}

// HalqaResponse is the view returned by the get/list endpoints.
type HalqaResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
