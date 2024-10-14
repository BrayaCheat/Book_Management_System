package controllers

import (
	"fmt"
	"server/src/configs"
	"server/src/models"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func CreateBook(c *fiber.Ctx) error {
	db := configs.DB
	book := new(models.Book)
	author := new(models.Author)

	//req from form data
	title := strings.TrimSpace(c.FormValue("title"))
	pubslishedDate := strings.TrimSpace(c.FormValue("published_date"))
	bookCover := strings.TrimSpace(c.FormValue("book_cover"))
	authorId, err := strconv.Atoi(c.FormValue("author_id"))

	if err != nil {
		return c.JSON("Invalid author's id")
	}

	//validation
	if title == "" || pubslishedDate == "" || bookCover == "" {
		return c.JSON("Fields can not be blank!")
	}

	//checking author exist || not
	if err := db.Where("id = ?", authorId).First(&author).Error; err != nil {
		return c.JSON("Author not found")
	}

	//checking duplicate
	if err := db.Where("title = ?", title).First(&book).Error; err == nil {
		return c.Status(409).JSON(
			fiber.Map{
				"message": "Book is already exist!",
			},
		)
	}

	//map to model
	book.Title = title
	book.PubslishedDate = pubslishedDate
	book.BookCover = bookCover
	book.AuthorID = uint(authorId)

	//creating
	if err := db.Create(&book).Error; err != nil {
		return c.JSON("Fail to save book")
	}

	return c.JSON(book)
}

func ListBooks(c *fiber.Ctx) error {

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

	sortBy := c.Query("sortBy", "title")
	orderBy := c.Query("orderBy", "asc")
	if orderBy != "asc" && orderBy != "desc" {
		orderBy = "asc"
	}

	offSet := (page - 1) * pageSize

	books := new([]models.Book)

	// Query the database
	if err := db.
		Offset(offSet).
		Limit(pageSize).
		Order(fmt.Sprintf("%s %s", sortBy, orderBy)).
		Find(&books).Error;
		err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to load books",
		})
	}
 
	// Return the response
	return c.Status(200).JSON(fiber.Map{
		"page":     page,
		"pageSize": pageSize,
		"data":     books,
		"sortBy":   sortBy,
		"orderBy":  orderBy,
		"length":   len(*books),
	})
}

func UpdateBook(c *fiber.Ctx) error {
	db := configs.DB
	id, err := strconv.Atoi(c.Params("id"))
	book := new(models.Book)

	if err != nil {
		return c.JSON("Invalid id")
	}

	if err := db.First(&book, id).Error; err != nil {
		return c.JSON("Could not find book")
	}

	title := strings.TrimSpace(c.FormValue("title"))
	published_date := strings.TrimSpace(c.FormValue("published_date"))
	book_cover := strings.TrimSpace(c.FormValue("book_cover"))

	if title == "" || published_date == "" || book_cover == "" {
		return c.JSON("Fields cant be blank")
	}

	book.Title = title
	book.PubslishedDate = published_date
	book.BookCover = book_cover

	if err := db.Save(&book).Error; err != nil {
		return c.JSON("Could not update the book")
	}

	return c.JSON(book)
}

func DeleteBook(c *fiber.Ctx) error {
	db := configs.DB
	book := new(models.Book)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON("Invalid id")
	}

	if err := db.Where("id = ?", id).First(&book).Error; err != nil {
		return c.JSON("Book is not exist in db")
	}

	if err := db.Delete(&book, id).Error; err != nil {
		return c.JSON("Could not delete")
	}

	return c.JSON(
		fiber.Map{
			"message": "Book deleted!",
		},
	)
}
