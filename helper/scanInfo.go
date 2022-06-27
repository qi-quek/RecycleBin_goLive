package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/structure"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ScanTemplate struct {
	WrongCategory bool
	WrongValue    bool
	// WeightMap     map[string]int
}

var InternalScanTemplate = ScanTemplate{}

func ScanInfo(res http.ResponseWriter, req *http.Request) {

	if !structure.IsLoggedin {
		http.Redirect(res, req, "/login", http.StatusSeeOther) //change this later
		return
	}

	// InternalScanTemplate.WeightMap = structure.RecycleWeightage

	c := make(chan []string)
	params := mux.Vars(req)

	//*go func so template can run
	//*while we wait for input
	if req.Method == "POST" {

		//*----------------------------------------commented away for change into webpage
		inputCat := req.FormValue("material")
		inputWeight := req.FormValue("weight")

		fmt.Println("test check for inputCat-", inputCat)

		catInternal := cases.Title(language.Und, cases.NoLower).String(strings.ToLower(inputCat))

		if _, ok := structure.RecycleWeightage[catInternal]; !ok {

			InternalScanTemplate.WrongCategory = true

			fmt.Println("test check1a for for entry")

			// 	//*this is for self test
			url := "/scan/" + params["id"]

			http.Redirect(res, req, url, http.StatusSeeOther)
			fmt.Println("test check1b for for entry")
			fmt.Println("test check for url value", url)
			return
		}

		InternalScanTemplate.WrongCategory = false

		var tempSlice []string

		//store the scanned line into variable input
		tempSlice = append(tempSlice, catInternal)

		//*---------commented away to test for input in browser

		if _, err := strconv.Atoi(inputWeight); err != nil {

			InternalScanTemplate.WrongValue = true
			// 	//*this is for the real integration
			// url := "/scan/" + structure.LoggedInVal

			//*this is for self test
			url := "/scan/" + params["id"]
			http.Redirect(res, req, url, http.StatusSeeOther)
			return

		}
		InternalScanTemplate.WrongValue = false

		tempSlice = append(tempSlice, inputWeight)
		//*up till here is manual user input--------------------------

		//*----------------------------------------------------------

		//Go rountine for each incoming input
		//sends in a slice of strings
		go func() {
			//*---------
			c <- tempSlice

			close(c) //close a channel after for loop has ended
			//*-------
		}()

		var output []string

		//*--
		for wordSlice := range c { //iterating through a channel gets the value from the channel and clears it
			//cannot assign to individual variables too
			output = wordSlice
		}
		fmt.Println("test word 1", output[0])
		fmt.Println("test word 2", output[1])

		values := map[string]string{"ItemCategory": output[0], "ItemWeight": output[1]}
		json_data, err := json.Marshal(values)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("check for values value-", values)

		//------------------------------------------------------------------------
		//*uncomment from here for testing
		postURL := "https://localhost:8080/user/" + params["id"]
		resp, err := http.Post(postURL, "application/json",
			bytes.NewBuffer(json_data))

		if err != nil {
			structure.Error.Println("error posting to own API")
		}

		var interMap map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&res)

		fmt.Println("intermap-json values", interMap["json"])

		//*http.Post(postURL, "application/json", bytes.NewBuffer(json_data))

		redirectURL := "https://localhost:8080/user/" + params["id"]

		http.Redirect(res, req, redirectURL, http.StatusSeeOther)

	}

	//------------------------------------------------------------------------
	structure.Tpl.ExecuteTemplate(res, "scanInfo.gohtml", InternalScanTemplate)

}
