package helper

import (
	"fmt"
	"server/structure"
	"strconv"
)

//To determine point distribution
func AllocatePoints(details structure.RecycleLoad) (int, error) {
	var points int
	fmt.Println("allocate points check one")
	fmt.Println(structure.RecycleWeightage)

	if weightage, ok := structure.RecycleWeightage[details.ItemCategory]; ok {
		val, err := strconv.Atoi(details.ItemWeight)
		if err != nil {
			structure.Error.Println("Weight entered is not a number")
			return 0, err
		}
		points = val / weightage
	}

	return points, nil
}

func ChannelInPut(tempSlice []string) <-chan []string {
	c := make(chan []string)
	c <- tempSlice
	// c <- input2

	close(c) //close a channel after for loop has ended

	return c
}

// //Function to check if user is already logged in
// //*edit this function, since no cookies will be required.
//due to local recyclebin
// func AlreadyLoggedIn(req *http.Request) bool {
// 	myCookie, err := req.Cookie("myCookie")
// 	if err != nil {
// 		return false
// 	}

// 	//returns true
// 	username := structure.MapSessions[myCookie.Value]
// 	_, ok := structure.MapUsers[username]
// 	return ok
// }
