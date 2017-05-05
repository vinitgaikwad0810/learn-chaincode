/*
Copyright IBM Corp 2016 All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
		 http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"os"
	"strings"
	//"bytes"
	//"strconv"
)

type Test struct {
	Objective      string `json:"objective"`
	ExpectedResult string `json:"expectedResult"`
	ActualResult   string `json:"actualResult"`
	Status         string `json:"status"`
}

type State struct {
	Text    string `json:"text"`
	Lat     string `json:"lat"`
	Lang    string `json:"lang"`
	Address string `json:"address"`
	Tests   []Test `json:"tests"`
}

type ProductSchema struct {
	ProductID   string  `json:"productId"`
	ProductName string  `json:"productName"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	QrCode      string  `json:"qrCode"`
	States      []State `json:"states"`
}

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
	Params   []Param `json:"params"`
}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "register" {
		return t.register(stub, args)
	} else if function == "putcontract" {
		return t.putcontract(stub, args)
	} else if function == "validate" {
		return t.validate(stub, args)
	}

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "getcontract" {
		return t.getcontract(stub, args)
	}
	// } else if function == "validate" {
	// 	return t.validate(stub, args)
	// }

	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) register(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var productId, productInfo string
	var err error
	fmt.Println("running register()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	productId = args[0] //rename for fun
	productInfo = args[1]
	err = stub.PutState(productId, []byte(productInfo)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

// putcontract - Put the received bytearray smatcontract in the json
func (t *SimpleChaincode) putcontract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[1] //rename for funsies
	value = args[2]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) statequery(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//var key, value string
	var err error
	fmt.Println("running statequery()")

	//Checking the number of arguments to be : inorder -> productId , product state json
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	//Fetch the latest state using the product id.

	//Insert

	statekey := args[0] //rename for funsies
	//statevalue := args[1]

	ProductTraceAsbytes, err := stub.GetState(statekey)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get ProductTrace for  " + statekey + "\"}"

		return nil, errors.New(jsonResp)
	}

	var f interface{} //Interface for marshalling the data received from blockchain contract used for comparison

	err_contract := json.Unmarshal(ProductTraceAsbytes, &f)
	if err_contract != nil {
		os.Exit(1)
	}

	var ouputAsBytes []byte
	producttrace_struct := f.(map[string]interface{})

	for k, v := range producttrace_struct {
		if k == "states" {

			//  fmt.Println(k, "is to be compared", v)
			states_values := v.(interface{})

			output, _ := json.Marshal(states_values)
			ouputAsBytes = []byte(output)

		}

	}

	return ouputAsBytes, nil
}

// getcontract - Get the smart Contract from the blockchain as bytearray
func (t *SimpleChaincode) getcontract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"

		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

//Smart Contract Validaton

func validateEvent(contractInfo string, eventInfo string, eventInfoStruct EventInfo) (bool, []Test, State) {

	var vData map[string]interface{}
	var vState map[string]interface{}
	var tests []Test
	var test Test
	var state State

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

	//fmt.Println("QRcode" + vState["qrcode"].(string))

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

			test = Test{
				Objective:      k,
				ExpectedResult: contractParams[k].(string),
				ActualResult:   eventParams[k].(string),
				Status:         "NOT VERIFIED",
			}

		} else {

			test = Test{
				Objective:      k,
				ExpectedResult: contractParams[k].(string),
				ActualResult:   eventParams[k].(string),
				Status:         "VERIFIED",
			}
		}
		tests = append(tests, test)
	}

	state = State{
		Text:    eventInfoStruct.Qrcode,
		Lat:     eventInfoStruct.Lat,
		Lang:    eventInfoStruct.Lng,
		Address: eventInfoStruct.Username,
		Tests:   tests,
	}

	return true, tests, state
}

func (t *SimpleChaincode) validate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//var key, jsonResp string
	//var err error

	type ValidateResponse struct {
		status string `json:"status"`
	}

	var validateResponse ValidateResponse

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	eventInfo := args[0]
	productId := args[1]

	fmt.Println("\n\n eventInfo ------------------------- " + eventInfo)

	fmt.Println("\n\n productType ------------------------- " + productId)

	valAsbytes, _ := stub.GetState(productId)

	var productSchema ProductSchema

	json.Unmarshal(valAsbytes, &productSchema)
	// n := bytes.IndexByte(valAsbytes, 0)
	// contractInfo := string(valAsbytes[:n])

	productInfo := string(valAsbytes[:])

	fmt.Println("\n productInfo is as follows " + productInfo)

	dec := json.NewDecoder(strings.NewReader(productInfo))

	var vProductInfo map[string]interface{}

	if err := dec.Decode(&vProductInfo); err != nil {
		fmt.Println("ERROR: " + err.Error())

	}

	productType := vProductInfo["category"]

	fmt.Println("\n Product Type is " + productType.(string))

	contractInfoAsbytes, _ := stub.GetState(productType.(string))

	// n := bytes.IndexByte(valAsbytes, 0)
	// contractInfo := string(valAsbytes[:n])

	contractInfo := string(contractInfoAsbytes[:])

	fmt.Println("\n\n Contract Retrieved " + contractInfo)

	arr := []byte(eventInfo)

	var eventInfoStruct EventInfo

	json.Unmarshal(arr, &eventInfoStruct)

	fmt.Printf("Event Info Struct is %+v\n", eventInfoStruct)

	_, _, state := validateEvent(contractInfo, eventInfo, eventInfoStruct)

	fmt.Println("JSON parsed into following struct \n")

	fmt.Printf("%+v\n", productSchema)

	productSchema.States = append(productSchema.States, state)

	fmt.Println("JSON modified into following struct \n")

	fmt.Printf("%+v\n", productSchema)

	productSchemaJSON, err := json.Marshal(productSchema)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Error")
	}

	err = stub.PutState(productId, []byte(string(productSchemaJSON[:]))) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}

	//fmt.Println("First element of interface array is " + vProductInfo["states"].([]interface{})[0].(string))

	validateResponse.status = "success"

	validateResponseAsBytes, _ := json.Marshal(validateResponse) //convert to array of bytes

	return validateResponseAsBytes, nil

	// if err != nil {
	// 	jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
	//
	// 	return nil, errors.New(jsonResp)
	// }
	//
	// fmt.Println(valAsbytes)
	//
	// return nil, nil

}

// // read - query function to read key/value pair
// func (t *SimpleChaincode) validate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
// 	var key, jsonResp string
// 	var err, err_state, err_contract error
//
// 	if len(args) != 3 {
// 		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
// 	}
//
// 	//parameters productid,contractkey,bytearray
// 	//data := `{"product_id":"IOT1124s","Contractid":"232241123","stake_holders":["Saurabh_id123","Vinit_Ajay123"],"sensor_value":"24","payment_percent":"20"}`
//
// 	StateJsonAsbytes := []byte(args[0])
// 	contractkey := args[1]
// 	//productid := args[2]
//
// 	ContractvalAsbytes, err := stub.GetState(contractkey)
// 	if err != nil {
// 		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
//
// 		return nil, errors.New(jsonResp)
// 	}
//
// 	var f interface{} //Interface for marshalling the data received from blockchain contract used for comparison.
// 	var g interface{} //Interface for receiving and marshalling the received data
//
// 	err_contract = json.Unmarshal(ContractvalAsbytes, &f)
// 	if err_contract != nil {
// 		os.Exit(1)
// 	}
//
// 	err_state = json.Unmarshal(StateJsonAsbytes, &g)
// 	if err_state != nil {
// 		os.Exit(1)
// 	}
//
// 	contract_json := f.(map[string]interface{})
//
// 	state_json := g.(map[string]interface{})
//
// 	// The Key value iteration can be done better for dynamicity as a seperate function. to loop over the two structs.
//
// 	var sensor_value, sensor_contract string
//
// 	for k, v := range contract_json {
// 		if k == "sensor_value" {
//
// 			fmt.Println(k, "is to be compared", v)
// 			sensor_value = v.(string)
//
// 		}
//
// 	}
//
// 	for k, v := range state_json {
// 		if k == "sensor_value" {
//
// 			fmt.Println(k, "is to be compared", v)
// 			sensor_contract = v.(string)
// 		}
//
// 	}
//
// 	val1, _ := strconv.Atoi(sensor_value)
// 	val2, _ := strconv.Atoi(sensor_contract)
//
// 	var exception string
//
// 	if val1 < val2 {
// 		exception = `{"result":"Exception: value Not acceptable","status":"failed"}`
// 	} else {
// 		exception = `{"result":"Success","status":"success"}`
// 	}
//
// 	exceptionAsBytes := []byte(exception)
//
// 	/*Section to validate the two jsons and put state only if data is validated*/
//
// 	//Smart Contract Rules :
//
// 	// case : blockchain.sensor_value==received.sensor_value
//
// 	// case : blockchain.expiry_max== received.expiry
//
// 	//if true : insert in to blockchain.
//
// 	return exceptionAsBytes, nil
// }
