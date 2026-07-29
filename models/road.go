package models

import "time"

// RoadRequest is the payload accepted by the create and update endpoints.
type RoadRequest struct {
	Name string `json:"name"`
}

// RoadResponse is the view returned by the get/list endpoints.
type RoadResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
