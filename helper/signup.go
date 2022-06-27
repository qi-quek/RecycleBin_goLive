package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"server/structure"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//sub struct to match response body of post method to webserver
//for verification of new account/user registeration
type dataRegExternal struct {
	Id        int    `json:"token"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Points    int    `json:"points"`
	LastLogin int    `json:"last_login"`
}

//main struct to match response body of post method to webserver
//for verification of new account/user registeration
type dataRegOutput struct {
	Ok   bool                   `json:"ok"`
	Msg  map[string]interface{} `json:"msg"`
	Data dataRegExternal        `json:"data"`
}

//creating type struct to declare as tpye for data
//to be sent into template
type internalRegStruct struct {
	WrongCredential bool
	NotEightNumber  bool
	SignUpSuccess   bool
	NameNumber      bool
	LessThanSix     bool
}

//struct variable to be sent in as data for execute template
//package level variable so the values will not constantly reset
//upon each handler call
var regDataInternal internalRegStruct

func SignUp(res http.ResponseWriter, req *http.Request) {
	//if user is already logged in
	//redirect back to log in page
	//which will subsequently direct to scan page
	//a user who is logged in, does not need to sign up and create an account
	if structure.IsLoggedin == true {
		//*change to structure.LoggedInVal when testing
		// url := "scan/" + structure.LoggedInVal

		//*remove the url below - it is only for self test
		url := "scan/" + "hi"
		http.Redirect(res, req, url, http.StatusSeeOther)
	}

	//initialise signupsuccess as false
	//for every call into signup function
	//signupsuccess is the struct field sent into template for
	//template message to be displayed
	regDataInternal.SignUpSuccess = false

	if req.Method == "POST" {

		//get  inputs by user
		userMobile := req.FormValue("mobile")
		userName := req.FormValue("name")
		userPassword := req.FormValue("password")

		//*-----------------------------------------------------------------------
		// //------------------1st check for empty fields-----------------------------
		// if userMobile == "" || userPassword == "" || userName == "" {
		// 	fmt.Println("test check - empty input")
		// 	regDataInternal.WrongCredential = true
		// 	//redirect back to login page if user has input empty credential
		// 	http.Redirect(res, req, "/signup", http.StatusSeeOther)
		// 	return
		// }

		// regDataInternal.WrongCredential = false
		// ///-------------------end of first check for empty input---------------------

		// ///-------------------2nd check for mobile input---------------------

		// if _, err := strconv.Atoi(userMobile); err != nil || len(userMobile) != 8 || string(userMobile[0]) != "8" && userMobile[:1] != "9" { //can use string(bytes[0]) or byte[:1] for string comparison
		// 	regDataInternal.NotEightNumber = true
		// 	fmt.Println("testcheck NotNumber-", regDataInternal.NotEightNumber)
		// 	//redirect back to login page if
		// 	http.Redirect(res, req, "/signup", http.StatusSeeOther)
		// 	return
		// }
		// regDataInternal.NotEightNumber = false
		// //-------------------end of 2nd check for mobile input---------------------

		// ///---------------------3rd check for name input; letters or whitespace only-------------

		// match, _ := regexp.MatchString("[^a-zA-Z ]", userName)

		// if match {
		// 	regDataInternal.NameNumber = true
		// 	http.Redirect(res, req, "/signup", http.StatusSeeOther)
		// 	return
		// }

		// regDataInternal.NameNumber = false

		// //-----------------------------end of 3rd check------------------------------
		// ///---------------------4th check, length of password longer than 6-------------

		// if len(userPassword) < 6 {
		// 	regDataInternal.LessThanSix = true
		// 	http.Redirect(res, req, "/signup", http.StatusSeeOther)
		// 	return
		// }
		// regDataInternal.LessThanSix = false

		// //-----------------------------end of 4th check------------------------------
		//*--------------------------------------------------------------------------
		//*---------------------------start of go rountine implementation for users----------------------------------
		//------------------1st check for empty fields-----------------------------

		//Go routine for first validation check
		//first validation check for empty strongs of all user inputs
		structure.Wg.Add(1)
		go func() {
			fmt.Println("1 entered")
			if userMobile == "" || userPassword == "" || userName == "" {
				fmt.Println("test check - empty input")
				regDataInternal.WrongCredential = true
			} else {
				regDataInternal.WrongCredential = false
			}
			structure.Wg.Done()
			fmt.Println("1 done")
		}()

		///-------------------end of first check for empty input---------------------

		///-------------------2nd check for mobile input---------------------

		//go routine for second validation check
		//second validation check for mobile input
		//a)all numbers, b)8 numbers only, c)number must start with an 8 or 9
		structure.Wg.Add(1)
		go func() {
			fmt.Println("2 entered")
			if _, err := strconv.Atoi(userMobile); err != nil || len(userMobile) != 8 || string(userMobile[0]) != "8" && userMobile[:1] != "9" { //can use string(bytes[0]) or byte[:1] for string comparison
				regDataInternal.NotEightNumber = true
			} else {
				regDataInternal.NotEightNumber = false
			}
			structure.Wg.Done()
			fmt.Println("2 done")
		}()

		//-------------------end of 2nd check for mobile input---------------------

		///---------------------3rd check for name input; letters or whitespace only-------------

		//go routine for third validation check
		//third validation check for username input
		//username must only contain alphabets and whitespace

		structure.Wg.Add(1)
		go func() {
			fmt.Println("3 entered")
			match, _ := regexp.MatchString("[^a-zA-Z ]", userName)
			if match {
				regDataInternal.NameNumber = true
			} else {
				regDataInternal.NameNumber = false
			}
			structure.Wg.Done()
			fmt.Println("3 done")
		}()

		//-----------------------------end of 3rd check------------------------------
		///---------------------4th check, length of password longer than 6-------------

		// structure.Wg.Add(1)
		// go func() {
		fmt.Println("4 entered")

		//fourth check - non go routine for concurrency
		//fourth check to ensure that password length is at least 6

		if len(userPassword) < 6 {
			regDataInternal.LessThanSix = true
		} else {
			regDataInternal.LessThanSix = false
		}
		// 	structure.Wg.Done()
		// 	// fmt.Println("4 done")
		// }()

		structure.Wg.Wait()

		if regDataInternal.LessThanSix == true || regDataInternal.NameNumber == true || regDataInternal.NotEightNumber == true || regDataInternal.WrongCredential == true {
			http.Redirect(res, req, "/signup", http.StatusSeeOther)
			return
		}
		//-----------------------------end of 4th check------------------------------
		//*------------------------------end of go routine implementation-------------------------------------------

		var newWord string

		//to split a string into an array of string
		//each element of the array of string is each word within the original string
		//separated by white space
		//"this is a fish" -> ["this","is","a","fish"]
		arrUsername := strings.Fields(userName)

		fmt.Println("strings.fields-", arrUsername)

		var newUserArr []string

		//changes all the words within the array into lower case
		//then uppercasing the first letter of each word
		//appending into a new array of strings for further usage
		for _, word := range arrUsername {
			newUserArr = append(newUserArr, cases.Title(language.Und, cases.NoLower).String(strings.ToLower(word)))
		}

		//joins the array of string back into a single string
		//usage to print out
		newWord = strings.Join(newUserArr[:], " ")
		//*------------------------------end of self test----------------------------

		//passing required data for url post method to webserver
		values := map[string]string{"phone": userMobile, "name": newWord, "password": userPassword}

		// fmt.Println("test check for values-map", values)

		//*------------------start of interfacing with dylan's webportal------------
		//marshalling of map to webserver
		json_data, err := json.Marshal(values)

		if err != nil {
			structure.Info.Println("json marshal fail")
			http.Redirect(res, req, "/signup", http.StatusSeeOther)
			return
		}

		//setting post URL
		postURL := "https://localhost:9091/api/v1/register"

		//posting to webserver
		resp, err := http.Post(postURL, "application/json",
			bytes.NewBuffer(json_data))

		//if post to webserver fails
		//redirect back to login page
		if err != nil {
			structure.Error.Println("connection failed")
			dataInternal.ErrorConnection = true
			http.Redirect(res, req, "/login", http.StatusSeeOther)
			return
		}

		reqBody, err := ioutil.ReadAll((resp.Body))

		var authDataResponse dataRegOutput

		json.Unmarshal(reqBody, &authDataResponse)

		//if response body from post request is true
		//user account has been created on database side
		if authDataResponse.Ok == true {
			structure.Info.Println("Account has been created")
			regDataInternal.SignUpSuccess = true
		}

		structure.Info.Println("response message-", authDataResponse.Msg)
		fmt.Println("test check-final")
		//*--------------end of interace with dylan's web portal-----------------------

	}

	structure.Tpl.ExecuteTemplate(res, "signUp.gohtml", regDataInternal)
}
