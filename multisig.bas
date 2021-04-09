//
//Multisig wallet! Written by thedudelebowski for the Dero smart contract competition!
//Version 2.0 - Competition entry 2021
//Check out https://github.com/lebowski1234/dero-multisig-2021 for setup and usage instructions. 


// This function is used to initialize parameters during install time
Function Initialize() Uint64
10 STORE("txCount", 0) 
20 userSetup(SIGNER())
30 RETURN 0
End Function


//Function userSetup: Contains all variables which must be set by the person deploying the contract.
Function userSetup(signer String) Uint64
10 STORE("numberOfOwners", 2) //Must match number of non-blank "owner" fields below. 
20 STORE("authorizationConsensus", 2) //Number of owners that must sign before a transaction is authorized. Can be 1 to No. of owners. 
30 STORE("owner1", signer) 
//Add additional owners below if more than 6 owners required. 
40 STORE("owner2", ADDRESS_RAW("deto1qxsq0qvvr9k9fvkphdvh97k97trfy75ypn3vwvf2fy9wky7vvg90u7qk9esx5"))
50 STORE("owner3", "")
60 STORE("owner4", "") 
70 STORE("owner5", "") 
80 STORE("owner6", "") 
90 STORE("error", "setup complete") 
100 RETURN 0
End Function



//Entrypoint functions:

//Function Deposit: This entrypoint function is called to deposit funds into the wallet. 
Function Deposit() Uint64
10 RETURN 0
End Function


//Function Send: Create a new transaction. To is a Dero address, Amount is the amount to send from the wallet balance. 
Function Send(To String, Amount Uint64) Uint64
10 DIM ownerNo, txCount as Uint64
50 LET ownerNo = sendChecks(To, Amount, SIGNER())
60 IF ownerNo != 0 THEN GOTO 80
70 RETURN 0 //Initial checks failed, exiting
80 storeTx(To, Amount, ownerNo)
90 storeSigners()
100 LET txCount = LOAD("txCount")
110 STORE("tx" + txCount + "_signer" + ownerNo, 1) //Transaction originator signs transaction automatically
120 LET txCount = txCount + 1
130 STORE("txCount", txCount)
140 RETURN 0
End Function


//Function Sign: The main sign function. Signs transaction ID, and sends Dero if required number of signatures has been reached. 
Function Sign(ID Uint64) Uint64
10 DIM ownerNo as Uint64
50 LET ownerNo = signChecks(SIGNER())
60 IF ownerNo != 0 THEN GOTO 80
70 RETURN 0 //Initial checks failed, exiting
80 sign(ID, ownerNo)
90 IF authorized(ID) == 0 THEN GOTO 110 //In this case, 0 = true, as authorized must return 0 to store values within function.
100 RETURN 0 //Not yet authorized, exiting
110 sendDero(ID)//Send Dero 
120 RETURN 0
End Function



//The following functions are called by the main functions, to break up the code into easy to read sections. 


//Function sendChecks: sequence of checks to perform before transaction is accepted.
Function sendChecks(to String, amount Uint64, signer String) Uint64
10 DIM ownerNo as Uint64
20 IF sendValid(to, amount) == 1 THEN GOTO 30 
25 RETURN 0
30 LET ownerNo = verifySigner(signer)
40 IF ownerNo != 0 THEN GOTO 60
50 RETURN 0
60 RETURN ownerNo //All checks passed, return owner No. to calling function. 
End Function


//Function signChecks: sequence of checks to perform before signing request is accepted.
Function signChecks(signer String) Uint64
10 DIM ownerNo as Uint64
20 LET ownerNo = verifySigner(signer)
30 IF ownerNo != 0 THEN GOTO 50
35 STORE("error", "failed at signChecks L35") //for debug only
40 RETURN 0
50 RETURN ownerNo //All checks passed, return owner No. to calling function. 
End Function


//Function sendValid: Checks whether Send transaction parameters are valid.
Function sendValid(s String, i Uint64) Uint64
10 IF IS_ADDRESS_VALID(ADDRESS_RAW(s)) == 1 THEN GOTO 40
20 STORE("error", "failed at sendValid L20 - recipient not a valid Dero address") //for debug only
30 RETURN 0 //Basic format check has failed, exit
40 IF i >0 THEN GOTO 70
50 STORE("error", "failed at sendValid L50 - Amount to send is zero, not a valid transaction") //for debug only
60 RETURN 0 //Basic format check has failed, exit
70 RETURN 1
End Function


//Function verifySigner: Check that signer is an owner. 
Function verifySigner(s String) Uint64
10 DIM inc, numberOfOwners as Uint64
30 LET numberOfOwners = LOAD("numberOfOwners") //ok up to here!
40 LET inc = 1
50 IF s == LOAD("owner" + inc) THEN GOTO 110 
60 IF inc == numberOfOwners THEN GOTO 90 //we have reached numberOfOwners, and not matched the signers address to an owner.
70 LET inc = inc + 1
80 GOTO 50
90 STORE("error", "failed at verifySigner L90 - Signer address not found in list of owners") //for debug only
100 RETURN 0 //Signer ownership check has failed, result is 0. Calling functon must exit on 0. 
110 RETURN inc //Signer is in list of owners, return owner index.
End Function


//Function storeTx: store a new transaction in the DB
Function storeTx(to String, amount Uint64, owner Uint64) Uint64
10 DIM txCount, ownerNo as Uint64
20 LET txCount = LOAD("txCount")
30 STORE("txIndex_"+txCount, txCount) 
40 STORE("recipient_"+txCount, to)
50 STORE("amount_"+txCount, amount)
60 STORE("sent_"+txCount, 0) //Not sent yet
70 STORE("signatures_"+txCount, 1) //Set to 1 immediately as originator signs as part of transaction creation
80 RETURN 0
End Function 


//Function storeSigners: setup (store) signer fields for a new transaction, based on current txCount.
Function storeSigners() Uint64
10 DIM txCount, ownerNo, numberOfOwners as Uint64
20 LET txCount = LOAD("txCount")
30 LET numberOfOwners = LOAD("numberOfOwners")
40 LET ownerNo = 1
50 STORE("tx" + txCount + "_signer" + ownerNo , 0)
60 IF ownerNo == numberOfOwners THEN GOTO 90
70 LET ownerNo = ownerNo + 1
80 GOTO 50
90 RETURN 0
End Function


//Function sendDero: Retrieve transaction from ID No, send Dero, mark transaction as sent. 
Function sendDero(ID Uint64) Uint64
10 DIM isSent, amount as Uint64
11 DIM to as String
20 LET isSent = LOAD("sent_" + ID)
30 IF isSent == 0 THEN GOTO 60
40 STORE("error", "failed at sendDero L40 - Transaction has already been sent") //for debug only
50 RETURN 0
60 LET amount = LOAD("amount_" + ID)
70 LET to = LOAD("recipient_" + ID)
90 SEND_DERO_TO_ADDRESS(ADDRESS_RAW(to), amount) 
100 STORE("sent_" + ID, 1) //mark tx as sent
110 RETURN 0
End Function


//Function authorized: Counts number of signatures for a transaction, and compares with consensus. Returns 1 if consensus reached.
Function authorized(ID Uint64) Uint64
10 DIM authCount, isSigned, ownerNo, numberOfOwners, authorizationConsensus as Uint64
20 LET numberOfOwners = LOAD("numberOfOwners")
30 LET authorizationConsensus = LOAD("authorizationConsensus")
40 LET ownerNo = 1
50 LET isSigned = LOAD("tx" + ID + "_signer" + ownerNo)
60 LET authCount = authCount + isSigned
70 IF ownerNo == numberOfOwners THEN GOTO 100
80 LET ownerNo = ownerNo + 1
90 GOTO 50
100 IF authCount >= authorizationConsensus THEN GOTO 140
120 RETURN 1
140 RETURN 0
End Function


//Function sign: check TX ID exists, then sign TX if owner has not already signed
Function sign(ID Uint64, owner Uint64) Uint64
10 DIM isSigned, signatures as Uint64
20 IF EXISTS("tx" + ID + "_signer" + owner) THEN GOTO 50
30 STORE("error", "failed at sign L30 - Transaction ID not found") //for debug only
40 RETURN 0
50 LET isSigned = LOAD("tx" + ID + "_signer" + owner)
60 IF isSigned == 0 THEN GOTO 90 //Transaction is not yet signed for this owner.
70 STORE("error", "failed at sign L30 - Transaction already signed for owner " + owner) //for debug only
80 RETURN 0
90 STORE("tx" + ID + "_signer" + owner, 1)//Sign transaction for this owner.
91 LET signatures = LOAD("signatures_"+ID) 
92 STORE("signatures_"+ID, signatures + 1) //increment signature count - for user interface
110 RETURN 0
End Function

