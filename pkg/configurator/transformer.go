package configurator

import (
	"github.com/gofiber/fiber/v2"
	"regexp"
	"strings"
)

type RequestTransformer struct {
	ModifyRequest  func(options any, ctx *fiber.Ctx)
	ModifyResponse func(options any, ctx *fiber.Ctx)
}

type RequestTransformerConfig struct {
	Type    string `json:"type"`
	Options any    `json:"options"`
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

var Transformers = map[string]RequestTransformer{
	"replacePath": {
		ModifyRequest: func(options any, ctx *fiber.Ctx) {
			opts := DeserializeOptions[struct {
				Pattern string `json:"pattern"`
				Value   string `json:"value"`
				Repl    string `json:"repl"` // Use when complex mode(regexp) enabled
				Complex bool   `json:"complex"`
			}](options)
			path := string(ctx.Request().URI().Path())
			if !opts.Complex {
				ctx.Path(strings.ReplaceAll(path, opts.Pattern, opts.Value))
			} else if ex := regexp.MustCompile(opts.Pattern); ex != nil {
				ctx.Path(ex.ReplaceAllString(path, opts.Repl))
			}
		},
	},
}
