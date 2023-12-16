package transformers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

// Definitions

type RequestTransformer struct {
	ModifyRequest  func(options any, ctx *fiber.Ctx)
	ModifyResponse func(options any, ctx *fiber.Ctx)
}

type RequestTransformerConfig struct {
	Type    string `json:"type" yaml:"type"`
	Options any    `json:"options" yaml:"options"`
}

func (v *RequestTransformerConfig) TransformRequest(ctx *fiber.Ctx) {
	for k, f := range Transformers {
		if k == v.Type {
			if f.ModifyRequest != nil {
				f.ModifyRequest(v.Options, ctx)
			}
			break
		}
	}
}

func (v *RequestTransformerConfig) TransformResponse(ctx *fiber.Ctx) {
	for k, f := range Transformers {
		if k == v.Type {
			if f.ModifyResponse != nil {
				f.ModifyResponse(v.Options, ctx)
			}
			break
		}
	}
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
	"replacePath": ReplacePath,
}
