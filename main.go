package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Product struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float32 `json:"price"`
}

var products = []Product{
	{Id: "1", Name: "product 1", Quantity: 10, Price: 10},
	{Id: "2", Name: "product 2", Quantity: 20, Price: 15.5},
	{Id: "3", Name: "product 3", Quantity: 3, Price: 2},
}

func getProducts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, products)
}

func getProductById(id string) (*Product, error) {
	for index, product := range products {
		if product.Id == id {
			return &products[index], nil
		}
	}
	return nil, errors.New("Product not found")
}

func getProduct(c *gin.Context) {
	id := c.Param("id")
	product, err := getProductById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, product)
}

func createProduct(c *gin.Context) {
	var newProduct Product
	var err error = c.BindJSON(&newProduct)
	if err != nil {
		return
	}
	products = append(products, newProduct)
	c.IndentedJSON(http.StatusOK, newProduct)
}

func checkoutProduct(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if ok == false {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter"})
		return
	}

	product, err := getProductById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	if product.Quantity-1 < 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Product is out of stock"})
		return
	}

	product.Quantity -= 1
	c.IndentedJSON(http.StatusOK, product)
}

func refillProduct(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if ok == false {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter"})
		return
	}

	product, err := getProductById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	product.Quantity += 1
	c.IndentedJSON(http.StatusOK, product)
}

func main() {
	router := gin.Default()
	router.POST("/products", createProduct)
	router.GET("/products", getProducts)
	router.GET("/products/:id", getProduct)
	router.PATCH("products/checkout", checkoutProduct)
	router.Run("localhost:8000")
}
