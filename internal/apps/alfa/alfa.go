package alfa

import (
	"log"

	"github.com/nebser/crypto-vote/internal/pkg/transaction"
	"github.com/nebser/crypto-vote/internal/pkg/websocket"

	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/wallet"
	"github.com/pkg/errors"
)

func Initialize(masterWallet wallet.Wallet, clientWallets wallet.Wallets, initBlockchain blockchain.InitBlockchainFn, addBlock blockchain.AddBlockFn) error {
	genesisTransaction, err := transaction.NewBaseTransaction(masterWallet, masterWallet.Address)
	if err != nil {
		return errors.Wrap(err, "Failed to generate genesis transaction")
	}
	genesisBlock, err := blockchain.NewBlock(nil, transaction.Transactions{*genesisTransaction})
	if err != nil {
		return errors.Wrap(err, "Failed to create genesis block")
	}
	tip, err := initBlockchain(*genesisBlock)
	if err != nil {
		errors.Wrap(err, "Failed to initialize blockchain")
	}
	baseTransactions := transaction.Transactions{}
	for _, w := range clientWallets {
		t, err := transaction.NewBaseTransaction(masterWallet, w.Address)
		if err != nil {
			return errors.Wrapf(err, "Failed to create transaction to wallet %#v", w)
		}
		baseTransactions = append(baseTransactions, *t)
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
	log.Println("STARTED CHOOSING BLOCK FORGER")
	if err := r(); err != nil {
		log.Printf("Failed to run forge finder. Error %s", err)
	}
	log.Println("FINISHED CHOOSING BLOCK FORGER")
}

func Runner(registeredNodes websocket.RegisteredNodesFn, unicastRandomly websocket.RandomUnicastFn, getTip blockchain.GetTipFn, getBlock blockchain.GetBlockFn, signer wallet.Signer) RunnerFn {
	return func() error {
		if len(registeredNodes()) < 2 {
			return errors.Errorf("Not enough nodes registered to perform block forging. Number of blocks %d\n", len(registeredNodes()))
		}
		height, err := blockchain.GetHeight(getTip, getBlock)
		if err != nil {
			return errors.Errorf("Error occurred while trying to retrieve blockchain height %s", err)
		}
		pong, err := websocket.Pong{
			Message: websocket.ForgeBlockMessage,
			Body: websocket.ForgeBlockBody{
				Height: height,
			},
		}.Signed(signer)
		if err != nil {
			return errors.Errorf("Failed to sign forge block message %s", err)
		}
		if err := unicastRandomly(pong); err != nil {
			return errors.Errorf("Failed to send forge block message %s", err)
		}
		return nil
	}
}
