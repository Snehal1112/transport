package client

import "context"

// Options the self referenceing function
// to set the config
type Options func(c *Connect)

// WithDatabase self referenceing function to set the database name.
func WithDatabase(database string) Options {
	return func(c *Connect) {
		c.databaseName = database
	}
}

// WithCtx self referenceing function to set the context
func WithCtx(ctx context.Context) Options {
	return func(c *Connect) {
		c.ctx = ctx
	}
}

// WithURL self referenceing function to set the url
func WithURL(url string) Options {
	return func(c *Connect) {
		c.url = url
	}
}

// WithLogLevel self referenceing function to
// set the log level for the logrus.
func WithLogLevel(level string) Options {
	return func(c *Connect) {
		c.loglevel = level
	}
}
