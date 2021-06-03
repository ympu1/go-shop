package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
	"html/template"
)

func login(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{})
}

func loginPost(c *gin.Context) {
	session := sessions.Default(c)

	var user User
	user.Name     = c.PostForm("name")
	user.Password = c.PostForm("password")

	if user.Name == "" || user.Password == "" {
		c.HTML(400, "login.html", gin.H{"message": globalConfig.FillFormMessage})
		return
	}

	if !user.checkUserExist() {
		c.HTML(401, "login.html", gin.H{"message": globalConfig.AccessDeneidMessage})
		return
	}

	session.Set(globalConfig.AdminSessionKey, user.Name)

	if err := session.Save(); err != nil {
		c.HTML(500, "login.html", gin.H{"message": globalConfig.SessionErrorMessage})
		return
	}

	c.Redirect(303, "/admin")
}

func adminHome(c *gin.Context) {
	c.HTML(200, "admin/index.html", gin.H{"products": getAllProducts()})
}

func adminPickAdd(c *gin.Context) {
	c.HTML(200, "admin/pick_add.html", gin.H{})
}

func adminPickAddPost(c *gin.Context) {
	var product Product

	product.fillFromForm(c)

	err := product.saveToDB()
	if err != nil {
		panic(err)
	}

	err = product.saveImage(c)
	if err != nil {
		panic(err)
	}

	c.HTML(200, "pick_add.html", gin.H{"show_succ_message": "true"})
}

func showErrorMessage(code int, error_text string, c *gin.Context) {
	c.String(code, error_text)
	c.Abort()
}

func adminRightRequred(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globalConfig.AdminSessionKey)
	if user == nil {
		c.Redirect(307, "/admin/login")
		c.AbortWithStatus(307)
		return
	}

	c.Next()
}

func adminUpdateSort(c *gin.Context) {
	c.Request.ParseForm()

	sortMap := make(map[int]int)

	for key, value := range c.Request.PostForm {

		id, err   := strconv.Atoi(key)
		if err != nil {
			panic(err)
		}
		sort, err := strconv.Atoi(value[0])
		if err != nil {
			panic(err)
		}

		sortMap[id] = sort
	}

	success := changeSort(sortMap)

	c.JSON(200, gin.H{"success": success})
}

func adminPickEdit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		showErrorMessage(404, "404", c);
	}

	product, err := getProductByID(id)
	if (err != nil) {
		showErrorMessage(404, "404", c);
	} else {
		c.HTML(200, "admin/pick_edit.html", gin.H{"product": &product})
	}
}

func adminPickEditPost(c *gin.Context) {
	var product Product

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		showErrorMessage(404, "404", c);
	}

	product.ID = int64(id)
	product.fillFromForm(c)

	err = product.updateToDB()
	if err != nil {
		panic(err)
	}

	_ = product.saveImage(c)

	c.HTML(200, "admin/pick_edit.html", gin.H{"product": &product, "show_succ_message": true})
}

func adminPickRemove(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		showErrorMessage(404, "404", c);
	}

	var product Product
	product.ID = int64(id)

	err = product.removeFromDB()
	if err != nil {
		panic(err)
	}

	product.removeImages()

	c.Redirect(307, "/admin")
}

func adminEditTexts(c *gin.Context) {
	homePage, err  := getPageByName("home") 
	if err != nil {
		panic(err)
	}

	aboutPage, err := getPageByName("about") 
	if err != nil {
		panic(err)
	}

	data := gin.H{}

	data["homeTitle"]       = homePage.Title
	data["homeDescription"] = homePage.Description
	data["homeText"]        = homePage.Text
	data["homeH1"]          = homePage.H1

	data["aboutTitle"]       = aboutPage.Title
	data["aboutDescription"] = aboutPage.Description
	data["aboutText"]        = aboutPage.Text
	data["aboutH1"]          = aboutPage.H1

	c.HTML(200, "admin/edit_texts.html", data)
}

func adminEditTextsSave(c *gin.Context) {
	var homePage, aboutPage Page

	homePage.Name        = "home"
	homePage.Title       = c.PostForm("home_title")
	homePage.Description = c.PostForm("home_description")
	homePage.H1          = c.PostForm("home_h1")
	homePage.Text        = template.HTML(c.PostForm("home_text"))

	aboutPage.Name        = "about"
	aboutPage.Title       = c.PostForm("about_title")
	aboutPage.Description = c.PostForm("about_description")
	aboutPage.H1          = c.PostForm("about_h1")
	aboutPage.Text        = template.HTML(c.PostForm("about_text"))

	err := homePage.updateToDB()
	if err != nil {
		showErrorMessage(503, "Ошибка сохранения", c)
	}

	err = aboutPage.updateToDB()
	if err != nil {
		showErrorMessage(503, "Ошибка сохранения", c)
	}

	c.Redirect(303, "/admin/edit_texts")
}