package transformers

import (
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

// Definitions

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type RequestTransformer struct {
	ModifyRequest  func(options any, ctx *fiber.Ctx) error
	ModifyResponse func(options any, ctx *fiber.Ctx) error
}

type RequestTransformerConfig struct {
	Type    string `json:"type" yaml:"type"`
	Options any    `json:"options" yaml:"options"`
}

func (v *RequestTransformerConfig) TransformRequest(ctx *fiber.Ctx) error {
	for k, f := range Transformers {
		if k == v.Type {
			if f.ModifyRequest != nil {
				return f.ModifyRequest(v.Options, ctx)
			}
			break
		}
	}
	return nil
}

func (v *RequestTransformerConfig) TransformResponse(ctx *fiber.Ctx) error {
	for k, f := range Transformers {
		if k == v.Type {
			if f.ModifyResponse != nil {
				return f.ModifyResponse(v.Options, ctx)
			}
			break
		}
	}
	return nil
}

// Helpers

func DeserializeOptions[T any](data any) T {
	var out T
	raw, _ := json.Marshal(data)
	_ = json.Unmarshal(raw, &out)
	return out
}

// Map of Transformers
// Every transformer need to be mapped here so that they can get work.

var Transformers = map[string]RequestTransformer{
	"replacePath":      ReplacePath,
	"compressResponse": CompressResponse,
}
