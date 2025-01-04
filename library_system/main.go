package main

//
//import (
//	"library-system/actions"
//	"log"
//	"net/http"
//	"os"
//)
//
//func main() {
//	app := actions.App()
//
//	port := os.Getenv("PORT")
//	if port == "" {
//		port = "3000"
//	}
//	log.Printf("Listening on port %s", port)
//	if err := http.ListenAndServe(":"+port, app); err != nil {
//		log.Fatal(err)
//	}
//}
