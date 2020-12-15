package api2go

import (
	"context"
)

// APIContextAllocatorFunc to allow custom context implementations
type APIContextAllocatorFunc func(*API) APIContexter

// APIContexter embedding context.Context and requesting two helper functions
type APIContexter interface {
	context.Context
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	Reset(ctx context.Context)
}

// APIContext api2go context for handlers
// It is a mutable implementation of a context.Context
type APIContext struct {
	keys map[string]interface{}
	context.Context
}

// Set a string key value in the context
func (c *APIContext) Set(key string, value interface{}) {
	if c.keys == nil {
		c.keys = make(map[string]interface{})
	}
	c.keys[key] = value
}

// Get a key value from the context
func (c *APIContext) Get(key string) (value interface{}, exists bool) {
	if c.keys != nil {
		value, exists = c.keys[key]
	}
	if !exists {
		value = c.Context.Value(key)
		exists = value != nil
	}
	return
}

// Reset resets all values on Context, making it safe to reuse
func (c *APIContext) Reset(ctx context.Context) {
	c.keys = nil
	c.Context = ctx
}

// Value implements net/context
func (c *APIContext) Value(key interface{}) interface{} {
	if keyAsString, ok := key.(string); ok {
		val, exists := c.Get(keyAsString)
		if exists {
			return val
		}
	}
	return c.Context.Value(key)
}

// Compile time check
var _ APIContexter = &APIContext{}

// ContextQueryParams fetches the QueryParams if Set
func ContextQueryParams(c *APIContext) map[string][]string {
	qp, ok := c.Get("QueryParams")
	if ok == false {
		qp = make(map[string][]string)
		c.Set("QueryParams", qp)
	}
	return qp.(map[string][]string)
}
