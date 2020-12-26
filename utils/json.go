package utils

import "github.com/tidwall/gjson"
func GetJsonValue(src string, filed string) string {
	value := gjson.Get(src, filed)
	return value.String()
}



func GetJsonArray(src string, filed string) *[]gjson.Result {
	results := gjson.GetMany(src, filed)
	return &results
}

