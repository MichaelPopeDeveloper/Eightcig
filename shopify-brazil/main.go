package main

import (
	"fmt"
)

func main() {
	// Configuration
	shopName := "your-shop-name"
	token := "your-access-token"
	app := goshopify.App{
		ApiKey:   "your-api-key",
		Password: "your-password",
	}
	apiVersion := "2021-04" // Adjust as per your required API version

	// Create a Shopify client
	client := goshopify.NewClient(app, shopName, token, goshopify.WithVersion(apiVersion))

	// Retrieve an order using the Order ID
	orderID := int64(12345678) // Replace with your actual Order ID
	order, err := client.Order.Get(orderID, nil)
	if err != nil {
		fmt.Printf("Error fetching order: %v", err)
		return
	}

	// Check if the order is being shipped to Brazil
	if order.ShippingAddress.Country == "Brazil" {
		// Apply free shipping by setting the shipping line price to zero
		order.ShippingLines[0].Price = 0

		// Update the order with the new shipping price
		_, err := client.Order.Update(order)
		if err != nil {
			fmt.Printf("Error updating order: %v", err)
			return
		}

		fmt.Println("Free shipping applied to the order shipped to Brazil!")
	} else {
		fmt.Println("Order is not being shipped to Brazil. No changes made.")
	}
}
