package vkapi

import (
	"net/url"
	"strconv"
	"strings"
)

// ConcatValues - Concatenating values
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

// ConcatInt64ToString - Concatenating array int64 to string
func ConcatInt64ToString(numbers []int64) string {
	var str []string
	for _, u := range numbers {
		str = append(str, strconv.FormatInt(u, 10))
	}

	return strings.Join(str, ",")
}
