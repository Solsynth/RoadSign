package transformers

import (
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var ReplacePath = Transformer{
	ModifyRequest: func(options any, ctx *fiber.Ctx) error {
		opts := DeserializeOptions[struct {
			Pattern string `json:"pattern" toml:"pattern"`
			Value   string `json:"value" toml:"value"`
			Repl    string `json:"repl" toml:"repl"` // Use when complex mode(regexp) enabled
			Complex bool   `json:"complex" toml:"complex"`
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
