package main

import (
	"database/sql"
	"github.com/foolin/goview/supports/ginview"
	"github.com/foolin/goview"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var globalConfig Config

func main() {
	err := globalConfig.fillFromYML()
	if err != nil {
		panic(err.Error())
	}

	db, err = sql.Open("mysql", globalConfig.DBAccess)
	if err != nil {
		panic(err.Error())
	}

	gvConfig := goview.Config {
		Root:      "templates",
		Extension: ".html",
		DisableCache: true,
	}

	r := gin.Default()
	r.HTMLRender = ginview.New(gvConfig)
	r.Use(sessions.Sessions(globalConfig.CookieSessionName, sessions.NewCookieStore([]byte(globalConfig.CookieSecretKey))))

	r.GET("/", home)
	r.GET("/about", page)
	r.GET("/get_products_json", getProductsJSON)
	r.GET("/product/:id", productPage)
	r.POST("/create_order", createOrder)
	r.GET("/create_order", createOrder)

	r.GET("/admin/login", login)
	r.POST("/admin/login", loginPost)

	private := r.Group("/admin")
	private.Use(adminRightRequred)
	{
		private.GET("/", adminHome)
		private.GET("/pick_add", adminPickAdd)
		private.POST("/pick_add", adminPickAddPost)
		private.POST("/update_sort", adminUpdateSort)
		private.GET("/edit_pick/:id", adminPickEdit)
		private.POST("/edit_pick/:id", adminPickEditPost)
		private.GET("/remove_pick/:id", adminPickRemove)
		private.GET("/edit_texts", adminEditTexts)
		private.POST("/edit_texts/save", adminEditTextsSave)
	}

	r.Run()
}
