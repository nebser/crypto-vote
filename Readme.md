# Crypto vote

Crypto vote project represents a simulation of electronic voting system that is based on a "proof-of-stake" blockchain.

## Prerequisites

1. go 1.14
2. bolt db - https://github.com/boltdb/bolt

## Compilation

I'd strongly suggest using Makefile for performing compilation because there are 6 applications in this project. Just run:

```
~$ make
```

## Basic concepts

In this system there are 3 types of nodes:

1. Alfa node - the main node in the blockchain system. It is in charge for creating the initial state of the blockchain which will enable voting for all users with a valid key-pair. Every new node that enters the system connects to this node to obtain addresses of all nodes in the current system. This node also has an http server that accepts voting requests from clients. Ideally, this node should be under the control of a governing body.
2. Party node - node that can forge new blocks into a blockchain. These nodes are being controlled by the parties who are subjects of the voting process. Forging new blocks is being done by creating a "stake" vote which transfers half of the current votes from the party to the alfa node. If the block is indeed valid that half of the current votes will be returned to the party as part of the next block in blockchain. On the other hand, if the block is in fact invalid, half of the current votes will not be returned to the party and that party will be excommunicated from the system.
3. Client node - node that can retrieve a copy of the blockchain. It cannot forge new blocks, but it can receive updates and validate any block and see if there are any irregularities. This node can be controlled by anyone who has a valid key-pair (in other words, anyone who has a right to vote)


## Applications

In this project there are 6 applications which can help you effectively simulate the voting process

### Key generator

Key generator is a key-pair generator used for generating all of the necessary key-pairs in the system - 1 key pair for alfa node, n key pairs for party nodes and m key-pairs for client nodes. This application accepts 5 options of which all have default values:

1. `alfa` - directory in which to create key pair for the alfa node; default value is `alfa`
2. `clients` - directory in which to create key pairs for clients (voters); default value is `clients`
3. `nodes` - directory in which to create key pairs for party nodes; default value is `nodes`
4. `clientsNumber` - number of key pairs to create for clients (voters); default value is `50`
5. `nodesNumber` - number of key pairs to create for nodes; default value is `5` 

To run key generator with default values type:
```
~$ ./key-generator
```

### Alfa node

Alfa node is the central node in the blockchain system. As soon as it starts it will print the initial blockchain state to the console output. 

Alfa node has a websocket server which communicates with the rest of the nodes in the system. All of the incoming nodes in the system will first register to alfa node and retrieve list of active nodes from it.

This application accepts 5 options which all have default values:

1. `new` - flag that indicates whether or not the node should initialize a new state of the blockchain; default value is `false`
2. `private` - path to private key file which the alfa node will use to sign request, blocks, etc; default value is `alfa/key.pem` (output of the `key` generator)
3. `public` - path to public key file which the alfa node will use as a part of it's address; default value is `alfa/key_pub.pem` (output of the key-generator)
4. `clients` - directory which contains voters public keys. This is necessary for the alfa node to create a transaction output that voters will use to actually create a vote; default value is `clients`
5. `nodes` - directory which contains public keys of nodes in control by parties. This is necessary for the alfa node to track requests from nodes created by parties; default value is `nodes`

To run a new alfa node type:
```
~$ ./alfa-node -new
```

### Client node

Client node is an application that can start a party node or client node based on the key-pair that is passed to it. As soon as it starts it will obtain the blockchain state from the alfa node and all of the running nodes in the system. The difference between party and client node is that the party node can forge new blocks where client node can only verify new blocks.

This application accepts 4 options:

1. `id` - internal id of the client node, must be an integer value greater than 0; there is no default value.
2. `new` - flag that indicates if the block should purge the blockchain it has locally or just take the missing blocks from the alfa node; default value is `false`.
3. `private` - path to private key file that the node will use for signing it's requests and forging new blocks (if it's a party node); default value is `nodes/key_id.pem`
4. `public` - path to public key file which will be used as a part of it's address; default value `nodes/key_id_pub.pem`

To run a new party node with a public key from the nodes directory type:
```
~$ ./client-node -new -id=1
```

### Poller

Poller is an application that polls the alfa node for a list of parties with the number of current votes and prints it to console output in an endless loop.

To run the poller type:
```
~$ ./poller
```

### Election

Election is an application that simulates voting process for all of the key-pairs it can find in the provided directory.

This application accepts a singe parameter:
- `clients` - directory of the key pairs for who to simulate the voting process; default value is `clients`

To the run the election application with default values:

```
~$ ./election
```

### Voter

Voter is an application that votes for a certain party during it's lifetime. It demonstrates an operation of a single voter. It is useful for debugging purposes

This application accepts 2 parameters:
1. `id` - id of the client that is voting, which is also the number of the key in `clients` directory
2. `choice` - number of the node for whom to vote which is also the number of the key in `nodes` directory

To run the voter with explicit parameters type:
```
~$ ./voter -id=1 -choice=1
```