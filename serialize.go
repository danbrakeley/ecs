package ecs

import (
	"bytes"
	"encoding/json"
)

// SerializeToJSON is a default JSON serializer that just uses encoding/json
func SerializeToJSON(buf *bytes.Buffer, thing interface{}) error {
	// TODO: if thing is a struct, iterate over fields, looking at tags.
	// for each "json" tag, see if that field is a Serializer, if so call Serialize, else call json.Marshal

	b, err := json.Marshal(thing)
	if err != nil {
		return err
	}
	_, err = buf.Write(b)
	if err != nil {
		return err
	}
	return nil
}
