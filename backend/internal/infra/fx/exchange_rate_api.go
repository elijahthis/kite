package fx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/elijahthis/kite/internal/domain"
)

const defaultAPIURL = "https://open.er-api.com/v6/latest/"

type rateCacheItem struct {
	rates     map[string]float64
	expiresAt time.Time
}

// fetches rates and caches them
type ERAPIProvider struct {
	client *http.Client
	mu     sync.RWMutex
	cache  map[string]rateCacheItem
	ttl    time.Duration
}

func NewERAPIProvider(ttl time.Duration) *ERAPIProvider {
	return &ERAPIProvider{
		client: &http.Client{Timeout: 5 * time.Second},
		cache:  make(map[string]rateCacheItem),
		ttl:    ttl,
	}
}

func (p *ERAPIProvider) GetRate(ctx context.Context, source, target domain.Currency) (float64, error) {
	sourceStr := source.String()
	targetStr := target.String()

	p.mu.RLock()
	item, exists := p.cache[sourceStr]
	p.mu.RUnlock()

	if exists && time.Now().Before(item.expiresAt) {
		if rate, ok := item.rates[targetStr]; ok {
			return rate, nil
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, defaultAPIURL+sourceStr, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create fx request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", domain.ErrRateFetch, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("%w: received status %d", domain.ErrRateFetch, resp.StatusCode)
	}

	var result struct {
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode fx response: %w", err)
	}

	p.mu.Lock()
	p.cache[sourceStr] = rateCacheItem{
		rates:     result.Rates,
		expiresAt: time.Now().Add(p.ttl),
	}
	p.mu.Unlock()

	rate, ok := result.Rates[targetStr]
	if !ok {
		return 0, fmt.Errorf("target currency %s not supported by provider", targetStr)
	}

	return rate, nil
}
