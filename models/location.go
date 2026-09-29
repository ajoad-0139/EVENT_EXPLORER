package models

import "net/http"

type Suggestion struct {
	PlaceID string `json:"placeId"`
	Text    string `json:"text"`
}

type City struct {
	City        string `json:"city"`
	CountryCode string `json:"countryCode"`
}

type AutocompleteResp struct {
	Suggestions []struct {
		PlacePrediction *struct {
			PlaceID string `json:"placeId"`
			Text    struct {
				Text    string `json:"text"`
				Matches []struct {
					EndOffset int `json:"endOffset"`
				} `json:"matches"`
			} `json:"text"`
		} `json:"placePrediction"`
	} `json:"suggestions"`
}

type GoogleClient struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

type AddressComponent struct {
	LongText  string   `json:"longText"`
	ShortText string   `json:"shortText"`
	Types     []string `json:"types"`
}

type PlaceResp struct {
	AddressComponents []AddressComponent `json:"addressComponents"`
}
