package utils

import "testing"

func TestParseTime(t *testing.T) {

	str1 := "2014-07-16T20:55:46Z"
	got1, _ := ParseTime(str1)

	want1 := int64(1405544146)
	if got1 != want1 {
		t.Errorf("Expected '%d', but got '%d'", want1, got1)
	}

	str2 := "teststring"
	got2, _ := ParseTime(str2)

	want2 := int64(0)
	if got2 != want2 {
		t.Errorf("Expected '%d', but got '%d'", want2, got2)
	}

}

//TODO :: Test functions for other parse functions
