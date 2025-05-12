package controller

import "github.com/mitchellh/mapstructure"

// Helper: decode struct to map[string]interface{} using mapstructure
func structToMap(input interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	config := &mapstructure.DecoderConfig{
		TagName:              "mapstructure",
		Result:               &result,
		WeaklyTypedInput:     true,
		IgnoreUntaggedFields: true,
		ZeroFields:           true,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}
	err = decoder.Decode(input)
	return result, err
}
