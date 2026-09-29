package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"event-explorer/models"
	"event-explorer/secrets"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// shared client, reuses connections instead of creating one per call
var googleClient = &http.Client{Timeout: 5 * time.Second}

func GetAutoCompletePlaces(ctx context.Context, input string, sessionToken string) ([]models.Suggestion, int, error) {
	payload, err := json.Marshal(map[string]any{
		"input":                input,
		"includedPrimaryTypes": []string{"(cities)"},
		"sessionToken":         sessionToken,
	})
	if err != nil {
		return nil, 400, err
	}

	base := strings.TrimRight(secrets.GetBaseUrls().GoogleBaseURL, "/")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/v1/places:autocomplete", bytes.NewReader(payload))
	if err != nil {
		return nil, 500, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", secrets.GetSecrets().GoogleKey)

	resp, err := googleClient.Do(req)
	if err != nil {
		return nil, 500, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, 500, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, 500, fmt.Errorf("google returned status %d", resp.StatusCode)
	}

	var out models.AutocompleteResp
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, 500, errors.New("invalid response from google")
	}

	suggestions := make([]models.Suggestion, 0, len(out.Suggestions))
	for _, s := range out.Suggestions {
		if s.PlacePrediction != nil {
			suggestions = append(suggestions, models.Suggestion{
				PlaceID: s.PlacePrediction.PlaceID,
				Text:    s.PlacePrediction.Text.Text,
			})
		}
	}
	return suggestions, resp.StatusCode, nil
}

func GetPlaceById(ctx context.Context, placeId string, sessionToken string) (*models.City, int, error) {
	base := strings.TrimRight(secrets.GetBaseUrls().GoogleBaseURL, "/")
	endpoint := base + "/v1/places/" + url.PathEscape(placeId) + "?sessionToken=" + url.QueryEscape(sessionToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, 500, err
	}
	req.Header.Set("X-Goog-Api-Key", secrets.GetSecrets().GoogleKey)
	req.Header.Set("X-Goog-FieldMask", "addressComponents")

	resp, err := googleClient.Do(req)
	if err != nil {
		return nil, 500, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, 500, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, 404, errors.New("place not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, 500, fmt.Errorf("google returned status %d", resp.StatusCode)
	}

	var out models.PlaceResp
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, 500, errors.New("invalid response from google")
	}

	// locality wins, the others are fallbacks only
	var cityName, fallback, countryCode string
	for _, c := range out.AddressComponents {
		for _, t := range c.Types {
			switch t {
			case "locality":
				cityName = c.LongText
			case "postal_town", "administrative_area_level_2":
				if fallback == "" {
					fallback = c.LongText
				}
			case "country":
				countryCode = c.ShortText
			}
		}
	}
	if cityName == "" {
		cityName = fallback
	}
	if cityName == "" || len(countryCode) != 2 {
		return nil, 422, errors.New("could not determine a city for this place")
	}

	city := &models.City{
		City:        cityName,
		CountryCode: countryCode,
	}
	return city, resp.StatusCode, nil
}
