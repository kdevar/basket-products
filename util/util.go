package util

func ConvertSToInterface(sarr []string) []interface{} {
	s := make([]interface{}, len(sarr))
	for i, v := range sarr {
		s[i] = v
	}
	return s
}

func ConvertIToInterface(sarr []int) []interface{} {
	s := make([]interface{}, len(sarr))
	for i, v := range sarr {
		s[i] = v
	}
	return s
}
