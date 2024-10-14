package controllers

import (
	"errors"
	"fmt"
	"server/src/configs"
	"server/src/models"
	"strconv"
	"strings"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateAuthor(c *fiber.Ctx) error {
	db := configs.DB
	author := new(models.Author)
	authorName := strings.TrimSpace(c.FormValue("name"))
	authorGender := strings.TrimSpace(c.FormValue("gender"))
	if authorName == "" && authorGender == "" {
		return c.JSON("Fields can not blank or whitespace")
	}

	if err := db.Where("name = ?", authorName).First(&author).Error; err == nil {
		return c.JSON("Author is already exist")
	}

	author.Name = authorName
	author.Gender = authorGender
	db.Create(&author)
	return c.JSON(author)
}

func ListAuthor(c *fiber.Ctx) error {
	db := configs.DB

	// Retrieve and validate query parameters
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	sortBy := c.Query("sortBy", "name")
	orderBy := c.Query("orderBy", "asc")
	if orderBy != "asc" && orderBy != "desc" {
		orderBy = "asc"
	}

	offSet := (page - 1) * pageSize

	authors := new([]models.Author)

	// Query the database
	if err := db.
		Preload("Books").
		Preload("AuthorAddress").
		Offset(offSet).
		Limit(pageSize).
		Order(fmt.Sprintf("%s %s", sortBy, orderBy)).
		Find(&authors).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to load authors",
		})
	}
 
	// Return the response
	return c.Status(200).JSON(fiber.Map{
		"page":     page,
		"pageSize": pageSize,
		"data":     authors,
		"sortBy":   sortBy,
		"orderBy":  orderBy,
		"length":   len(*authors),
	})
}

func GetAuthor(c *fiber.Ctx) error {
	db := configs.DB
	id, err := strconv.Atoi(c.Params("id"))
	author := new(models.Author)
	if err != nil {
		return c.JSON("Invalid id")
	}
	if err := db.
		Preload("Books").
		Preload("AuthorAddress").
		Find(&author, id).
		Error; err != nil {
		return c.JSON("Author is not exist!")
	}
	return c.JSON(author)
}

func UpdateAuthor(c *fiber.Ctx) error {
	db := configs.DB
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON("invalid id")
	}

	author := new(models.Author)

	if err := db.First(&author, id).Error; err != nil {
		return c.JSON("Could not find author!")
	}

	name := strings.TrimSpace(c.FormValue("name"))
	gender := strings.TrimSpace(c.FormValue("gender"))

	if name == "" || gender == "" {
		return c.JSON("Field can not be blank!")
	}

	author.Name = name
	author.Gender = gender

	if err := db.Save(&author).Preload("AuthorAddress").Error; err != nil {
		return c.JSON("Could not update the author")
	}

	return c.Status(201).JSON(
		fiber.Map{
			"message": "Author updated!",
			"data":    author,
			"status":  200,
		},
	)
}

func DeleteAuthor(c *fiber.Ctx) error {
	db := configs.DB
	id, _ := strconv.Atoi(c.Params("id"))
	author := new(models.Author)

	if err := db.First(&author, id).Error; err != nil {
		return c.JSON("Could not find author!")
	}

	if err := db.Delete(&author, id).Error; err != nil {
		return c.JSON("Could not delete author!")
	}

	return c.Status(fiber.StatusNoContent).JSON(
		fiber.Map{
			"message": "Author deleted!",
		},
	)
}

func CreateAddress(c *fiber.Ctx) error {
	db := configs.DB
	author := new(models.Author)
	authorAddress := new(models.AuthorAddress)
	street := strings.TrimSpace(c.FormValue("street"))
	town := strings.TrimSpace(c.FormValue("town"))
	city := strings.TrimSpace(c.FormValue("city"))
	country := strings.TrimSpace(c.FormValue("country"))
	authorId, err := strconv.Atoi(c.FormValue("author_id"))

	if err != nil {
		return c.JSON("Invalid author's id")
	}

	if street == "" || town == "" || city == "" || country == "" {
		return c.JSON("Fields cant be blank")
	}

	//author is exist or not
	if err := db.Where("id = ?", authorId).First(&author).Error; err != nil {
		return c.JSON("Author is not exist")
	}

	//author is already has the address
	if err := db.Where("author_id = ?", authorId).First(&authorAddress).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"error": "Author already has an address"})
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check existing address"})
	}

	// Create the address
	authorAddress.Street = street
	authorAddress.Town = town
	authorAddress.City = city
	authorAddress.Country = country
	authorAddress.AuthorID = uint(authorId)

	if err := db.Create(&authorAddress).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save author's address"})
	}

	return c.Status(200).JSON(
		fiber.Map{
			"message": "Author's address created",
			"data":    authorAddress,
		},
	)
}

func ListAddress(c *fiber.Ctx) error {
	db := configs.DB
	addresses := new([]models.AuthorAddress)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))
	sortBy := c.Query("sortBy", "street")
	orderBy := c.Query("orderBy", "asc")

	offSet := (page - 1) * pageSize

	if err := db.Limit(pageSize).Offset(offSet).Order(fmt.Sprintf("%s %s", sortBy, orderBy)).Find(&addresses).Error; err != nil {
		return c.JSON("Failed to load address list")
	}

	return c.Status(200).JSON(
		fiber.Map{
			"page":     page,
			"pageSize": pageSize,
			"data":     addresses,
			"sortBy":   sortBy,
			"orderBy":  orderBy,
			"length":   len(*addresses),
		},
	)
}

func GetAddress(c *fiber.Ctx) error {
	db := configs.DB
	id, err := strconv.Atoi(c.Params("id"))
	authorAddress := new(models.AuthorAddress)

	if err != nil {
		return c.JSON("Invalid id")
	}

	if err := db.First(&authorAddress, id).Error; err != nil {
		return c.JSON("Could not find address")
	}

	return c.JSON(authorAddress)
}

func UpdateAddress(c *fiber.Ctx) error {
	db := configs.DB
	authorAddress := new(models.AuthorAddress)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON("Invalid id")
	}

	if err := db.First(&authorAddress, id).Error; err != nil {
		return c.JSON("Could not find author's address!")
	}

	street := strings.TrimSpace(c.FormValue("street"))
	town := strings.TrimSpace(c.FormValue("town"))
	city := strings.TrimSpace(c.FormValue("city"))
	country := strings.TrimSpace(c.FormValue("country"))

	if street == "" || town == "" || city == "" || country == "" {
		return c.JSON("Fields cant be blank")
	}

	authorAddress.Street = street
	authorAddress.Town = town
	authorAddress.City = city
	authorAddress.Country = country

	if err := db.Save(&authorAddress).Error; err != nil {
		return c.JSON("Failed to update author's address")
	}

	return c.JSON(authorAddress)
}

func DeleteAddress(c *fiber.Ctx) error {
	db := configs.DB
	id, err := strconv.Atoi(c.Params("id"))
	authorAddress := new(models.AuthorAddress)

	if err != nil {
		return c.JSON("Invalid id")
	}

	if err := db.First(&authorAddress, id).Error; err != nil {
		return c.JSON("Could not find author's address!")
	}

	if err := db.Delete(&authorAddress, id).Error; err != nil {
		return c.JSON("Could not delete the address")
	}

	return c.Status(200).JSON(
		fiber.Map{
			"message": "Address deleted!",
		},
	)
}
