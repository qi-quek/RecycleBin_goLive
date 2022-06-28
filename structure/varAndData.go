package structure

import (
	"html/template"
	"log"
	"sync"
)

var Tpl *template.Template

//Map declaration for recyclables material type
//key with type of recyclables material
//value for weightage for point system calculation
var RecycleWeightage map[string]int

// var BoolMap = make(map[string]bool)

//variable to for login check
var IsLoggedin bool

//global value for mobile first 4 digits
// var LoggedInVal string

var TempMap = make(map[string]RecycleLoad)

//Error type log.Logger struct
var Error *log.Logger

var LoggedInVal string

//Info type log.Logger struct
var Info *log.Logger

type RecycleLoad struct {
	//*currently just using itemcategory and itemweight
	//*uncomment username and nric if need to use down the road
	// UserName     string  `json:"UserName"`
	// NRIC         string  `json:"NRIC"`
	ItemCategory string `json:"ItemCategory"`
	ItemWeight   string `json:"ItemWeight"`
}

// var JsonData []byte

var Wg sync.WaitGroup

//jwt token
var JwtToken string

type User struct {
	Username string
	Password []byte
	token    string
}
