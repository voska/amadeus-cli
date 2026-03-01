package output

import (
	"encoding/json"
	"os"
)

func WriteJSON(data any, pretty bool) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	if pretty {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(data)
}
