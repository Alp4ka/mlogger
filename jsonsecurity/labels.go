package jsonsecurity

import (
	"fmt"
	"strings"
	"sync"
)

// MaskerLabel defines the way we mask the data in config.
type MaskerLabel string

const MaskerLabelPassword MaskerLabel = "PASSWORD"

var passwordMasker = func(key string, value interface{}) (MaskResult, error) {
	return MaskResult{
		Key:   key,
		Value: strings.Repeat(MaskSymbol, len(fmt.Sprint(value))),
	}, nil
}

const MaskerLabelCVV MaskerLabel = "CVV"

var cvvMasker = func(key string, value interface{}) (MaskResult, error) {
	const mask = MaskSymbol + MaskSymbol + MaskSymbol

	return MaskResult{
		Key:   key,
		Value: mask,
	}, nil
}

const MaskerLabelName MaskerLabel = "NAME"

var nameMasker = func(key string, value interface{}) (MaskResult, error) {
	const leadingLength = 1
	const trailingLength = 1

	if ret, err := _maskBetween(fmt.Sprint(value), leadingLength, trailingLength); err != nil {
		return MaskResult{}, err
	} else {
		return MaskResult{
			Key:   key,
			Value: ret,
		}, nil
	}
}

const MaskerLabelCardNumber MaskerLabel = "CARD_NUMBER"

var cardMasker = func(key string, value interface{}) (MaskResult, error) {
	const leadingLength = 6
	const trailingLength = 4

	if ret, err := _maskBetween(fmt.Sprint(value), leadingLength, trailingLength); err != nil {
		return MaskResult{}, err
	} else {
		return MaskResult{
			Key:   key,
			Value: ret,
		}, nil
	}
}

const MaskerLabelPhoneNumber MaskerLabel = "PHONE_NUMBER"

var phoneNumberMasker = func(key string, value interface{}) (MaskResult, error) {
	const leadingLength = 2
	const trailingLength = 4

	if ret, err := _maskBetween(fmt.Sprint(value), leadingLength, trailingLength); err != nil {
		return MaskResult{}, err
	} else {
		return MaskResult{
			Key:   key,
			Value: ret,
		}, nil
	}
}

const MaskerLabelEmail MaskerLabel = "EMAIL"

var emailMasker = func(key string, value interface{}) (MaskResult, error) {
	var leadingLength, trailingLength int

	strValue := fmt.Sprint(value)
	lenValue := len(strValue)

	atIdx := strings.Index(strValue, "@")
	if atIdx != -1 {
		leadingLength = 1
		trailingLength = lenValue - atIdx
	} else {
		leadingLength = lenValue
	}

	if ret, err := _maskBetween(strValue, leadingLength, trailingLength); err != nil {
		return MaskResult{}, err
	} else {
		return MaskResult{
			Key:   key,
			Value: ret,
		}, nil
	}
}

// UTILS

func _maskBetween(value string, leadingLen int, trailingLen int) (string, error) {
	var runeLeading, runeTrailing string

	lenValue := len(value)
	lenConstraint := leadingLen + trailingLen

	masked := strings.Repeat(MaskSymbol, lenValue)
	if lenValue >= lenConstraint {
		runeLeading, runeTrailing = value[0:leadingLen], value[lenValue-trailingLen:]
		masked = _noWhitespaceRegexCompiled.ReplaceAllString(value[leadingLen:lenValue-trailingLen], MaskSymbol)
	}

	return runeLeading + masked + runeTrailing, nil
}

var (
	_mapLabelsMu sync.RWMutex
	_mapLabels   map[MaskerLabel]MaskerFunc
)

func init() {
	_mapLabelsMu.RLock()
	_mapLabels = make(map[MaskerLabel]MaskerFunc, 6)
	_mapLabels[MaskerLabelCVV] = cvvMasker
	_mapLabels[MaskerLabelPassword] = passwordMasker
	_mapLabels[MaskerLabelEmail] = emailMasker
	_mapLabels[MaskerLabelName] = nameMasker
	_mapLabels[MaskerLabelPhoneNumber] = phoneNumberMasker
	_mapLabels[MaskerLabelCardNumber] = cardMasker

	_mapLabelsMu.RUnlock()
}
