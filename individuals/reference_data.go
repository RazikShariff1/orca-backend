package individuals

import "main/models"

// Static, hardcoded reference data for halqa, masjid, road and address until
// those get their own tables and CRUD endpoints. Individuals only stores the
// ids (h_id, m_id, r_id, address_id); these maps resolve them for responses.

var halqas = map[int]models.Halqa{
	1: {Id: 1, Name: "Halqa 1"},
	2: {Id: 2, Name: "Halqa 2"},
	3: {Id: 3, Name: "Halqa 3"},
}

var masjids = map[int]models.Masjid{
	1: {Id: 1, Name: "Masjid 1"},
	2: {Id: 2, Name: "Masjid 2"},
	3: {Id: 3, Name: "Masjid 3"},
}

var roads = map[int]models.Road{
	1: {Id: 1, Name: "Road 1"},
	2: {Id: 2, Name: "Road 2"},
	3: {Id: 3, Name: "Road 3"},
}

var addresses = map[int]models.Address{
	1: {Id: 1, RoadId: 1, City: "City 1", State: "State 1", Country: "Country 1"},
	2: {Id: 2, RoadId: 2, City: "City 2", State: "State 2", Country: "Country 2"},
	3: {Id: 3, RoadId: 3, City: "City 3", State: "State 3", Country: "Country 3"},
}
