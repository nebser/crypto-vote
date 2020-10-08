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

1. alfa - directory in which to create key pair for the alfa node; default value is "alfa"
2. clients - directory in which to create key pairs for clients (voters); default value is "clients"
3. nodes - directory in which to create key pairs for party nodes; default value is "nodes"
4. clientsNumber - number of key pairs to create for clients (voters); default value is 50
5. nodesNumber - number of key pairs to create for nodes; default value is 5  

To run it with default values:
```
~$ ./key-generator
```
