package serializer

import (
	jsoniter "github.com/json-iterator/go"
)

func JSON() jsoniter.API {
	return jsoniter.ConfigCompatibleWithStandardLibrary
}
