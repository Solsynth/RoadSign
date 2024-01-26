package navi

import (
	"code.smartsheep.studio/goatworks/roadsign/pkg/navi/transformers"

	"github.com/gofiber/fiber/v2"
)

type RoadApp struct {
	Regions []*Region   `json:"regions"`
	Traces  []RoadTrace `json:"traces"`
}

func (v *RoadApp) Forward(ctx *fiber.Ctx, dest *Destination) error {
	switch dest.GetType() {
	case DestinationHypertext:
		return makeUnifiedResponse(ctx, dest)
	case DestinationStaticFile:
		return makeFileResponse(ctx, dest)
	default:
		return fiber.ErrBadGateway
	}
}

type RequestTransformerConfig = transformers.TransformerConfig
