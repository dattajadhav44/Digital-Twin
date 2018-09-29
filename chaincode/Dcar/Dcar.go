package main

 /* Imports
  * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
  * 2 specific Hyperledger Fabric specific libraries(packages) for Smart Contracts
  */
 import (
	 "bytes"
	 "encoding/json"
	 "fmt"
	 "strconv"
 
	 "github.com/hyperledger/fabric/core/chaincode/shim"
	 sc "github.com/hyperledger/fabric/protos/peer"
 )
 
 // Define the Smart Contract structure
 type SmartContract struct {
 }
 
 // Define the car structure, with 6 properties.  Structure tags are used by encoding/json library
 type Car struct {
	 Vin    string `json:"vin"`
	 Owner  string `json:"owner"`
	 Colour string `json:"colour"`
	 Model  string `json:"model"`
	 Brand  string `json:"brand"`
	 Milleage  string `json:"milleage"`
 }
 
 /*
  * The Init method is called when the Smart Contract "car" is instantiated by the blockchain network
  * Best practice is to have any Ledger initialization in separate function -- see initLedger()
  */
 func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	 return shim.Success(nil)
 }
 
 /*
  * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
  * The calling application program has also specified the particular smart contract function to be called, with arguments
  */
 func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
 
	 // Retrieve the requested Smart Contract function and arguments
	 function, args := APIstub.GetFunctionAndParameters()
	 // Route to the appropriate handler function to interact with the ledger appropriately
	 if function == "queryCar" {
		 return s.queryCar(APIstub, args)
	 } else if function == "initLedger" {
		 return s.initLedger(APIstub)
	 } else if function == "createCar" {
		 return s.createCar(APIstub, args)
	 } else if function == "queryAllCars" {
		 return s.queryAllCars(APIstub)
	 } else if function == "changeCarOwner" {
		 return s.changeCarOwner(APIstub, args)
	 } else if function == "changeCarMilleageAndColour" {
		return s.changeCarOwner(APIstub, args)
	}
 
	 return shim.Error("Invalid Smart Contract function name.")
 }
 
 func (s *SmartContract) queryCar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 1 {
		 return shim.Error("Incorrect number of arguments. Expecting 1")
	 }
 
	 carAsBytes, _ := APIstub.GetState(args[0])
	 return shim.Success(carAsBytes)
 }
 
 func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	 cars := []Car{
		 Car{Vin: "WDAPF4CC2JP603165",Owner: "Dattatray", Colour: "red", Model: "Sprinter", Brand: "Mercedes-Benz", Milleage: "25"},
		 Car{Vin: "WBS8M9C51J5K98915",Owner: "Bhausaheb", Colour: "blue", Model: "M3", Brand: "BMW", Milleage: "20"},
		 Car{Vin: "2C4RC1BG5JR236293",Owner: "George", Colour: "red", Model: "Pacifica", Brand: "Chrysler", Milleage: "25"},
		 Car{Vin: "3VW5DAAT6JM516495",Owner: "Robert", Colour: "blue", Model: "Beetle", Brand: "Volkswagen", Milleage: "30"},
		 Car{Vin: "JF2GTAMC2JH253258",Owner: "Cosmin", Colour: "white", Model: "Crosstrek", Brand: "Subaru", Milleage: "25"},
	 }
 
	 i := 0
	 for i < len(cars) {
		 fmt.Println("i is ", i)
		 carAsBytes, _ := json.Marshal(cars[i])
		 APIstub.PutState("CAR"+strconv.Itoa(i), carAsBytes)
		 fmt.Println("Added", cars[i])
		 i = i + 1
	 }
 
	 return shim.Success(nil)
 }
 
 func (s *SmartContract) createCar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 6 {
		 return shim.Error("Incorrect number of arguments. Expecting 6")
	 }
 
	 var car = Car{Vin: args[1], Owner: args[2], Colour: args[3], Model: args[4], Brand: args[5], Milleage: args[6]} 
	  
	 carAsBytes, _ := json.Marshal(car)
	 APIstub.PutState(args[0], carAsBytes)
 
	 return shim.Success(nil)
 }
 
 func (s *SmartContract) queryAllCars(APIstub shim.ChaincodeStubInterface) sc.Response {
 
	 startKey := "CAR0"
	 endKey := "CAR10"
 
	 resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	 if err != nil {
		 return shim.Error(err.Error())
	 }
	 defer resultsIterator.Close()
 
	 // buffer is a JSON array containing QueryResults
	 var buffer bytes.Buffer
	 buffer.WriteString("[")
 
	 bArrayMemberAlreadyWritten := false
	 for resultsIterator.HasNext() {
		 queryResponse, err := resultsIterator.Next()
		 if err != nil {
			 return shim.Error(err.Error())
		 }
		 // Add a comma before array members, suppress it for the first array member
		 if bArrayMemberAlreadyWritten == true {
			 buffer.WriteString(",")
		 }
		 buffer.WriteString("{\"Key\":")
		 buffer.WriteString("\"")
		 buffer.WriteString(queryResponse.Key)
		 buffer.WriteString("\"")
 
		 buffer.WriteString(", \"Record\":")
		 // Record is a JSON object, so we write as-is
		 buffer.WriteString(string(queryResponse.Value))
		 buffer.WriteString("}")
		 bArrayMemberAlreadyWritten = true
	 }
	 buffer.WriteString("]")
 
	 fmt.Printf("- queryAllCars:\n%s\n", buffer.String())
 
	 return shim.Success(buffer.Bytes())
 }
 
 func (s *SmartContract) changeCarOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 2 {
		 return shim.Error("Incorrect number of arguments. Expecting 2")
	 }
 
	 carAsBytes, _ := APIstub.GetState(args[0])
	 car := Car{}
 
	 json.Unmarshal(carAsBytes, &car)
	 car.Owner = args[1]
 
	 carAsBytes, _ = json.Marshal(car)
	 APIstub.PutState(args[0], carAsBytes)
 
	 return shim.Success(nil)
 }

 func (s *SmartContract) changeCarMilleageAndColor(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	carAsBytes, _ := APIstub.GetState(args[0])
	car := Car{}

	json.Unmarshal(carAsBytes, &car)
	car.Colour = args[1]
	car.Milleage = args[2]
	 

	carAsBytes, _ = json.Marshal(car)
	APIstub.PutState(args[0], carAsBytes)

	return shim.Success(nil)
}
 
 // The main function is only relevant in unit test mode. Only included here for completeness.
 func main() {
 
	 // Create a new Smart Contract
	 err := shim.Start(new(SmartContract))
	 if err != nil {
		 fmt.Printf("Error creating new Smart Contract: %s", err)
	 }
 }