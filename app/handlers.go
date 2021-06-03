package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func home(c *gin.Context) {
	page, err := getPageByName("home")
	if err != nil {
		showErrorMessage(404, "404", c)
	}
	data := gin.H{}
	data["title"]       = page.Title
	data["description"] = page.Description
	data["text"]        = page.Text
	data["h1"]          = page.H1
	data["products"]    = getAllProducts()

	c.HTML(200, "index.html", data)
}

func page(c *gin.Context) {
	page, err := getPageByName("about")
	if err != nil {
		showErrorMessage(404, "404", c)
	}
	data := gin.H{}
	data["title"]       = page.Title
	data["description"] = page.Description
	data["text"]        = page.Text
	data["h1"]          = page.H1

	c.HTML(200, "about.html", data)
}

func getProductsJSON(c *gin.Context) {
	c.JSON(200, gin.H{"products": getAllProducts()})
}

func productPage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		showErrorMessage(404, "404", c)
	}

	product, err := getProductByID(id)
	if (err != nil) {
		showErrorMessage(404, "404", c)
	} else {
		c.HTML(200, "product.html", gin.H{"product": &product})
	}
}

func createOrder(c *gin.Context) {
	var order Order

	order.userName = c.PostForm("user_name")
	order.userInf  = c.PostForm("user_inf")

	if len(order.userName) == 0 || len(order.userInf) == 0 {
		c.JSON(400, gin.H{"success": false})
		return
	}

	err := order.sendNotification()

	if err != nil {
		c.JSON(400, gin.H{"success": false})
		return
	}

	c.JSON(200, gin.H{"success": true})
}