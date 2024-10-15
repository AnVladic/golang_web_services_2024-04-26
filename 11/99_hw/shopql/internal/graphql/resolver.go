package graphql

import "hw11_shopql/internal"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Catalogs map[int]*internal.Catalog
	Sellers  map[int]*internal.Seller
}
