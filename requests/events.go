package requests

import (
	"context"
	"encoding/json"
	"event-explorer/models"
	"event-explorer/secrets"
	"net/http"
	"net/url"
	"time"
)

// reusable events client
var eventsClient = &http.Client{Timeout: 5 * time.Second}

func GetConcurrentEvents(ctx context.Context, city string, countryCode string, classificationName string, limit string) ([]models.EventResponse, int, string) {

	baseURL := secrets.GetBaseUrls().TicketmasterBaseURL + "events.json"

	reqURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, 400, "invalid url"
	}

	query := reqURL.Query()
	query.Set("city", city)
	query.Set("countryCode", countryCode)
	query.Set("classificationName", classificationName)
	query.Set("size", limit)
	query.Set("apikey", secrets.GetSecrets().TicketmasterKey)
	reqURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, 500, "can not create new request"
	}

	resp, err := eventsClient.Do(req)
	if err != nil {
		return nil, 500, err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, "ticketmaster returned " + resp.Status
	}

	var result models.EventsListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, http.StatusInternalServerError, "can not parse response"
	}

	events := result.Embedded.Events
	if events == nil {
		events = []models.EventResponse{}
	}

	return events, http.StatusOK, ""
}

func GetEvents(ctx context.Context, city string, countryCode string, limit string) ([]models.EventCategory, int, string) {

	categories := []string{"Music", "Sports"}

	resultChan := make(chan models.FetchedResult, len(categories))

	for _, category := range categories {
		go func(category string) {
			events, status, errMsg := GetConcurrentEvents(ctx, city, countryCode, category, limit)
			resultChan <- models.FetchedResult{Category: category, Events: events, Status: status, ErrMsg: errMsg}
		}(category)
	}

	sections := make([]models.EventCategory, 0, len(categories))

	for i := 0; i < len(categories); i++ {
		result := <-resultChan
		if result.Status != http.StatusOK {
			return nil, result.Status, result.ErrMsg
		}

		sections = append(sections, models.EventCategory{
			Name:   result.Category,
			Events: result.Events,
		})
	}

	close(resultChan)

	return sections, http.StatusOK, ""
}

//get single event details

func GetSingleEvent(ctx context.Context, eventId string) (*models.EventResponse, int, string) {

	baseURL := secrets.GetBaseUrls().TicketmasterBaseURL + "events/" + url.PathEscape(eventId) + ".json"

	reqURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, 400, "invalid url"
	}

	query := reqURL.Query()
	query.Set("apikey", secrets.GetSecrets().TicketmasterKey)
	reqURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, 500, "can not create new request"
	}

	resp, err := eventsClient.Do(req)
	if err != nil {
		return nil, 500, err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, resp.StatusCode, "event not found"
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, "ticketmaster returned " + resp.Status
	}

	var result models.EventResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, http.StatusInternalServerError, "can not parse response"
	}

	return &result, http.StatusOK, ""
}
