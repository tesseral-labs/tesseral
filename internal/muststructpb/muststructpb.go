package muststructpb

import "google.golang.org/protobuf/types/known/structpb"

func MustNewValue(v interface{}) *structpb.Value {
	val, err := structpb.NewValue(v)
	if err != nil {
		panic(err)
	}
	return val
}

func MustNewStruct(m map[string]interface{}) *structpb.Struct {
	val, err := structpb.NewStruct(m)
	if err != nil {
		panic(err)
	}
	return val
}
