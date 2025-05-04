package dragonfly

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultMaxRetries      = 3
	_defaultMinRetryBackoff = 500 * time.Millisecond
	_defaultMaxRetryBackoff = 5 * time.Second
)

type Client struct {
	redis *redis.Client
}

// New creates a new cache client.
func New(opts ...Option) (*Client, error) {
	opt := options{
		readTimeout:     _defaultReadTimeout,
		writeTimeout:    _defaultWriteTimeout,
		maxRetries:      _defaultMaxRetries,
		minRetryBackoff: _defaultMinRetryBackoff,
		maxRetryBackoff: _defaultMaxRetryBackoff,
	}

	for _, o := range opts {
		o.apply(&opt)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:            opt.address,
		Password:        opt.password,
		DB:              opt.db,
		ReadTimeout:     opt.readTimeout,
		WriteTimeout:    opt.writeTimeout,
		MaxRetries:      opt.maxRetries,
		MinRetryBackoff: opt.minRetryBackoff,
		MaxRetryBackoff: opt.maxRetryBackoff,
	})

	return &Client{
		redis: redisClient,
	}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.redis.Ping(ctx).Err()
}

func (c *Client) Close() error {
	return c.redis.Close()
}

func (c *Client) Pipeline() redis.Pipeliner {
	return c.redis.Pipeline()
}

func (c *Client) GetBytes(ctx context.Context, key string) ([]byte, error) {
	return c.redis.Get(ctx, key).Bytes()
}

func (c *Client) GeoAdd(ctx context.Context, key string, locations ...*redis.GeoLocation) error {
	return c.redis.GeoAdd(ctx, key, locations...).Err()
}

func (c *Client) GeoSearchLocation(ctx context.Context, key string, query *redis.GeoSearchLocationQuery) ([]redis.GeoLocation, error) {
	return c.redis.GeoSearchLocation(ctx, key, query).Result()
}

func (c *Client) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.redis.SMembers(ctx, key).Result()
}

func (c *Client) GeoRadius(ctx context.Context, key string, lng, lat, radius float64) ([]string, error) {
	cmd := c.redis.Do(ctx, "GEOSEARCH", key, "FROMLONLAT", lng, lat, "BYRADIUS", radius, "m")

	if err := cmd.Err(); err != nil {
		return nil, fmt.Errorf("cmd.Err: %w", err)
	}

	results, err := cmd.StringSlice()
	if err != nil {
		return nil, fmt.Errorf("cmd.StringSlice: %w", err)
	}

	return results, nil
}

func (c *Client) GeoRadiusWithDist(ctx context.Context, key string, lng, lat, radius float64) ([]redis.GeoLocation, error) {
	cmd := c.redis.Do(ctx, "GEOSEARCH", key, "FROMLONLAT", lng, lat, "BYRADIUS", radius, "m", "WITHCOORD", "WITHDIST")

	if err := cmd.Err(); err != nil {
		return nil, fmt.Errorf("cmd.Err: %w", err)
	}

	val := cmd.Val()
	// TODO val.([][]any) ???
	items, ok := val.([]any)
	if !ok {
		return nil, fmt.Errorf("val is not a slice")
	}

	results := []redis.GeoLocation{}

	// TODO change errors message
	for _, item := range items {
		arr, ok := item.([]any)
		if !ok {
			return nil, fmt.Errorf("item is not a slice")
		}

		if len(arr) < 3 {
			return nil, fmt.Errorf("item is not a slice")
		}

		var (
			name string
			dist float64
			lon  float64
			lat  float64
		)

		name, ok = arr[0].(string)
		if !ok {
			return nil, fmt.Errorf("item is not a slice")
		}
		dist, ok = arr[1].(float64)
		if !ok {
			return nil, fmt.Errorf("item is not a slice")
		}

		coord, ok := arr[2].([]any)
		if !ok {
			return nil, fmt.Errorf("item is not a slice")
		}

		lon, ok = coord[0].(float64)
		if !ok {
			return nil, fmt.Errorf("item is not a slice")
		}
		lat, ok = coord[1].(float64)
		if !ok {
			return nil, fmt.Errorf("item is not a slice")
		}

		results = append(results, redis.GeoLocation{
			Name:      name,
			Longitude: lon,
			Latitude:  lat,
			Dist:      dist,
		})
	}

	return results, nil
}

func (c *Client) GetTransportInRadiusDirect(ctx context.Context, key string, lng, lat, radius float64) ([]redis.GeoLocation, error) {
	return c.GeoRadiusWithDist(ctx, key, lng, lat, radius)
}

func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.redis.Set(ctx, key, value, expiration).Err()
}

func (c *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.redis.HGetAll(ctx, key).Result()
}

func (c *Client) HSet(ctx context.Context, key string, values ...any) error {
	return c.redis.HSet(ctx, key, values...).Err()
}

func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.redis.Expire(ctx, key, expiration).Err()
}
