package helper

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"server/structure"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type dataStruct struct {
	Weight   string
	Material string
	Points   int
}

type dropData struct {
	Item   string `json:"item"`
	Phone  string `json:"phone"`
	Points int    `json:"points"`
	Weight int    `json:"weight"`
}

type dropPostResp struct {
	Ok   bool     `json:"ok"`
	Msg  string   `json:"msg"`
	Data dropData `json:"data"`
}

var data = dataStruct{}

func DropOff(res http.ResponseWriter, req *http.Request) {

	// if !structure.IsLoggedin {
	// 	http.Redirect(res, req, "/login", http.StatusSeeOther) //change this later
	// 	return
	// }

	params := mux.Vars(req)

	_, err := req.Cookie(params["id"])

	if err != nil {
		http.Redirect(res, req, "/login", http.StatusSeeOther) //change this later
		return
	}

	var receivedReqInfo dropPostResp

	//mux.vars() returns route variable of current request

	fmt.Println("test check  initial data check", data)

	// //*----------------------
	// params := mux.Vars(req)
	// //*----------------------

	//check for json
	//from body content type
	if req.Header.Get("Content-type") == "application/json" {

		if req.Method == "POST" {

			//per instance var
			var newDropOff structure.RecycleLoad

			//read from http request
			//incoming Post command from scanner in this case
			reqBody, err := ioutil.ReadAll((req.Body))
			//successful call returns err = nil
			//it does not treat src to End Of File(EOF) as an error

			if err == nil {
				//*check structure.ReycleLoad fields and decide what is wanted
				json.Unmarshal(reqBody, &newDropOff)

				fmt.Println("test value of .itemcat-", newDropOff.ItemCategory)
				fmt.Println("test value of .itemcat-", newDropOff.ItemWeight)

				if newDropOff.ItemCategory == "" || newDropOff.ItemWeight == "" {
					res.WriteHeader(http.StatusUnprocessableEntity)
					res.Write([]byte("422 - Please supply scanner input"))
					http.Redirect(res, req, "https://localhost:8080/logout", http.StatusSeeOther)
					return
				}

				///To Uppercase first letter of each itemcategory
				//reassign back to Itemcategory
				newDropOff.ItemCategory = cases.Title(language.Und, cases.NoLower).String(strings.ToLower(newDropOff.ItemCategory))
				fmt.Println("test xxxx-", newDropOff.ItemCategory)

				//check if key is in map
				if _, ok := structure.RecycleWeightage[newDropOff.ItemCategory]; !ok {
					res.WriteHeader(http.StatusUnprocessableEntity)
					res.Write([]byte("423 - Please drop off correct recyclable"))
					structure.Error.Println("Incorrect Recyclable material")
					http.Redirect(res, req, "https://localhost:8080/login", http.StatusSeeOther)
					return
				} else {
					fmt.Println("Test check map check success")
				}

			}

			points, err := AllocatePoints(newDropOff)
			if err != nil {
				structure.Info.Panicln("Not a number")
				fmt.Println(err)
				return
			}

			data = dataStruct{Weight: newDropOff.ItemWeight,
				Material: newDropOff.ItemCategory,
				Points:   points}

			fmt.Println("test check, end of self test for drop off")
			//*------------------end of self test for drop off--------------

			//*---------start of post to dylan----------------------------
			itemWeight, err := strconv.Atoi(newDropOff.ItemWeight)

			values := map[string]interface{}{"item_cat": newDropOff.ItemCategory, "points": points, "wgt_in_grams": itemWeight}

			json_data, err := json.Marshal(values)

			if err != nil {
				log.Fatal(err)
				structure.Info.Println("error converting to json_data")
			}

			// //*-----changed to dylan's url-----to add API key-------------------
			//------------------------------------------------------------------------\
			url := "https://localhost:9091/api/v1/users/" + params["id"] + "/transactions"

			//*------------------------
			client := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			}

			// preparing the information to send out
			apiReq, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))

			if err != nil {
				structure.Error.Panicln(`[TRANX-CTL] transaction data parsing, failed`)
				return
			}
			// set all the required header attributes of this POST request.

			authorizationString := "Bearer " + structure.JwtToken //put variable
			apiReq.Header.Set("Content-Type", "application/json")
			apiReq.Header.Set("Authorization", authorizationString)
			// apiReq.Header.Set("username", API_USERNAME)

			resp, err := client.Do(apiReq)

			if err != nil {
				structure.Error.Println(`[TRANX-CTL] transaction data parsing, failed`)
				return
			}

			reqBodyNew, err := ioutil.ReadAll((resp.Body))

			if err != nil {
				structure.Error.Println("reading request body fail")
			}

			json.Unmarshal(reqBodyNew, &receivedReqInfo)

			structure.Info.Println(receivedReqInfo.Msg)

			//*------------------------

			// structure.Info.Panicln(apiReq, "response from post")

			if err != nil {
				log.Fatal(err)
			}
			//*---------------------------------------------

		}
	}
	fmt.Println("test check", data)

	structure.Tpl.ExecuteTemplate(res, "dropOff.gohtml", data)
}
