package content

import "github.com/quangd42/silicon_valley_trail/internal/model"

func DefaultRoute() []model.Location {
	return []model.Location{
		{
			ID:        "san-jose",
			Name:      "San Jose",
			Desc:      "Garage HQ and the start of the run.",
			Latitude:  37.3382,
			Longitude: -121.8863,
		},
		{
			ID:        "santa-clara",
			Name:      "Santa Clara",
			Desc:      "Corporate campuses and investor drive-bys.",
			Latitude:  37.3541,
			Longitude: -121.9552,
		},
		{
			ID:        "sunnyvale",
			Name:      "Sunnyvale",
			Desc:      "Quiet streets with too much ambition per block.",
			Latitude:  37.3688,
			Longitude: -122.0363,
		},
		{
			ID:        "mountain-view",
			Name:      "Mountain View",
			Desc:      "Product gravity gets stronger here.",
			Latitude:  37.3861,
			Longitude: -122.0839,
		},
		{
			ID:        "palo-alto",
			Name:      "Palo Alto",
			Desc:      "Every coffee shop sounds like a pitch deck.",
			Latitude:  37.4419,
			Longitude: -122.1430,
		},
		{
			ID:        "menlo-park",
			Name:      "Menlo Park",
			Desc:      "Sand Hill energy without Sand Hill certainty.",
			Latitude:  37.4530,
			Longitude: -122.1817,
		},
		{
			ID:        "redwood-city",
			Name:      "Redwood City",
			Desc:      "A calm stretch before the last push.",
			Latitude:  37.4852,
			Longitude: -122.2364,
		},
		{
			ID:        "san-mateo",
			Name:      "San Mateo",
			Desc:      "Close enough to feel the pressure now.",
			Latitude:  37.5630,
			Longitude: -122.3255,
		},
		{
			ID:        "south-san-francisco",
			Name:      "South San Francisco",
			Desc:      "The air shifts; the finish is finally visible.",
			Latitude:  37.6547,
			Longitude: -122.4077,
		},
		{
			ID:        "san-francisco",
			Name:      "San Francisco",
			Desc:      "Final pitch territory.",
			Latitude:  37.7749,
			Longitude: -122.4194,
		},
	}
}
