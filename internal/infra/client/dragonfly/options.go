package dragonfly

import "time"

type options struct {
	address         string
	password        string
	db              int
	readTimeout     time.Duration
	writeTimeout    time.Duration
	maxRetries      int
	minRetryBackoff time.Duration
	maxRetryBackoff time.Duration
}

type Option interface {
	apply(*options)
}

type addressOption string

func (o addressOption) apply(opts *options) {
	opts.address = string(o)
}

func WithAddress(address string) Option {
	return addressOption(address)
}

type passwordOption string

func (o passwordOption) apply(opts *options) {
	opts.password = string(o)
}

func WithPassword(password string) Option {
	return passwordOption(password)
}

type dbOption int

func (o dbOption) apply(opts *options) {
	opts.db = int(o)
}

func WithDB(db int) Option {
	return dbOption(db)
}

type readTimeoutOption time.Duration

func (o readTimeoutOption) apply(opts *options) {
	opts.readTimeout = time.Duration(o)
}

func WithReadTimeout(timeout time.Duration) Option {
	return readTimeoutOption(timeout)
}

type writeTimeoutOption time.Duration

func (o writeTimeoutOption) apply(opts *options) {
	opts.writeTimeout = time.Duration(o)
}

func WithWriteTimeout(timeout time.Duration) Option {
	return writeTimeoutOption(timeout)
}

type maxRetriesOption int

func (o maxRetriesOption) apply(opts *options) {
	opts.maxRetries = int(o)
}

func WithMaxRetries(maxRetries int) Option {
	return maxRetriesOption(maxRetries)
}

type minRetryBackoffOption time.Duration

func (o minRetryBackoffOption) apply(opts *options) {
	opts.minRetryBackoff = time.Duration(o)
}

func WithMinRetryBackoff(minRetryBackoff time.Duration) Option {
	return minRetryBackoffOption(minRetryBackoff)
}

type maxRetryBackoffOption time.Duration

func (o maxRetryBackoffOption) apply(opts *options) {
	opts.maxRetryBackoff = time.Duration(o)
}

func WithMaxRetryBackoff(maxRetryBackoff time.Duration) Option {
	return maxRetryBackoffOption(maxRetryBackoff)
}
