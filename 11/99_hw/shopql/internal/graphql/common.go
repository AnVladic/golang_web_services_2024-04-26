package graphql

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"hw11_shopql/internal"
)

func ModelCatalogToCatalog(modelCatalog *internal.Catalog) *Catalog {
	return &Catalog{
		ID:   modelCatalog.Id,
		Name: &modelCatalog.Name,
	}
}

func ModelItemToItem(modelItem *internal.Item) *Item {
	return &Item{
		ID:   modelItem.Id,
		Name: &modelItem.Name,
	}
}

func ModelSellerToSeller(modelSeller *internal.Seller) *Seller {
	return &Seller{
		ID:   modelSeller.Id,
		Name: &modelSeller.Name,
	}
}

func ModelCartItemsToCartItems(modelsCartItem []*internal.CartItem) []*CartItem {
	var cartItems []*CartItem
	for _, cartItem := range modelsCartItem {
		if cartItem.Count <= 0 {
			continue
		}

		cartItems = append(cartItems, &CartItem{
			Item:     ModelItemToItem(cartItem.Item),
			Quantity: cartItem.Count,
		})
	}
	return cartItems
}

func AuthorizedDirective(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	session := ctx.Value("session")
	if session == nil {
		return nil, fmt.Errorf("User not authorized")
	}

	return next(ctx)
}
