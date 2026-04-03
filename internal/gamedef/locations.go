package gamedef

import "github.com/quangd42/silicon_valley_trail/internal/model"

func DefaultRoute() []model.Location {
	return []model.Location{
		{
			ID:        "san-jose",
			Name:      "San Jose",
			Desc:      "Garage HQ and the start of the run.",
			Latitude:  37.3394,
			Longitude: -121.8939,
		},
		{
			ID:        "santa-clara",
			Name:      "Santa Clara",
			Desc:      "Corporate campuses and investor drive-bys.",
			Latitude:  37.354,
			Longitude: -121.954,
		},
		{
			ID:        "sunnyvale",
			Name:      "Sunnyvale",
			Desc:      "Quiet streets with too much ambition per block.",
			Latitude:  37.3689,
			Longitude: -122.0353,
		},
		{
			ID:        "mountain-view",
			Name:      "Mountain View",
			Desc:      "Product gravity gets stronger here.",
			Latitude:  37.3861,
			Longitude: -122.0828,
		},
		{
			ID:        "palo-alto",
			Name:      "Palo Alto",
			Desc:      "Every coffee shop sounds like a pitch deck.",
			Latitude:  37.4419,
			Longitude: -122.1419,
		},
		{
			ID:        "menlo-park",
			Name:      "Menlo Park",
			Desc:      "Sand Hill energy without Sand Hill certainty.",
			Latitude:  37.4539,
			Longitude: -122.1811,
		},
		{
			ID:        "redwood-city",
			Name:      "Redwood City",
			Desc:      "A calm stretch before the last push.",
			Latitude:  37.4853,
			Longitude: -122.2353,
		},
		{
			ID:        "san-mateo",
			Name:      "San Mateo",
			Desc:      "Close enough to feel the pressure now.",
			Latitude:  37.5631,
			Longitude: -122.3244,
		},
		{
			ID:        "south-san-francisco",
			Name:      "South San Francisco",
			Desc:      "The air shifts; the finish is finally visible.",
			Latitude:  37.6547,
			Longitude: -122.4067,
		},
		{
			ID:        "san-francisco",
			Name:      "San Francisco",
			Desc:      "Final meeting territory.",
			Latitude:  37.775,
			Longitude: -122.4183,
		},
	}
}
