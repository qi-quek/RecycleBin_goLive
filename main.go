package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"server/helper"
	"server/structure"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func init() {

	//parse templates to check if it matches the template format
	//ParseGlob is to parse multiple files at once
	//template.Must takes the reponse of template.ParseFiles
	//and does error checking
	structure.Tpl = template.Must(template.ParseGlob("./template/*"))

	_ = godotenv.Load("app.env")

	structure.RecycleWeightage = map[string]int{
		"Metal":   10,
		"Glass":   5,
		"Paper":   2,
		"Plastic": 1,
	} //*4 materials

	file, err := os.OpenFile("errors.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	structure.Error = log.New(io.MultiWriter(file, os.Stderr), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	structure.Info = log.New(io.MultiWriter(file, os.Stderr), "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	SERVER_ADDRESS := os.Getenv("SERVER_ADDRESS")

	muxRouter := mux.NewRouter() //provide new mux router instance

	muxRouter.PathPrefix("/image").Handler(http.StripPrefix("/image", http.FileServer(http.Dir("image"))))

	muxRouter.HandleFunc("/signup", helper.SignUp)
	muxRouter.HandleFunc("/", helper.Start)
	muxRouter.HandleFunc("/login", helper.Login)
	muxRouter.HandleFunc("/scan/{id}", helper.ScanInfo)
	muxRouter.HandleFunc("/user/{id}", helper.DropOff)
	muxRouter.HandleFunc("/logout", helper.LogOut)

	// muxRouter.HandleFunc("/api/{NRIC}", #input function here).Methods("GET", "DELETE", "POST", "PUT")

	http.ListenAndServeTLS(SERVER_ADDRESS, "./ssl/localhost.cert.pem", "./ssl/localhost.key.pem", muxRouter)

}
