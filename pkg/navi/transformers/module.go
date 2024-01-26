package transformers

import (
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

// Definitions

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Transformer struct {
	ModifyRequest  func(options any, ctx *fiber.Ctx) error
	ModifyResponse func(options any, ctx *fiber.Ctx) error
}

type TransformerConfig struct {
	Type    string `json:"type" toml:"type"`
	Options any    `json:"options" toml:"options"`
}

func (v *TransformerConfig) TransformRequest(ctx *fiber.Ctx) error {
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

func (v *TransformerConfig) TransformResponse(ctx *fiber.Ctx) error {
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

var Transformers = map[string]Transformer{
	"replacePath":      ReplacePath,
	"compressResponse": CompressResponse,
}
