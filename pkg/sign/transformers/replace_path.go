package transformers

import (
	"github.com/gofiber/fiber/v2"
	"regexp"
	"strings"
)

var ReplacePath = RequestTransformer{
	ModifyRequest: func(options any, ctx *fiber.Ctx) error {
		opts := DeserializeOptions[struct {
			Pattern string `json:"pattern" yaml:"pattern"`
			Value   string `json:"value" yaml:"value"`
			Repl    string `json:"repl" yaml:"repl"` // Use when complex mode(regexp) enabled
			Complex bool   `json:"complex" yaml:"complex"`
		}](options)
		path := string(ctx.Request().URI().Path())
		if !opts.Complex {
			ctx.Path(strings.ReplaceAll(path, opts.Pattern, opts.Value))
		} else if ex := regexp.MustCompile(opts.Pattern); ex != nil {
			ctx.Path(ex.ReplaceAllString(path, opts.Repl))
		}
		return nil
	},
}
