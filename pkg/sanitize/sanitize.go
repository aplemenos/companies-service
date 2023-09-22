package sanitize

import (
	"encoding/json"

	"github.com/microcosm-cc/bluemonday"
)

// SanitizeJSON sanitizes a body request
func SanitizeJSON(s interface{}) ([]byte, error) {
	sanitize(s)
	return json.MarshalIndent(s, "", "    ")
}

func sanitize(data interface{}) {
	sanitizer := bluemonday.UGCPolicy()

	switch d := data.(type) {
	case map[string]interface{}:
		for k, v := range d {
			switch tv := v.(type) {
			case string:
				d[k] = sanitizer.Sanitize(tv)
			case map[string]interface{}:
				sanitize(tv)
			case []interface{}:
				sanitize(tv)
			case nil:
				delete(d, k)
			}
		}
	case []interface{}:
		if len(d) > 0 {
			switch d[0].(type) {
			case string:
				for i, s := range d {
					d[i] = sanitizer.Sanitize(s.(string))
				}
			case map[string]interface{}:
				for _, t := range d {
					sanitize(t)
				}
			case []interface{}:
				for _, t := range d {
					sanitize(t)
				}
			}
		}
	}
}
