package models

type EventResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	URL        string `json:"url"`
	Info       string `json:"info"`
	PleaseNote string `json:"pleaseNote"`

	Dates struct {
		Start struct {
			LocalDate string `json:"localDate"`
			LocalTime string `json:"localTime"`
		} `json:"start"`
		Timezone string `json:"timezone"`
	} `json:"dates"`

	Classifications []struct {
		Segment struct {
			Name string `json:"name"`
		} `json:"segment"`
	} `json:"classifications"`

	Embedded struct {
		Venues []struct {
			Name string `json:"name"`
			City struct {
				Name string `json:"name"`
			} `json:"city"`
			Country struct {
				CountryCode string `json:"countryCode"`
			} `json:"country"`
		} `json:"venues"`
	} `json:"_embedded"`
}

type EventsListResponse struct {
	Embedded struct {
		Events []EventResponse `json:"events"`
	} `json:"_embedded"`
}

type EventCategory struct {
	Name   string
	Events []EventResponse
}

type FetchedResult struct {
	Category string
	Events   []EventResponse
	Status   int
	ErrMsg   string
}
