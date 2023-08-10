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
	return &Masker{cfg: cfg}, nil
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

func maskLabel(label MaskerLabel, key string, value interface{}) (MaskResult, error) {
	_mapLabelsMu.RLock()
	defer _mapLabelsMu.RUnlock()

	if masker, ok := _mapLabels[label]; ok {
		return masker(key, value)
	}

	return defaultMasker(key, value)
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
