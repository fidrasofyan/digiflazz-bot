package job

import (
	"context"
	"log"

	"github.com/fidrasofyan/digiflazz-bot/database"
	"github.com/fidrasofyan/digiflazz-bot/internal/service"
)

func PopulateProducts(ctx context.Context) error {
	// Get prepaid products from Digiflazz
	log.Println("PopulateProducts: fetching prepaid products...")
	res, err := service.DigiflazzGetPrepaidPriceList(ctx)
	if err != nil {
		return err
	}
	log.Printf("PopulateProducts: fetched %d prepaid products\n", len(res.Data))

	// Start transaction
	tx, err := database.DBConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // Always defer rollback (will do nothing if already committed)

	qtx := database.Sqlc.WithTx(tx)

	// Purge all prepaid products first
	if err := qtx.DeleteAllPrepaidProducts(ctx); err != nil {
		return err
	}

	// Insert into database
	log.Println("PopulateProducts: inserting prepaid products...")
	for _, product := range res.Data {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err = qtx.InsertPrepaidProduct(ctx, &database.InsertPrepaidProductParams{
			Name:                product.ProductName,
			Category:            product.Category,
			Brand:               product.Brand,
			Type:                product.Type,
			SellerName:          product.SellerName,
			Price:               product.Price,
			BuyerSkuCode:        product.BuyerSKUCode,
			BuyerProductStatus:  product.BuyerProductStatus,
			SellerProductStatus: product.SellerProductStatus,
			UnlimitedStock:      product.UnlimitedStock,
			Stock:               product.Stock,
			Multi:               product.Multi,
			StartCutOff:         &product.StartCutOff,
			EndCutOff:           &product.EndCutOff,
			Description:         &product.Description,
		})
		if err != nil {
			return err
		}
	}

	log.Printf("PopulateProducts: %d prepaid products inserted", len(res.Data))

	return tx.Commit()
}
