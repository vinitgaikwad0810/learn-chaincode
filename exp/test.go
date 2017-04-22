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
	//"strconv"
	"strings"
)

type Param struct {
	SensorValue    string `json:"sensorValue"`
	PhValue        string `json:"phValue"`
	AlcoholContent string `json:"alcoholContent"`
	Temperature    string `json:"temperature"`
	Humidity       string `json:"humidity"`
}

type EventInfo struct {
	Qrcode   string  `json:"qrcode"`
	Lat      string  `json:"lat"`
	Lng      string  `json:"lng"`
	Username string  `json:"username"`
	Params   Param `json:"params"`
}

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

	fmt.Println("----" + vState["lat"].(string))

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

// 	eventInfo := `{
//   "qrcode": "3fdsf-324-234-fds5",
//   "lat": "23.8859",
//   "lng": "45.0792",
//   "username": "awaise@gmail.com",
//   "params": {
//     "sensorValue": "23",
//     "phValue": "0.7",
//     "alcoholContent": "67",
//     "temperature": "24",
//     "humidity": "23"
//   }
// }`

	eventInfo := `{"qrcode":"3fdsf-324-234-fds5","lat":"23.8859","lng":"45.0792","username":"awaise@gmail.com","params":{"sensorValue":"23","phValue":"0.7","alcoholContent":"67","temperature":"24","humidity":"23"}}`

//	s, _ := strconv.Unquote(eventInfo)

	arr := []byte(eventInfo)

	var eventInfoStruct EventInfo

	json.Unmarshal(arr, &eventInfoStruct)

	fmt.Printf("%+v\n", eventInfoStruct)

	ret := validateEvent(contractInfo, eventInfo)

	fmt.Println(ret)

}
