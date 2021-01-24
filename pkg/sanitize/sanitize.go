package sanitize

import (
	"bytes"

	"github.com/microcosm-cc/bluemonday"

	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/serializer"
)

var sanitizer *bluemonday.Policy

func init() {
	sanitizer = bluemonday.UGCPolicy()
}

func SanitizeJSON(s []byte) ([]byte, error) {
	d := serializer.JSON().NewDecoder(bytes.NewReader(s))
	d.UseNumber()
	var i interface{}
	err := d.Decode(&i)
	if err != nil {
		return nil, err
	}
	sanitize(i)
	return serializer.JSON().MarshalIndent(i, "", "    ")
}

func sanitize(data interface{}) {
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
