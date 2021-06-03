package main

import (
	"encoding/json"
	"fmt"
	"os"
	"io"
	"strconv"
	"github.com/sunshineplan/imgconv"
	"github.com/gin-gonic/gin"
)

type Product struct {
	ID          int64             `json:"id"`
	Price       int               `json:"price"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Attributes  map[string]string `json:"attributes"`
	Image       string            `json:"image"`
	Thumb       string            `json:"thumb"`
	URL         string            `json:"url"`
}

func (product *Product) addAtributes(attr_names []string, attr_values []string) {
	product.Attributes = make(map[string]string)

	for i := range attr_names {
		product.Attributes[attr_names[i]] = attr_values[i]
	}
}

func (product *Product) saveToDB() error {
	stmt, err := db.Prepare("INSERT INTO products(name, description, price, attributes) values (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(product.Name, product.Description, product.Price, product.getAttributesInJSON())
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	product.ID = id

	return nil
}

func (product *Product) updateToDB() error {
	stmt, err := db.Prepare("UPDATE products SET name=?, description=?, price=?, attributes=? WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(product.Name, product.Description, product.Price, product.getAttributesInJSON(), product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (product *Product) getAttributesInJSON() string {
	b, err := json.Marshal(product.Attributes)
	if err != nil {
		panic(err.Error())
	}


	return string(b)
}

func (product *Product) saveImage(c *gin.Context) error {

	var tempImagePath string
	file, err := c.FormFile("image")
	if err != nil {
		return err
	} else {
		tempImagePath = globalConfig.TempUploadPath + file.Filename
		c.SaveUploadedFile(file, tempImagePath)
	}

	src, err := imgconv.Open(tempImagePath)
	if err != nil {
		return err
	}

	thumb := imgconv.Resize(src, imgconv.ResizeOption{Width: globalConfig.ThumbWidth})

	stringID := strconv.Itoa(int(product.ID))

	photoOutputFile, err := os.Create(globalConfig.ImageUploadPath + stringID + ".jpg")
	if err != nil {
		return err
	}

	thumbOutputFile, err := os.Create(globalConfig.ImageUploadPath + stringID + globalConfig.ThumbPostfix + ".jpg")
	if err != nil {
		return err
	}

	err = imgconv.Write(photoOutputFile, src, imgconv.FormatOption{Format: imgconv.JPEG})
	if err != nil {
		return err
	}

	err = imgconv.Write(thumbOutputFile, thumb, imgconv.FormatOption{Format: imgconv.JPEG})
	if err != nil {
		return err
	}

	defer photoOutputFile.Close() 
	defer thumbOutputFile.Close() 
	defer os.Remove(tempImagePath)
	
	return nil
}

func (product *Product) fillFromForm(c *gin.Context) error {
	intPrice, err := strconv.Atoi(c.PostForm("price"))
	if err != nil {
		intPrice = 0
	}

	product.Name = c.PostForm("name")
	product.Price = intPrice
	product.Description = c.PostForm("description")

	attrNames := c.PostFormArray("attr_name[]")
	attrValues := c.PostFormArray("attr_value[]")

	if len(attrNames) > 0 && len(attrNames) == len(attrValues) {
		product.addAtributes(attrNames, attrValues)
	}

	return nil
}

func (product *Product) createAttrFromJSON(attrJSON string) {
	json.Unmarshal([]byte(attrJSON), &product.Attributes)
}

func (product *Product) getImageLink() string {
	return "/img/pictures/" + strconv.Itoa(int(product.ID)) + ".jpg"
}

func (product *Product) getThumbLink() string {
	return "/img/pictures/" + strconv.Itoa(int(product.ID)) + globalConfig.ThumbPostfix + ".jpg"
}

func (product *Product) getURL() string {
	return "/product/" + strconv.Itoa(int(product.ID));
}

func (product *Product) removeFromDB() error {
	stmt, err := db.Prepare("DELETE FROM products WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (product *Product) removeImages() error {
	stringID := strconv.Itoa(int(product.ID))

	err := os.Remove(globalConfig.ImageUploadPath + stringID + ".jpg")
	if err != nil {
		return err
	}

	err = os.Remove(globalConfig.ImageUploadPath + stringID + globalConfig.ThumbPostfix + ".jpg")
	if err != nil {
		return err
	}

	return nil
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func getAllProducts() []Product {
	rows, err := db.Query("SELECT id, name, description, price, attributes FROM products ORDER BY sort")

	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var productsSlice []Product

	for rows.Next() {
		var product Product
		var attrJSON string

		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &attrJSON)

		if err != nil {
			panic(err.Error())
		}

		product.createAttrFromJSON(attrJSON)
		product.Image = product.getImageLink()
		product.Thumb = product.getThumbLink()
		product.URL  = product.getURL()

		productsSlice = append(productsSlice, product)
	}

	return productsSlice
}

func getProductByID(id int) (Product, error) {
	var product Product
	var attrJSON string
	product.ID = int64(id)

	row := db.QueryRow("SELECT name, description, price, attributes FROM products WHERE id = ?", product.ID)
	err := row.Scan(&product.Name, &product.Description, &product.Price, &attrJSON);
	if err != nil {
		return product, err
	}

	product.createAttrFromJSON(attrJSON)
	product.Image = product.getImageLink()
	product.Thumb = product.getThumbLink()
	product.URL  = product.getURL()

	return product, nil
}

func changeSort(sortMap map[int]int) bool {
	sqlQuery := "INSERT INTO products(id, sort) VALUES "

	i := 0
	sortMapLen := len(sortMap)
	for id, sort := range sortMap {
		sqlQuery += fmt.Sprintf("(%v, %v)", id, sort)
		i++

		if i < sortMapLen {
			sqlQuery += ", "
		}
	}

	sqlQuery += " ON DUPLICATE KEY UPDATE id=VALUES(id), sort=VALUES(sort);"

	_, err := db.Exec(sqlQuery)

	if err != nil {
		panic(err.Error())
	}

	return true
}