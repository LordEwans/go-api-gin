package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func getBooks(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, books)
}

func bookById(ctx *gin.Context) {
	id := ctx.Param("id")
	book, err := getBookById(id)

	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"reply": "Book not found!"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, book)
}

func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("Book not found!")
}

func checkOutBook(ctx *gin.Context) {
	id, ok := ctx.GetQuery("id")

	if !ok {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"reply": "Missing query param"})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"reply": "Book not found!"})
		return
	}

	if book.Quantity <= 0 {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"reply": "Book not available!"})
		return
	}

	book.Quantity--
	ctx.IndentedJSON(http.StatusOK, book)
}

func returnBook(ctx *gin.Context) {
	id, ok := ctx.GetQuery("id")

	if !ok {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"reply": "Missing query param"})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"reply": "Book not found!"})
		return
	}

	book.Quantity++
	ctx.IndentedJSON(http.StatusOK, book)
}

func createBooks(ctx *gin.Context) {
	var newBook book

	if err := ctx.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	ctx.IndentedJSON(http.StatusCreated, newBook)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBooks)
	router.PATCH("/checkout", checkOutBook)
	router.PATCH("/return", returnBook)
	router.Run(":8080")
}
