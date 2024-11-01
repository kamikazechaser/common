package httputil

import (
	"encoding/json"
	"net/http"
)

// 10KB is the default maximum body size.
const defaultBodySize = 10 << 10

// BindJSON binds the JSON input to the target struct returning any error encountered.
func BindJSON(
	w http.ResponseWriter,
	req *http.Request,
	target any,
) error {
	req.Body = http.MaxBytesReader(w, req.Body, defaultBodySize)

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(target); err != nil {
		return err
	}

	return nil
}

// JSON writes the response as a JSON with the given HTTP status code and struct.
func JSON(w http.ResponseWriter, httpCode int, res any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	if res == nil {
		return nil
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		return err
	}

	return nil
}
