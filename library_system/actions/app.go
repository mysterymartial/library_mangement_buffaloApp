package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/pop/v6"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"library-system/controllers"
	"library-system/repositories/repository"
	"library-system/services"
	"log"
	"net/http"
	"sync"
)

var (
	app          *buffalo.App
	appOnce      sync.Once
	sessionStore sessions.Store
	ENV          = envy.Get("GO_ENV", "development")
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found. Using environment variables")
	}
}

func initSessionStore() {
	sessionKey := []byte("12345678901234567890123456789012")
	store := sessions.NewCookieStore(sessionKey)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30,
		HttpOnly: true,
		Secure:   ENV == "production",
		SameSite: http.SameSiteLaxMode,
	}
	sessionStore = store
}

func App() *buffalo.App {
	appOnce.Do(func() {
		initSessionStore()

		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessionStore,
			SessionName:  "_library_system_session",
			Addr:         ":3000",
		})

		// Apply CORS middleware
		app.Use(CORS)

		app.Use(SecurityHeaders)
		app.Use(SetCurrentUser)

		db := pop.Connections[ENV]

		userRepo := repository.NewUserRepository(db)
		bookRepo := repository.NewBookRepository(db)
		loanRepo := repository.NewLoanRepository(db)

		userService := &services.UserServices{
			UserRepo: userRepo,
			BookRepo: bookRepo,
			LoanRepo: loanRepo,
		}
		bookService := &services.BookServices{
			BookRepo: bookRepo,
		}

		userController := controllers.NewUserController(userService, sessionStore)
		bookController := controllers.NewBookController(bookService)

		bookGroup := app.Group("/books")
		bookGroup.GET("/", bookController.GetAllBooks)
		bookGroup.OPTIONS("/", func(c buffalo.Context) error {
			c.Response().WriteHeader(http.StatusOK)
			return nil
		})
		bookGroup.POST("/add", bookController.AddBook)
		bookGroup.OPTIONS("/add", func(c buffalo.Context) error {
			c.Response().WriteHeader(http.StatusOK)
			return nil
		})
		bookGroup.DELETE("/remove/{id}", bookController.RemoveBook)
		bookGroup.OPTIONS("/remove/{id}", func(c buffalo.Context) error {
			c.Response().WriteHeader(http.StatusOK)
			return nil
		})
		bookGroup.GET("/search", bookController.SearchBook)
		bookGroup.OPTIONS("/search", func(c buffalo.Context) error {
			c.Response().WriteHeader(http.StatusOK)
			return nil
		})
		bookGroup.PUT("/update", bookController.UpdateBook)
		bookGroup.OPTIONS("/update", func(c buffalo.Context) error {
			c.Response().WriteHeader(http.StatusOK)
			return nil
		})
		bookGroup.GET("/getBookById/{id}", bookController.GetBookByID)
		bookGroup.OPTIONS("/getBookById/{id}", func(c buffalo.Context) error {
			c.Response().WriteHeader(http.StatusOK)
			return nil
		})

		userGroup := app.Group("/users")
		userGroup.POST("/register", userController.RegisterUser)
		userGroup.OPTIONS("/register", func(c buffalo.Context) error {
			c.Response().WriteHeader(http.StatusOK)
			return nil
		})

		protectedGroup := userGroup.Group("/")
		//protectedGroup.Use(Authorize)
		protectedGroup.POST("/checkout", userController.CheckoutBook)
		protectedGroup.OPTIONS("/checkout", func(c buffalo.Context) error {
			c.Response().WriteHeader(http.StatusOK)
			return nil
		})
		protectedGroup.POST("/return", userController.ReturnBook)
		protectedGroup.OPTIONS("/return", func(c buffalo.Context) error {
			c.Response().WriteHeader(http.StatusOK)
			return nil
		})
		protectedGroup.POST("/reserve", userController.ReserveBook)
		protectedGroup.OPTIONS("/reserve", func(c buffalo.Context) error {
			c.Response().WriteHeader(http.StatusOK)
			return nil
		})

		app.ServeFiles("/", packr.New("public", "../public"))
		app.GET("/", HomeHandler)

		app.GET("/user-dashboard", UserDashboardHandler)
		app.GET("/librarian-dashboard", LibrarianDashboardHandler)
	})

	return app
}
