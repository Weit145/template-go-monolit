package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const MaxJSONBodySize int64 = 64 << 10

func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, MaxJSONBodySize)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return fmt.Errorf("decode JSON body: %w", err)
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return fmt.Errorf("decode JSON body: multiple JSON values")
		}
		return fmt.Errorf("decode trailing JSON body: %w", err)
	}

	return nil
}
