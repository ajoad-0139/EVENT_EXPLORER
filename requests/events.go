package requests

import (
	"context"
	"encoding/json"
	"event-explorer/models"
	"event-explorer/secrets"
	"net/http"
	"net/url"
	"sync"
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
	results := make([]models.FetchedResult, len(categories))

	var wg sync.WaitGroup
	for i, category := range categories {
		wg.Add(1)
		go func(i int, category string) {
			defer wg.Done()
			events, status, errMsg := GetConcurrentEvents(ctx, city, countryCode, category, limit)
			results[i] = models.FetchedResult{Events: events, Status: status, ErrMsg: errMsg} // each goroutine writes only its own index
		}(i, category)
	}
	wg.Wait()

	// group by category, keeping the order of `categories`
	sections := make([]models.EventCategory, 0, len(categories))
	for i, r := range results {
		if r.Status != http.StatusOK {
			return nil, r.Status, r.ErrMsg
		}
		sections = append(sections, models.EventCategory{
			Name:   categories[i],
			Events: r.Events,
		})
	}

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
