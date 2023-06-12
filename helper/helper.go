package helper

func ConvertInterfaceToString(is interface{}) []string {
	return is.([]string)
}
