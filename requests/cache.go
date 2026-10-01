package requests

import (
	"fmt"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/cache"
)

const (
	ListTTL   = 5 * time.Minute
	SingleTTL = 10 * time.Minute
)

var (
	ListCache   cache.Cache
	SingleCache cache.Cache
)

func init() {

	var err error
	// "interval" = seconds between expired-item GC sweeps
	ListCache, err = cache.NewCache("memory", `{"interval":60}`)
	if err != nil {
		panic(err)
	}
	SingleCache, err = cache.NewCache("memory", `{"interval":60}`)
	if err != nil {
		panic(err)
	}
}

func ListKey(city, countryCode, limit string) string {
	return fmt.Sprintf("events:list:%s:%s:%s",
		strings.ToLower(strings.TrimSpace(city)),
		strings.ToUpper(strings.TrimSpace(countryCode)),
		limit)
}

func SingleKey(eventId string) string {
	return "events:single:" + eventId
}
