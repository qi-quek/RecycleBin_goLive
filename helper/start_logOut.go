package helper

import (
	"net/http"
	"server/structure"
)

//Function for /start handler
//to display points gained per respective recyclable material
func Start(res http.ResponseWriter, req *http.Request) {
	displayWeightage := make(map[string]float32)

	//to create a map that matches intended data to be shown
	for key, value := range structure.RecycleWeightage {
		displayWeightage[key] = 1 / float32(value)
	}
	structure.Tpl.ExecuteTemplate(res, "start.gohtml", displayWeightage)
}

//function for /logout handler
//to logout and clear all global variable
func LogOut(res http.ResponseWriter, req *http.Request) {

	//redirect back to sign in page if user is not logged in
	if !structure.IsLoggedin {
		http.Redirect(res, req, "/", http.StatusSeeOther) //change this later
		return
	}

	//*Change all global var back to false or empty string
	//clear IsLoggedin status
	structure.IsLoggedin = false

	//clear loggedinval
	structure.LoggedInVal = ""

	//clear data
	data = dataStruct{}

	//clear JwtToken
	structure.JwtToken = ""

	//clear Internscan
	structure.Tpl.ExecuteTemplate(res, "logOut.gohtml", nil)
}
