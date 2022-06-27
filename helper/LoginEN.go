package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"server/structure"
)

type dataAuthExternal struct {
	Id    int    `json:"id"`
	Token string `json:"token"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type dataAuthOuput struct {
	Ok   bool                   `json:"ok"`
	Msg  map[string]interface{} `json:"msg"`
	Data dataAuthExternal       `json:"data"`
}

type dataInputStruct struct {
	WrongCredential bool
	ErrorConnection bool
	NotEightNumber  bool
	AuthFail        bool
}

var dataInternal = dataInputStruct{}

func Login(res http.ResponseWriter, req *http.Request) {

	if structure.IsLoggedin == true {
		//*change to structure.LoggedInVal when testing
		url := "scan/" + structure.LoggedInVal

		// //*remove the url below - it is only for self test
		// url := "scan/" + "hi"

		http.Redirect(res, req, url, http.StatusSeeOther)
	}

	//*from here on it is getting input from template/user--------------------------------------------------------
	if req.Method == "POST" {
		userMobile := req.FormValue("mobile")
		userPassword := req.FormValue("password")

		///-------------------first check for empty input---------------------

		// structure.Wg.Add(1)
		// go func() {
		// 	fmt.Println("1 entered")
		// 	if userMobile == "" || userPassword == "" || len(userPassword) < 7 {
		// 		dataInternal.WrongCredential = true
		// 	} else {
		// 		dataInternal.WrongCredential = false
		// 	}
		// 	structure.Wg.Done()
		// 	fmt.Println("1 done")
		// }()

		// ///-------------------end of first check for empty input---------------------

		// ///-------------------2nd check for not a number input---------------------
		// // _, err := strconv.Atoi(userMobile)

		// if _, err := strconv.Atoi(userMobile); err != nil || len(userMobile) != 8 || string(userMobile[0]) != "8" && userMobile[:1] != "9" { //can use string(bytes[0]) or byte[:1] for string comparison

		// 	dataInternal.NotEightNumber = true

		// 	return
		// } else {
		// 	dataInternal.NotEightNumber = false
		// }

		// //----------------end of 2nd check for not a number input------------------

		// if dataInternal.NotEightNumber == true || dataInternal.WrongCredential == true {
		// 	http.Redirect(res, req, "/login", http.StatusSeeOther)
		// 	return
		// }

		//*----------------------code for posting to main webserver----------------------------

		values := map[string]interface{}{"phone": userMobile, "password": userPassword}

		json_data, err := json.Marshal(values)

		if err != nil {
			structure.Info.Println("json marshal fail")
			http.Redirect(res, req, "/login", http.StatusSeeOther)
			return
		}

		postURL := "https://localhost:9091/api/v1/auth"
		resp, err := http.Post(postURL, "application/json",
			bytes.NewBuffer(json_data))

		if err != nil {
			structure.Info.Println("connection failed")
			dataInternal.ErrorConnection = true
			http.Redirect(res, req, "/login", http.StatusSeeOther)
			return
		}

		dataInternal.ErrorConnection = false

		// *--------------------------------2nd part of posting to main webserver----------------------------------

		// read value from response body

		reqBody, err := ioutil.ReadAll((resp.Body))

		//declare variable for data receiving info from post request body
		var loginDataVerify dataAuthOuput

		if err == nil {
			//unmarshal response body into var loginDataVerify
			json.Unmarshal(reqBody, &loginDataVerify)

			fmt.Println(loginDataVerify.Data.Id)

			if loginDataVerify.Ok == true {

				fmt.Println()

				//change global login variable to true
				structure.IsLoggedin = true

				//set var NotNumber to false
				dataInternal.NotEightNumber = false

				//set var WrongCredential to false
				dataInternal.WrongCredential = false

				//set var ErrorConnection to false
				dataInternal.ErrorConnection = false

				structure.LoggedInVal = fmt.Sprint(loginDataVerify.Data.Id)

				structure.JwtToken = loginDataVerify.Data.Token

				//update Id value into url string, for further calling into
				//scan handler
				url := "scan/" + structure.LoggedInVal
				http.Redirect(res, req, url, http.StatusSeeOther)
				return

			} else {
				//if not verified
				//authentication boolean var changed to true
				dataInternal.AuthFail = true
			}

		} else {
			dataInternal.ErrorConnection = true
		}

	}

	//*if authentication from Singpass is true,
	//*structure.IsLogin = true
	structure.Tpl.ExecuteTemplate(res, "loginWork.gohtml", dataInternal)
}
