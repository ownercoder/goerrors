package goerrors

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	httpRequestContextKey = "log_http_request"
	fingerprintContextKey = "log_fingerprint"
	tagsContextKey        = "log_tags"
	stackContextKey       = "log_stack"
)

type Context struct {
	ctx context.Context
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Context) Err() error {
	return c.ctx.Err()
}

func (c *Context) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

func (c *Context) AddStack(stack []byte) *Context {
	c.ctx = context.WithValue(c.ctx, stackContextKey, stack)

	return c
}

func (c *Context) AddHTTPRequest(r *http.Request) *Context {
	c.ctx = context.WithValue(c.ctx, httpRequestContextKey, r)

	return c
}

func (c *Context) AddFingerprint(fingerprint []string) *Context {
	c.ctx = context.WithValue(c.ctx, fingerprintContextKey, fingerprint)

	return c
}

func (c *Context) AddTags(tags map[string]string) *Context {
	c.ctx = context.WithValue(c.ctx, tagsContextKey, tags)

	return c
}

func CreateContext(ctx context.Context) *Context {
	return &Context{ctx: ctx}
}

func getEntryStack(entry *logrus.Entry) []byte {
	if entry.Context == nil {
		return nil
	}

	if stack, ok := entry.Context.Value(stackContextKey).([]byte); ok {
		return stack
	}

	return nil
}

func getEntryHTTPRequest(entry *logrus.Entry) *http.Request {
	if entry.Context == nil {
		return nil
	}

	if r, ok := entry.Context.Value(httpRequestContextKey).(*http.Request); ok {
		return r
	}

	return nil
}

func getEntryFingerprint(entry *logrus.Entry) []string {
	if entry.Context == nil {
		return nil
	}

	if fingerprint, ok := entry.Context.Value(fingerprintContextKey).([]string); ok {
		return fingerprint
	}

	return nil
}

func getEntryTags(entry *logrus.Entry) map[string]string {
	if entry.Context == nil {
		return nil
	}

	if tags, ok := entry.Context.Value(tagsContextKey).(map[string]string); ok {
		return tags
	}

	return nil
}
