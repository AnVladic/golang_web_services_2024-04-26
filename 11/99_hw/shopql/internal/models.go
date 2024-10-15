package internal

import (
	"strconv"
)

type MarketplaceData struct {
	Catalog *Catalog  `json:"catalog"`
	Sellers []*Seller `json:"sellers"`
}

type Catalog struct {
	Id       int        `json:"id"`
	Name     string     `json:"name"`
	Children []*Catalog `json:"childs"`
	Parent   *Catalog   `json:"parent"`
	Items    []*Item    `json:"items"`
}

type Item struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	InStock  int      `json:"in_stock"`
	SellerId int      `json:"seller_id"`
	Parent   *Catalog `json:"-"`
}

type Seller struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Deals int    `json:"deals"`
}

type User struct {
	Email    string      `json:"email"`
	Password []byte      `json:"password"`
	Token    string      `json:"token"`
	Cart     []*CartItem `json:"cart"`
}

type Session struct {
	Id   int   `json:"id"`
	User *User `json:"user"`
}

type CartItem struct {
	Item  *Item `json:"item"`
	Count int   `json:"count"`
}

func FindCatalog(catalogs map[int]*Catalog, id string) *Catalog {
	idInt, _ := strconv.Atoi(id)
	catalog := catalogs[idInt]
	if catalog == nil {
		return nil
	}
	return catalog
}

func FindSeller(sellers map[int]*Seller, id string) *Seller {
	idInt, _ := strconv.Atoi(id)
	seller := sellers[idInt]
	if seller == nil {
		return nil
	}
	return seller
}

func GetItemBySellerId(catalogs map[int]*Catalog, sellerId int, offset int, limit int) []*Item {
	var items []*Item
	for _, catalog := range catalogs {
		for _, item := range catalog.Items {
			if item.SellerId == sellerId {
				offset--
				if offset < 0 {
					items = append(items, item)
				}
				if len(items) == limit {
					break
				}
			}
		}
	}
	return items
}

func (c *Catalog) ToMap() map[int]*Catalog {
	catalogs := map[int]*Catalog{}
	c.SetCatalogToMap(&catalogs)
	return catalogs
}

func (c *Catalog) SetCatalogToMap(catalogs *map[int]*Catalog) {
	(*catalogs)[c.Id] = c
	for _, item := range c.Items {
		item.Parent = c
	}
	for _, child := range c.Children {
		child.SetCatalogToMap(catalogs)
		child.Parent = c
	}
}

func (m *MarketplaceData) SellersToMap() map[int]*Seller {
	sellers := map[int]*Seller{}
	for _, seller := range m.Sellers {
		sellers[seller.Id] = seller
	}
	return sellers
}

func FindItemById(catalogs map[int]*Catalog, id int) *Item {
	var item *Item
	for _, catalog := range catalogs {
		for _, currentItem := range catalog.Items {
			if currentItem.Id == id {
				item = currentItem
			}
		}
	}
	return item
}
