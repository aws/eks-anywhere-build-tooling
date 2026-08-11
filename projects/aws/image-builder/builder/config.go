package builder

import (
	"encoding/json"
	"log"
	"reflect"
	"sort"
	"strings"
)

// WarnOnUnknownConfigKeys logs a warning for every top-level key in the raw
// config that the given config type has no field for. Those keys are dropped by
// json.Unmarshal and never reach Packer, so without this they fail silently.
//
// This warns rather than errors on purpose. Unknown keys have always been
// accepted, so rejecting them now would break configs that build today.
func WarnOnUnknownConfigKeys(rawConfig []byte, config interface{}, configFlag string) {
	var provided map[string]json.RawMessage
	if err := json.Unmarshal(rawConfig, &provided); err != nil {
		// The typed unmarshal reports malformed JSON with a better message.
		return
	}

	known := knownJSONKeys(reflect.TypeOf(config))

	var unknown []string
	for key := range provided {
		if _, ok := known[key]; !ok {
			unknown = append(unknown, key)
		}
	}
	if len(unknown) == 0 {
		return
	}
	sort.Strings(unknown)

	log.Printf("Warning: ignoring unrecognized key(s) in %s: %s. "+
		"These values are dropped and will not be passed to Packer.",
		configFlag, strings.Join(unknown, ", "))
}

// knownJSONKeys collects the json tag names of t, descending into embedded
// structs because the hypervisor configs compose their fields that way.
func knownJSONKeys(t reflect.Type) map[string]struct{} {
	keys := make(map[string]struct{})
	collectJSONKeys(t, keys)
	return keys
}

func collectJSONKeys(t reflect.Type, keys map[string]struct{}) {
	for t != nil && (t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface) {
		t = t.Elem()
	}
	if t == nil || t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			collectJSONKeys(field.Type, keys)
			continue
		}

		name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		if name == "" || name == "-" {
			continue
		}
		keys[name] = struct{}{}
	}
}
