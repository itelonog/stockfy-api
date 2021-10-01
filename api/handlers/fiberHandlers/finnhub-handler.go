package fiberHandlers

import (
	"stockfyApi/usecases"
	"stockfyApi/usecases/asset"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type FinnhubApi struct {
	ApplicationLogic usecases.Applications
	Api              asset.ExternalApiRepository
}

func (finn *FinnhubApi) GetSymbol(c *fiber.Ctx) error {
	var err error

	symbolLookup, err := finn.ApplicationLogic.AssetApp.
		AssetVerificationExistence(c.Query("symbol"), c.Query("country"),
			finn.Api)
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if err := c.JSON(&fiber.Map{
		"success":      true,
		"symbolLookup": symbolLookup,
		"message":      "Symbol Lookup via Finnhub returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

func (finn *FinnhubApi) GetSymbolPrice(c *fiber.Ctx) error {
	var err error

	symbolPrice, err := finn.ApplicationLogic.AssetApp.AssetVerificationPrice(
		c.Query("symbol"), c.Query("country"), finn.Api)
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if err := c.JSON(&fiber.Map{
		"success":      true,
		"symbolLookup": symbolPrice,
		"message":      "Symbol Lookup via Finnhub returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return err

}
