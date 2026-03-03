// Package contextkey provides strongly-typed keys for storing and retrieving values in context.Context.
package contextkey

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "eneba ctx key " + k.name
}

// CID is the context key used to store and retrieve correlation IDs.
var CID = &contextKey{"cid"}
