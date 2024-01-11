package jsonsecurity

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Masker struct allows you to use a single Masker object to manipulate json masking.
// cfg stores the config that is applied for all masking operations executed via Mask method.
type Masker struct {
	cfg Config

	lowercaseTriggers map[string][]TriggerOpts
}

// NewMasker returns a new Masker instance.
func NewMasker(cfg Config) (*Masker, error) {
	lowercaseTriggers := make(map[string][]TriggerOpts, 0)
	for k, v := range cfg.Triggers {
		lowerK := strings.ToLower(k)
		v.original = k

		if _, ok := lowercaseTriggers[lowerK]; !ok {
			lowercaseTriggers[lowerK] = []TriggerOpts{v}
		} else {
			lowercaseTriggers[lowerK] = append(lowercaseTriggers[lowerK], v)
		}
	}

	return &Masker{
		cfg:               cfg,
		lowercaseTriggers: lowercaseTriggers,
	}, nil
}

// Mask masks data given in parameters using specified config.
//
// Example:
//
// Data = `{"password": "qwerty123", "email": "example@example.com"}`
//
// Using PASSWORD label for "password" and EMAIL for "email" we will reach the next result:
//
// Output = `{"password": "*********", "email": "e******@example.com"}`
func (m *Masker) Mask(data string) (string, error) {
	var dataAny any

	err := json.Unmarshal([]byte(data), &dataAny)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal initial data: %+v", err)
	}

	result, err := m.walkthrough(dataAny, 0)
	if err != nil {
		return "", fmt.Errorf("walkthrough error: %+v", err)
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("fail while marshalling result map: %+v", err)
	}

	return string(bytes), nil
}

// walkthrough falls down to the json structure considering every key-value pair may be represented as map[string]interface{}
// It gets trigger options in config using key and replaces initial value using the value from walkthrough result.
func (m *Masker) walkthrough(layer interface{}, currentDepth int) (interface{}, error) {
	if currentDepth > m.cfg.MaxDepth {
		return nil, fmt.Errorf("max recursion depth reached: %d", m.cfg.MaxDepth)
	}

	switch t := layer.(type) {
	case map[string]interface{}:
		for k, v := range t {
			var err error

			triggerOpt, ok := m.getTriggerOpts(k)
			if !ok {
				t[k], err = m.walkthrough(v, currentDepth+1)
				if err != nil {
					return nil, fmt.Errorf("fail while walking through key '%s', depth '%d': %+v", k, currentDepth, err)
				}
				continue
			}

			err = mask(t, k, triggerOpt)
			if err != nil {
				return nil, fmt.Errorf("fail while masking key '%s', depth '%d': %+v", k, currentDepth, err)
			}
		}
	case []interface{}:
		for i, v := range t {
			var err error

			t[i], err = m.walkthrough(v, currentDepth+1)
			if err != nil {
				return nil, fmt.Errorf("fail while masking array index '%d', depth '%d': %+v", i, currentDepth, err)
			}
		}
	default:
	}
	return layer, nil
}

// getTriggerOpts returns TriggerOpts for specified label-string.
func (m *Masker) getTriggerOpts(key string) (TriggerOpts, bool) {
	lowerKey := strings.ToLower(key)
	if triggers, ok := m.lowercaseTriggers[lowerKey]; ok && len(triggers) > 0 {
		for _, trigger := range m.lowercaseTriggers[lowerKey] {
			if trigger.CaseSensitive && trigger.original == key {
				return trigger, true
			} else if !trigger.CaseSensitive {
				return trigger, true
			}
		}
	}

	return TriggerOpts{}, false
}

// mask masks the key-value pair located in dict using provided TriggerOpts.
func mask(dict map[string]interface{}, key string, opts TriggerOpts) error {
	if !opts.ShouldAppear {
		delete(dict, key)
		return nil
	}

	var err error

	_mapLabelsMu.RLock()
	masker, ok := _mapLabels[opts.MaskMethod]
	_mapLabelsMu.RUnlock()

	if !ok {
		masker = defaultMasker
	}

	ret, err := masker(key, dict[key])
	if err != nil {
		return err
	}

	dict[ret.Key] = ret.Value
	return nil
}

// MaskResult represents the result of masking.
type MaskResult struct {
	// Key masked key. (Maybe useful in future)
	Key string

	// Value masked value.
	Value interface{}
}

// MaskerFunc masks given key-value pair and returns MaskResult structure which defines further behavior of masking algorithm.
type MaskerFunc func(key string, value interface{}) (MaskResult, error)

func defaultMasker(key string, value interface{}) (MaskResult, error) {
	return MaskResult{Key: key, Value: value}, nil
}
