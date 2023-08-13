package utils

import (
	"encoding/json"
	"github.com/ArkamFahry/uploadnexus/server/errors"
)

func ParseRequestBody(body []byte, data interface{}) error {
	const Op errors.Op = "utils.ParseRequestBody"

	err := json.Unmarshal(body, data)
	if err != nil {
		return errors.NewError(Op, "invalid request body", err)
	}

	return nil
}
