package main

import (
	"bytes"
	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"net/http"
)

// Utitility to convert JSON object in body
func UnmarshalRequest(c context.Context, r *http.Request, value interface{}) error {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(r.Body)
	err := json.Unmarshal(buffer.Bytes(), value)
	if err != nil {
		log.Errorf(c, "Error while decoing JSON: %v", err)
		log.Infof(c, "JSON: %v", buffer.String())
		return err
	}
	return nil
}
