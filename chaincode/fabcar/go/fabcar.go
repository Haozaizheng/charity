/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
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

// Define the car structure, with 4 properties.  Structure tags are used by encoding/json library
type Deal struct {
	Name   string `json:"name"`
	Org  string `json:"org"`
	Sum string `json:"sum"`
	Date  string `json:"date"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
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
	if function == "queryDeal" {
		return s.queryDeal(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createDeal" {
		return s.createDeal(APIstub, args)
	} else if function == "queryAllDeals" {
		return s.queryAllDeals(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryDeal(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	dealAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(dealAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	deals := []Deal{
		Deal{Name: "T", Org: "Prius", Sum: "10", Date: "2019-1-1"},
		Deal{Name: "F", Org: "Mustang", Sum: "10", Date: "2019-1-1"},
		Deal{Name: "H", Org: "Tucson", Sum: "10", Date: "2019-1-1"},
		Deal{Name: "V", Org: "Passat", Sum: "10", Date: "2017-1-1"},
		Deal{Name: "T", Org: "S", Sum: "10", Date: "2017-2-1"},
		Deal{Name: "P", Org: "205", Sum: "10", Date: "2018-1-1"},
		Deal{Name: "C", Org: "S22L", Sum: "10", Date: "2018-1-1"},
		Deal{Name: "F", Org: "Punto", Sum: "10", Date: "2017-1-1"},
		Deal{Name: "T", Org: "Nano", Sum: "10", Date: "2018-1-1"},
		Deal{Name: "H", Org: "Barina", Sum: "10", Date: "2017-1-1"},
	}

	i := 0
	for i < len(deals) {
		fmt.Println("i is ", i)
		dealAsBytes, _ := json.Marshal(deals[i])
		APIstub.PutState("DEAL"+strconv.Itoa(i), dealAsBytes)
		fmt.Println("Added", deals[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createDeal(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var deal = Deal{Name: args[1], Org: args[2], Sum: args[3], Date: args[4]}

	dealAsBytes, _ := json.Marshal(deal)
	APIstub.PutState(args[0], dealAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllDeals(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "DEAL0"
	endKey := "DEAL999"

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

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
