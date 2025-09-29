package data

import (
	"order-food/internal/models"
)

var catalog = map[string]models.Product{
	"10": {ID: "10", Name: "Chicken Waffle", Price: 7.99, Category: "Waffle"},
	"11": {ID: "11", Name: "Veg Waffle", Price: 5.49, Category: "Waffle"},
	"21": {ID: "21", Name: "Margherita Pizza", Price: 8.75, Category: "Pizza"},
	"31": {ID: "31", Name: "Iced Coffee", Price: 2.10, Category: "Beverage"},
	"41": {ID: "41", Name: "Chocolate Shake", Price: 3.40, Category: "Beverage"},
}

func AllProducts() []models.Product {
	out := []models.Product{}
	for _, p := range catalog {
		out = append(out, p)
	}
	return out
}

func ProductByID(id string) (models.Product, bool) {
	p, ok := catalog[id]
	return p, ok
}
