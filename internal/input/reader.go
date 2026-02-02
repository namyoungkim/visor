package input

import (
	"encoding/json"
	"io"
)

// Parse reads JSON from stdin and returns a Session.
// Returns an empty Session on any error (graceful fallback).
func Parse(r io.Reader) *Session {
	var session Session
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&session); err != nil {
		return &Session{}
	}
	return &session
}
