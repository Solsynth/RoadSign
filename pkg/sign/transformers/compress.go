package transformers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

var CompressResponse = RequestTransformer{
	ModifyResponse: func(options any, ctx *fiber.Ctx) error {
		opts := DeserializeOptions[struct {
			Level int `json:"level"`
		}](options)

		var fctx = func(c *fasthttp.RequestCtx) {}
		var compressor fasthttp.RequestHandler
		switch opts.Level {
		// Best Speed Mode
		case 1:
			compressor = fasthttp.CompressHandlerBrotliLevel(fctx,
				fasthttp.CompressBrotliBestSpeed,
				fasthttp.CompressBestSpeed,
			)
		// Best Compression Mode
		case 2:
			compressor = fasthttp.CompressHandlerBrotliLevel(fctx,
				fasthttp.CompressBrotliBestCompression,
				fasthttp.CompressBestCompression,
			)
		// Default Mode
		default:
			compressor = fasthttp.CompressHandlerBrotliLevel(fctx,
				fasthttp.CompressBrotliDefaultCompression,
				fasthttp.CompressDefaultCompression,
			)
		}

		compressor(ctx.Context())

		return nil
	},
}
