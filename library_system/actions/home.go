package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"library-system/repositories/repository"
	"library-system/services"
	"log"
	"net/http"
)

func HomeHandler(c buffalo.Context) error {
	bookService := services.NewBookServices(repository.NewBookRepository(pop.Connections[ENV]))
	books, err := bookService.GetAllBooks()
	if err != nil {
		log.Println("Error fetching books:", err)
		return c.Render(http.StatusInternalServerError, r.JSON(map[string]string{"error": "Failed to fetch books"}))
	}

	c.Set("books", books)
	c.Set("routes", app.Routes())
	c.Set("rootPath", func() string {
		return "/"
	})
	c.Set("t", func(key string) string {
		return key
	})

	return c.Render(http.StatusOK, r.HTML("home/landing.plush.html"))
}

func LibrarianDashboardHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("pages/librarian-dashboard.plush.html"))
}

func UserDashboardHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("pages/user-dashboard.plush.html"))
}
