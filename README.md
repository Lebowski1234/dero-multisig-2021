# Multi-Signature Wallet - 2021 

This smart contract was originally written for the Dero Stargate Smart Contract competition in 2019. Re-released in 2021 for the dARCH — Decentralized Architecture Competition Series - Event 0: [https://forum.dero.io/t/darch-decentralized-architecture-competition-series/1318](https://forum.dero.io/t/darch-decentralized-architecture-competition-series/1318)

![screenshot](https://github.com/Lebowski1234/dero-multisig-2021/blob/bfc308d6532e831f65ebe7ac445350f1bc5861ec/screenshot.png)

## Disclaimer

This smart contract was written for the Dero Stargate testnet. It has not been extensively tested for security vulnerabilities, or peer reviewed. It may require modifications to function correctly on the main network, once the main network is released. Use at your own risk.

## Description

The multisig wallet allows a group of individuals to collectively control a single account. For example, a multisig wallet could be used to administer funds in an online community, with multiple stake holders who must all agree on how funds are spent.

The multisig wallet stores a balance of Dero, and is owned by multiple owners with equal rights. Any owner can create a new send transaction, with a recipient (Dero address) and an amount to send (Dero). Before the transaction is sent by the wallet, multiple owners must approve the transaction by signing it. The number of signatures required is defined when the wallet is set up, and can be any number between 1 and the number of owners.

When the required number of signatures has been received, the wallet executes the transaction automatically, sending the Dero to the recipient. It is up to the owners to make sure sufficient funds are available in the wallet before signing a transaction. If sending the transaction fails due to insufficient funds, the transaction must be signed again after sufficient funds are available in the wallet.

Owner 1 deploys the wallet, but does not have any special priviledges once the wallet is deployed. Owner 1 must deploy the wallet initially. After the wallet is deployed, the rules of the wallet cannot be changed: number of owners, owner addresses, and number of signatures required to authorize a transaction.

## Wallet Setup

Before deployment, the following parameters must be configured by the person deploying the wallet. This person automatically becomes Owner 1.

* Number of owners
* Owner addresses
* Authorization consensus (how many owners must sign a transaction before it is sent)

All user defined parameters are grouped under the userSetup function.

## Wallet Deployment

Once setup has been complete, the contract is deployed:

```
curl --request POST --data-binary @multisig.bas http://127.0.0.1:40403/install_sc
```

The Dero HE (testnet) daemon and wallet must both be running first. For a more complete guide, check out the official Dero documentation:

[https://github.com/deroproject/derohe/tree/main/guide](https://github.com/deroproject/derohe/tree/main/guide)



## User Interface

The UI is written in Go and must be compiled before use. This is not covered here (plenty of beginner guides available elsewhere). 

```
go get github.com/Lebowski1234/dero-multisig-2021
```

The UI has a single dependancy:

```
go get github.com/dixonwille/wmenu
```

Note: All development was done on Windows 10. I have not had a chance to test on Linux yet, although there is no reason why it shouldn't work. 

Once compiled, the smart contract address must be copied into a text file in the same directory as the binary, named 'scid.txt'. The Dero wallet RPC address can also be put into the second line of the text file, to allow running multiple wallets with different port numbers, for testing. For example:

```
c2c555b3bc90b305aee1653dc7cef75082a6d889465a73a68ba427366166797f
http://127.0.0.1:40403/json_rpc
```

The user interface has the following functions:

- Deposit dero into the wallet
- Create a new transaction to send somebody Dero from the multisig wallet
- Sign a transaction created by another owner
- View a list of all transactions, including status

Note: all Dero values are entered and displayed using the as-yet unnamed Dero sub-denomination: 100000 of these units = 1 Dero. 


## Contact
I can be reached in the Dero project Discord channel (thedudelebowski#1775). 

If you found this code useful, any Dero donations are most welcome :) dERoVYHj6uBU4xjXVbn35ZiszZznGP2yZfnxqRSZZWvSbhjBaay8GC7cz8TTC54yfAChAjXCk6akeDh9Nmg8gEjm2G9Jb3wHg1
