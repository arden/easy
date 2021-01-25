package stringx

import "testing"

func TestDecodeToInt64(t *testing.T) {

}

func TestEncodeFromInt64(t *testing.T) {
	value := "一"
	str, _ := DecodeToInt64(value)
	println(str)

	value2 := EncodeFromInt64(7)
	println(value2)
}
