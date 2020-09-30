package alfa

import (
	"fmt"
	"log"

	"github.com/nebser/crypto-vote/internal/pkg/party"
	"github.com/nebser/crypto-vote/internal/pkg/transaction"
	"github.com/nebser/crypto-vote/internal/pkg/websocket"

	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/wallet"
	"github.com/pkg/errors"
)

func Initialize(masterWallet wallet.Wallet, nodeWallets, clientWallets wallet.Wallets, addBlock blockchain.AddBlockFn, saveParty party.SavePartyFn) error {
	genesisTransaction, err := transaction.NewBaseTransaction(masterWallet, masterWallet.Address, 100*transaction.VoteValue)
	if err != nil {
		return errors.Wrap(err, "Failed to generate genesis transaction")
	}
	genesisBlock, err := blockchain.NewBlock(nil, transaction.Transactions{*genesisTransaction})
	if err != nil {
		return errors.Wrap(err, "Failed to create genesis block")
	}
	tip, err := addBlock(*genesisBlock)
	if err != nil {
		errors.Wrap(err, "Failed to initialize blockchain")
	}
	baseTransactions := transaction.Transactions{}
	for _, w := range append(nodeWallets, clientWallets...) {
		t, err := transaction.NewBaseTransaction(masterWallet, w.Address, transaction.VoteValue)
		if err != nil {
			return errors.Wrapf(err, "Failed to create transaction to wallet %#v", w)
		}
		baseTransactions = append(baseTransactions, *t)
	}
	for i, wallet := range nodeWallets {
		p := party.Party{
			Name:    fmt.Sprintf("Party Number: %d", i),
			Address: wallet.Address,
		}
		if err := saveParty(p); err != nil {
			return errors.Wrapf(err, "Failed to save party %#v", p)
		}
	}
	block, err := blockchain.NewBlock(tip, baseTransactions)
	if err != nil {
		return errors.Wrap(err, "Failed to create block of base transactions")
	}
	if _, err := addBlock(*block); err != nil {
		return errors.Wrapf(err, "Failed to add block %#v", *block)
	}
	return nil

}

type RunnerFn func() error

func (r RunnerFn) Run() {
	log.Println("STARTED RUNNER")
	if err := r(); err != nil {
		log.Printf("Failed to run runner. Error %s", err)
	}
	log.Println("FINISHED RUNNER")
}

func Runner(registeredNodes websocket.RegisteredNodesFn, unicastRandomly websocket.RandomUnicastFn, getTip blockchain.GetTipFn, getBlock blockchain.GetBlockFn) RunnerFn {
	return func() error {
		if len(registeredNodes()) < 2 {
			return errors.Errorf("Not enough nodes registered to perform block forging. Number of blocks %d\n", len(registeredNodes()))
		}
		height, err := blockchain.GetHeight(getTip, getBlock)
		if err != nil {
			return errors.Errorf("Error occurred while trying to retrieve blockchain height %s", err)
		}
		pong := websocket.Pong{
			Message: websocket.ForgeBlockMessage,
			Body: websocket.ForgeBlockBody{
				Height: height,
			},
		}
		if err != nil {
			return errors.Errorf("Failed to sign forge block message %s", err)
		}
		if err := unicastRandomly(pong); err != nil {
			return errors.Errorf("Failed to send forge block message %s", err)
		}
		return nil
	}
}

func Cleaner(
	getTransactions transaction.GetTransactionsFn,
	isReturnStakeTransaction transaction.IsReturnStakeTransactionFn,
	getTip blockchain.GetTipFn,
	getBlock blockchain.GetBlockFn,
	addBlock blockchain.AddBlockFn,
	broadcast websocket.BroadcastFn,
) RunnerFn {
	return func() error {
		txs, err := getTransactions()
		if err != nil {
			return errors.Wrap(err, "Failed to retrieve transactions")
		}
		log.Printf("Found transactions %d", len(txs))
		log.Printf("Found transactions %s", txs)
		if len(txs) != 1 || !isReturnStakeTransaction(txs[0]) {
			log.Println("Cleaner unnecessary")
			return nil
		}
		height, err := blockchain.GetHeight(getTip, getBlock)
		if err != nil {
			return errors.Wrap(err, "Failed to retrieve blockchain height")
		}
		block, err := blockchain.NewBlock(getTip(), transaction.Transactions{txs[0]})
		if err != nil {
			return errors.Wrap(err, "Failed to create new block")
		}
		if _, err := addBlock(*block); err != nil {
			return errors.Wrapf(err, "Failed to add block to blockchain")
		}
		broadcast(websocket.Pong{
			Message: websocket.BlockForgedMessage,
			Body: websocket.BlockForgedBody{
				Height: height + 1,
				Block:  *block,
			},
		})
		return nil
	}
}
