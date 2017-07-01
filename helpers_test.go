package vkapi

import (
	"net/url"
	"testing"
)

func TestConcatValues(t *testing.T) {
	val1 := url.Values{}
	val2 := url.Values{}
	val3 := url.Values{}

	val1.Set("key", "val1")
	val1.Set("val1", "key")

	val2.Set("key", "val2")
	val1.Set("val2", "key")

	val3.Set("key", "val3")
	val1.Set("val3", "key")

	res1 := ConcatValues(false, val1, val2, val3)
	if res1.Encode() != "key=val3&val1=key&val2=key&val3=key" {
		t.Fatal(res1.Encode())
	}

	res2 := ConcatValues(true, val1, val2, val3)
	if res2.Encode() != "key=val1%2Cval2%2Cval3&val1=key&val2=key&val3=key" {
		t.Fatal(res2.Encode())
	}
}

func TestConcatInt64ToString(t *testing.T) {
	mass := []int64{5, 15, -10, 70}
	str := ConcatInt64ToString(mass...)
	if str != "5,15,-10,70" {
		t.Fatal(str)
	}
}
