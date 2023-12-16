package transformers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/samber/lo"
)

var CompressResponse = RequestTransformer{
	ModifyResponse: func(options any, ctx *fiber.Ctx) error {
		opts := DeserializeOptions[struct {
			Level int `json:"level"`
		}](options)
		level := lo.Ternary(opts.Level < 0 || opts.Level > 2, 0, opts.Level)
		ware := compress.New(compress.Config{Level: compress.Level(level)})

		return ware(ctx)
	},
}
