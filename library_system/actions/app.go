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
			SessionName:  sessionName,
		})

		app.Use(SecurityHeaders)
		app.Use(CORS)
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
		bookGroup.POST("/add", bookController.AddBook)
		bookGroup.DELETE("/remove/{id}", bookController.RemoveBook)
		bookGroup.GET("/search", bookController.SearchBook)
		bookGroup.PUT("/update", bookController.UpdateBook)
		bookGroup.GET("/getBookById/{id}", bookController.GetBookByID)

		userGroup := app.Group("/users")
		userGroup.POST("/register", userController.RegisterUser)

		protectedGroup := userGroup.Group("/")
		protectedGroup.Use(Authorize)
		protectedGroup.POST("/checkout", userController.CheckoutBook)
		protectedGroup.POST("/return", userController.ReturnBook)
		protectedGroup.POST("/reserve", userController.ReserveBook)

		app.ServeFiles("/", packr.New("public", "../public"))
		app.GET("/", HomeHandler)

		app.GET("/user-dashboard", UserDashboardHandler)
		app.GET("/librarian-dashboard", LibrarianDashboardHandler)
	})

	return app
}
