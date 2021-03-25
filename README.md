# Multi-Signature Wallet - 2021 

This smart contract was originally written for the Dero Stargate Smart Contract competition in 2019. Re-released in 2021 for the dARCH — Decentralized Architecture Competition Series - Event 0!

The multisig wallet allows a group of individuals to collectively control a single account. For example, a multisig wallet could be used to administer funds in an online community, with multiple stake holders who must all agree on how funds are spent.

The multisig wallet stores a balance of Dero, and is owned by multiple owners with equal rights. Any owner can create a new send transaction, with a recipient (Dero address) and an amount to send (Dero). Before the transaction is sent by the wallet, multiple owners must approve the transaction by signing it. The number of signatures required is defined when the wallet is set up, and can be any number between 1 and the number of owners.

When the required number of signatures has been received, the wallet executes the transaction automatically, sending the Dero to the recipient. It is up to the owners to make sure sufficient funds are available in the wallet before signing a transaction. If sending the transaction fails due to insufficient funds, the transaction must be signed again after sufficient funds are available in the wallet.

Owner 1 deploys the wallet, but does not have any special priviledges once the wallet is deployed. Owner 1 must set up the wallet initially (see below). After the wallet is deployed, the rules of the wallet cannot be changed: number of owners, owner addresses, and number of signatures required to authorize a transaction.

For a full description, refer to the original version readme at [https://github.com/Lebowski1234/dero-multisig](https://github.com/Lebowski1234/dero-multisig). Note that all curl commands in the original readme are no longer valid, and the referenced UI will not work either without major updates. 
