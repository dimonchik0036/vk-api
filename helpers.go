package vkapi

import "net/url"

//ConcatValues - Concatenating values
func ConcatValues(unite bool, values ...url.Values) (result url.Values) {
	result = url.Values{}
	for _, v := range values {
		for key := range v {
			value := v.Get(key)
			if oldVal := result.Get(key); unite && oldVal != "" {
				value = oldVal + "," + value
			}

			result.Set(key, value)
		}
	}

	return result
}
