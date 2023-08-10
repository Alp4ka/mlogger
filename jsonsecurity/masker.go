package jsonsecurity

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

var (
	_globalMaskerMu = sync.RWMutex{}
	_globalMasker   *Masker
)

// Masker struct allows you to use a single Masker object to manipulate json masking.
// cfg stores the config that is applied for all masking operations executed via Mask method.
type Masker struct {
	cfg Config

	lowercaseTriggers map[string][]TriggerOpts
}

func GlobalMasker() *Masker {
	var m *Masker

	_globalMaskerMu.RLock()
	m = _globalMasker
	_globalMaskerMu.RUnlock()

	return m
}

func ReplaceGlobals(masker *Masker) {
	_globalMaskerMu.Lock()
	_globalMasker = masker
	_globalMaskerMu.Unlock()
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

// Deprecated: Mask masks data given in parameters using specified config.
//
// Example:
//
// Data = `{"password": "qwerty123", "email": "example@example.com"}`
//
// Using PASSWORD label for "password" and EMAIL for "email" we will reach the next result:
//
// Output = `{"password": "*********", "email": "e******@example.com"}`
// TODO(Gorkovets Roman): Definitely needs rework due to multiple serialization and deserialization. Inefficient af.
// TODO(Gorkovets Roman): Decompose.
func (m *Masker) Mask(data []byte) ([]byte, error) {
	mapData := make(map[string]interface{}, 0)

	err := json.Unmarshal(data, &mapData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %+v", err)
	}

	// Collecting lowercase keys and making copy of initial map into a resulting one.
	mapResult := make(map[string]interface{}, len(mapData))
	mapLCData := make(map[string][]modifiedKey, len(mapData))
	for k, v := range mapData {
		lowerK := strings.ToLower(k)
		if _, ok := mapLCData[lowerK]; !ok {
			mapLCData[lowerK] = []modifiedKey{{Modified: lowerK, Original: k}}
		} else {
			mapLCData[lowerK] = append(mapLCData[lowerK], modifiedKey{Modified: lowerK, Original: k})
		}
		mapResult[k] = v
	}

	// Iterating over triggers.
	var (
		triggerErr error
		res        MaskResult
	)
	for trigger, opts := range m.cfg.Triggers {
		if opts.ShouldAppear {
			if opts.CaseSensitive {
				if val, ok := mapData[trigger]; ok {
					res, triggerErr = maskLabel(opts.MaskMethod, trigger, val)
					if triggerErr != nil {
						return nil, fmt.Errorf("fail while masking trigger '%s': %+v", trigger, err)
					}
					mapResult[res.Key] = res.Value
				}
			} else {
				if modVals, ok := mapLCData[strings.ToLower(trigger)]; ok {
					for _, modVal := range modVals {
						res, triggerErr = maskLabel(opts.MaskMethod, modVal.Original, mapData[modVal.Original])
						if triggerErr != nil {
							return nil, fmt.Errorf("fail while masking trigger '%s': %+v", trigger, err)
						}
						mapResult[res.Key] = res.Value
					}
				}
			}
		} else {
			if opts.CaseSensitive {
				delete(mapResult, trigger)
			} else {
				if modVals, ok := mapLCData[strings.ToLower(trigger)]; ok {
					for _, modVal := range modVals {
						delete(mapResult, modVal.Original)
					}
				}
			}
		}
	}

	bytes, err := json.Marshal(mapResult)
	if err != nil {
		return nil, fmt.Errorf("fail while marshalling resulting map: %+v", err)
	}

	return bytes, nil
}

func (m *Masker) Mask2(data []byte) ([]byte, error) {
	mapData := make(map[string]interface{}, 0)

	err := json.Unmarshal(data, &mapData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal initial data: %+v", err)
	}

	//mapResult := make(map[string]interface{}, len(mapData))

	result, err := m.walkthrough(mapData, 0)
	if err != nil {
		return nil, fmt.Errorf("walkthrough error: %+v", err)
	}

	mapResult := result.(map[string]interface{})
	bytes, err := json.Marshal(mapResult)
	if err != nil {
		return nil, fmt.Errorf("fail while marshalling result map: %+v", err)
	}

	return bytes, nil
}

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

			err = mask(t, k, v, triggerOpt)
			if err != nil {
				return nil, fmt.Errorf("fail while masking key '%s', depth '%d': %+v", k, currentDepth, err)
			}
		}
	}
	return layer, nil
}

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

func maskLabel(label MaskerLabel, key string, value interface{}) (MaskResult, error) {
	_mapLabelsMu.RLock()
	defer _mapLabelsMu.RUnlock()

	if masker, ok := _mapLabels[label]; ok {
		return masker(key, value)
	}

	return defaultMasker(key, value)
}

func mask(dict map[string]interface{}, key string, value interface{}, opts TriggerOpts) error {
	if !opts.ShouldAppear {
		fmt.Println(123)
		delete(dict, key)
		return nil
	}

	_mapLabelsMu.RLock()
	defer _mapLabelsMu.RUnlock()

	var err error

	masker, ok := _mapLabels[opts.MaskMethod]
	if !ok {
		masker = defaultMasker
	}

	ret, err := masker(key, value)
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
