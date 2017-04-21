package main

import (

	// Standard library packages
	"encoding/json"
	"fmt"
	// "net/http"
	// // Third party packages
	// "github.com/julienschmidt/httprouter"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	// "net/url"
	//"reflect"
	"strings"
)

func validateEvent(contractInfo string, eventInfo string) bool {

	var vData map[string]interface{}
	var vState map[string]interface{}

	dec := json.NewDecoder(strings.NewReader(contractInfo))

	if err := dec.Decode(&vData); err != nil {
		fmt.Println("ERROR: " + err.Error())

	}

	dec = json.NewDecoder(strings.NewReader(eventInfo))

	if err := dec.Decode(&vState); err != nil {
		fmt.Println("ERROR: " + err.Error())

	}

	contractParams := vData["params"].(map[string]interface{})

	eventParams := vState["params"].(map[string]interface{})

	fmt.Println("\n\n eventParams--------------------------")

	for k, v := range eventParams {

		fmt.Println(k + "-" + v.(string))
	}

	fmt.Println("\n\n contractParams--------------------------")

	for k, v := range contractParams {

		fmt.Println(k + "-" + v.(string))

		fmt.Println(eventParams[k])

		if eventParams[k] == nil || eventParams[k] != contractParams[k] {

			fmt.Println("(" + k + "," + eventParams[k].(string) + ") is absent or different from what is expected in smart contract")
			return false
		}
	}

	return true
}

func main() {

	// 	data := `{
	//   "contactId": "2343",
	//   "productType": "DrugB",
	//   "params": [
	//     {
	//       "sensorValue": "23"
	//     },
	//     {
	//       "phValue": "0.7"
	//     },
	//     {
	//       "alcoholContent": "67"
	//     },
	//     {
	//       "temperature": "24"
	//     },
	//     {
	//       "humidity": "23"
	//     }
	//   ]
	// }`

	contractInfo := `{
  "contactId": "2343",
  "productType": "DrugB",
  "params": {
    "sensorValue": "23",
    "phValue": "0.7",
    "alcoholContent": "67",
    "temperature": "24",
    "humidity": "23"
  }
}`

	eventInfo := `{

   "params":{

	 "sensorValue": "23",
 	"phValue": "0.7",
 	"alcoholContent": "67",
 	"temperature": "24",
	"humidity": "24"

 }


 }`

	ret := validateEvent(contractInfo, eventInfo)

	fmt.Println(ret)

}
