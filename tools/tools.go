package tools

func DereferenceOrNil(value *string) interface{} {
	if value == nil {
		return nil
	}
	return *value
}
