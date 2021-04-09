//Simple user interface for multisig wallet. 
//https://github.com/Lebowski1234/dero-multisig-2021

package main

import (
	"bufio"	
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"net/http"	
	"strconv"	
	"strings"
	"text/tabwriter"
	"github.com/dixonwille/wmenu"
	
)

//Struct for calling an SC
type CallDeroSC struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
	Params  struct {
		ScDeroDeposit int    `json:"sc_dero_deposit,omitempty"`
		Scid          string `json:"scid"`
		ScRPC         []struct {
			Name     string `json:"name"`
			Datatype string `json:"datatype"`
			Value    interface{} `json:"value"`
		} `json:"sc_rpc"`
	} `json:"params"`
}


//Struct for getting SC parameters
type GetDeroSCParams struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
	Params  struct {
		Scid       string   `json:"scid"`
		Code       bool     `json:"code"`
		Keysstring []string `json:"keysstring"`
	} `json:"params"`
}

type FunctionParams struct {
	Name		string
	Datatype	string
	Value		string
}

//Daemon response
type DaemonResult struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  struct {
		Valuesstring []string `json:"valuesstring"`
		Balance      int      `json:"balance"`
		Code         string   `json:"code"`
		Status       string   `json:"status"`
	} `json:"result"`
}


var daemonURL = "http://127.0.0.1:40402/json_rpc"
var walletURL = "http://127.0.0.1:40403/json_rpc" //default
var scid = ""

func main() {
	loadSCID()
	
	mm := mainMenu()
	err := mm.Run()
	if err != nil {
		fmt.Println(err)
		mm.Run()		
	}
}

//Scan file containing SCID and wallet address
//Expects SCID on line 1, walletURL on line 2
func loadSCID() {
	file, err := os.Open("scid.txt") 
    if err != nil { 
        log.Fatal("Failed to open scid.txt") 
    } 
	scanner := bufio.NewScanner(file) 
    scanner.Split(bufio.ScanLines)
	
	count := 0
	for scanner.Scan() { 
       	line := scanner.Text() 
		if count == 0 {
			scid = line
		}
		if count == 1 {
			if line == "" {
				return
			}
			walletURL = line
			
		}
		count +=1	
		if count == 2 {
			return
		}
	}
	
}


//Get txIndex and balance
func getTxIndex() (string, int, error) {
	var c GetDeroSCParams
	c.Jsonrpc = "2.0"
	c.ID = "0"
	c.Method = "getsc"
	var params struct{Scid string `json:"scid"`; Code bool `json:"code"`; Keysstring []string `json:"keysstring"`}
	params.Scid = scid
	params.Code = false
	
	keyString := make([]string,1)
	keyString[0] = "txCount"
	params.Keysstring = keyString
	c.Params = params
	
	j, err := json.Marshal(c)
	if err != nil {
		return "", 0, err
	}
	
	body := bytes.NewReader(j)
	
	result, err := rpcPost(body, daemonURL)
	if err != nil {
		return "", 0, err
	}
		
	var r DaemonResult
	err = json.Unmarshal([]byte(result), &r)
	if err != nil {
		return "", 0, err
	}
	
	//catch bad response
	if len(r.Result.Valuesstring) == 0 {
		return "", 0, errors.New("Bad response")
	}
	
	//catch key not found
	if strings.Contains(r.Result.Valuesstring[0], "NOT AVAILABLE") {
		return "", 0, errors.New("Key not found")
	}
	
	//note: assume if no errors, balance is available. To do: Full error checking.
	
	return r.Result.Valuesstring[0], r.Result.Balance, nil
		
}

//Display transactions
func displayTransactions() {
	transactions, err := getTransactions()
	if err != nil {
		log.Println(err)
		return
	}
	
	txCount, err := strconv.Atoi(transactions["txCount"])
	if err != nil {
		log.Println(err)
		return
	}
			
	fmt.Println("")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "Tx ID\tAmount\tRecipient\tNumber of Signatures\tSent\t")
	if txCount >0 {
		for i := 0; i <txCount; i++ {
			n := strconv.Itoa(i)
			sent := "No"
			if transactions["sent_" + n] == "1" {
				sent = "Yes"
			}
			s :=  n + "\t" + transactions["amount_" + n] + "\t" + transactions["recipient_" + n] + "\t" + transactions["signatures_" + n] + "\t" + sent + "\t"
			fmt.Fprintln(w, s)
			
		}
	
	}
	
	w.Flush()
	fmt.Println("")
	fmt.Printf("Multisig Wallet Balance: %s\n",transactions["balance"]) //To do in future: convert to float
	fmt.Printf("Number of Transactions: %s\n",transactions["txCount"]) 
	fmt.Printf("Number of Owners: %s\n",transactions["numberOfOwners"]) 
	fmt.Printf("Number of Signatures Required to Authorize a Transaction: %s\n",transactions["authorizationConsensus"])
	fmt.Println("")
	
	fmt.Println("Press 'Enter' to return to menu...")
  	bufio.NewReader(os.Stdin).ReadBytes('\n')
	
	
}

//Get SC params via RPC
func getTransactions() (map[string]string, error){
	transactions := make(map[string]string)
	
	//get txCount and balance
	txIndex, balance, err := getTxIndex()
	if err != nil {
		return transactions, err
	}
	
	txCount, err := strconv.Atoi(txIndex)
	if err != nil {
		return transactions, err
	}
		
	//make list of keys
	keys := make([]string,2)
	keys[0] = "numberOfOwners"
	keys[1] = "authorizationConsensus"
	
	for i := 0; i < txCount; i++ {
		n := strconv.Itoa(i)
		keys = append(keys, "recipient_" + n)
		keys = append(keys, "amount_" + n)
		keys = append(keys, "sent_" + n)
		keys = append(keys, "signatures_" + n)
	}
	
	//get keys from daemon and populate map
	var c GetDeroSCParams
	c.Jsonrpc = "2.0"
	c.ID = "0"
	c.Method = "getsc"
	var params struct{Scid string `json:"scid"`; Code bool `json:"code"`; Keysstring []string `json:"keysstring"`}
	params.Scid = scid
	params.Code = false
	
	for i := 0; i < len(keys); i++ {
		keyString := make([]string,1)
		keyString[0] = keys[i]
		params.Keysstring = keyString
		c.Params = params
	
		j, err := json.Marshal(c)
		if err != nil {
			return transactions, err
		}
	
		//log.Printf("Sending: %s",string(j))
		body := bytes.NewReader(j)
	
		result, err := rpcPost(body, daemonURL)
		if err != nil {
			return transactions, err
		}
		
		var r DaemonResult
		err = json.Unmarshal([]byte(result), &r)
		if err != nil {
			return transactions, err
		}
	
		//catch bad response
		if len(r.Result.Valuesstring) == 0 {
			return transactions, errors.New("Bad response")
		}
	
		//catch key not found
		if strings.Contains(r.Result.Valuesstring[0], "NOT AVAILABLE") {
			return transactions, errors.New("Key not found")
		
		}
		
		//populate map
		transactions[keys[i]] = r.Result.Valuesstring[0]
				
	}
	
	//add balance and txCount
	transactions["balance"] = strconv.Itoa(balance) 
	transactions["txCount"] = strconv.Itoa(txCount)
	
	return transactions, nil
}

//Builds SC function call, sends to Dero wallet. 
//amount = dero amount to send: 100000 = 1 Dero. 
func buildFunctionCall(scid string, amount int, funcName string, paramsList []FunctionParams, paramsExist bool) {
	var c CallDeroSC
	c.Jsonrpc = "2.0"
	c.ID = "0"
	c.Method = "scinvoke"
	if amount > 0 {
		c.Params.ScDeroDeposit = amount
	}
	c.Params.Scid = scid

	paramsLength := 1
	if paramsExist {
		paramsLength = len(paramsList)+1
	}
	
	scrpc := make([]struct{Name string `json:"name"`; Datatype string `json:"datatype"`; Value interface{} `json:"value"`},paramsLength)
	
	//entrypoint is always function name
	scrpc[0].Name = "entrypoint"
	scrpc[0].Datatype = "S"
	scrpc[0].Value = funcName
	
	//add function parameters and values, if we have any
	if paramsLength > 1{
		for i := 0; i < len(paramsList); i++ {	
			scrpc[i+1].Name = paramsList[i].Name
			//if we have a Uint64:
			if paramsList[i].Datatype == "Uint64" {
				scrpc[i+1].Datatype = "U"
				x, err := strconv.Atoi(paramsList[i].Value)
				if err != nil {
					log.Fatal(err)
				}
				scrpc[i+1].Value = x
			} else { //must be a string
				scrpc[i+1].Datatype = "S"
				scrpc[i+1].Value = paramsList[i].Value
			}
		
		}
	}
	c.Params.ScRPC = scrpc
	
	j, err := json.Marshal(c)
	if err != nil {
		log.Println(err)
		return
	}
	
	log.Printf("Sending: %s",string(j))
	body := bytes.NewReader(j)
	
	result, err := rpcPost(body, walletURL)
	if err != nil {
		log.Println(err)
		fmt.Println("")
		return
	}
	log.Println(result)
	fmt.Println("")
}

//rpcPost: Send RPC request, return response body as string 
func rpcPost(body *bytes.Reader, url string) (string, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	return string(response), err
}

//For menu: user enters Dero amount
func getDeroAmount() (int, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var text string
	fmt.Print("Enter Dero Amount (100000 = 1 x Dero): ")
	scanner.Scan()
	text = scanner.Text()
	wmenu.Clear()	
	a := text
	if a == "" { //if user hits enter, amount = 0 - allows quick transactions where there is no dero value required
		return 0, nil
	}
	x, err := strconv.Atoi(a)
	if err != nil {
		return x, err
	}  
	return x, nil
}

//For menu: user enters recipient address
func getRecipient() string {
	scanner := bufio.NewScanner(os.Stdin)
	var text string
	fmt.Print("Enter Recipient Dero Address: ")
	scanner.Scan()
	text = scanner.Text()
	wmenu.Clear()	
	return text
}

//For menu: user enters dero amount for multisig transaction
func getSendAmount() string {
	scanner := bufio.NewScanner(os.Stdin)
	var text string
	fmt.Print("Enter Amount to Send from Multisig Wallet (100000 = 1 x Dero): ")
	scanner.Scan()
	text = scanner.Text()
	wmenu.Clear()	
	return text
}

//For menu: user enters transaction ID
func getTxID() string {
	scanner := bufio.NewScanner(os.Stdin)
	var text string
	fmt.Print("Enter Transaction ID to Authorize: ")
	scanner.Scan()
	text = scanner.Text()
	wmenu.Clear()	
	return text
}


//Main menu 
func mainMenu() *wmenu.Menu {
	menu := wmenu.NewMenu("Multisig Wallet Menu - Using Dero Wallet address: " + walletURL + " - Choose an option!")
	
	//Deposit Dero in wallet
	menu.Option("Deposit Dero in Multisig Wallet", nil, false, func(opt wmenu.Opt) error { 
		wmenu.Clear()
				
		amount, err := getDeroAmount()		
		if err != nil {
			mm := mainMenu()
			return mm.Run()
		}
		
		params := make([]FunctionParams,0)
		buildFunctionCall(scid, amount, "Deposit", params, false)
		
		mm := mainMenu()
		return mm.Run()
	})
	
	//Send Dero from wallet
	menu.Option("Send Dero from Multisig Wallet (Create Transaction)", nil, false, func(opt wmenu.Opt) error { 
		wmenu.Clear()

		recipient := getRecipient()
		amount := getSendAmount()
		
		params := make([]FunctionParams,2)
						
		var p1 FunctionParams
		p1.Name = "To"
		p1.Datatype = "String"
		p1.Value = recipient
		
		var p2 FunctionParams
		p2.Name = "Amount"
		p2.Datatype = "Uint64"
		p2.Value = amount
		
		params[0] = p1
		params[1] = p2
		
		buildFunctionCall(scid, 0, "Send", params, true)
		
		mm := mainMenu()
		return mm.Run()
	})
	
	//Sign transaction
	menu.Option("Sign (Authorize) a Multisig Wallet Transaction", nil, false, func(opt wmenu.Opt) error { 
		wmenu.Clear()

		id := getTxID()
		
		params := make([]FunctionParams,1)
						
		var p1 FunctionParams
		p1.Name = "ID"
		p1.Datatype = "Uint64"
		p1.Value = id
				
		params[0] = p1
		
		buildFunctionCall(scid, 0, "Sign", params, true)
		
		mm := mainMenu()
		return mm.Run()
	})
	
	//Display transactions
	menu.Option("Display Multisig Wallet Transactions", nil, false, func(opt wmenu.Opt) error { 
		wmenu.Clear()
		displayTransactions()
		mm := mainMenu()
		return mm.Run()
	})
	
	//Exit
	menu.Option("Exit", nil, false, func(opt wmenu.Opt) error { 
		wmenu.Clear()
		return nil //Exit	
		
	})
			
	return menu
}



