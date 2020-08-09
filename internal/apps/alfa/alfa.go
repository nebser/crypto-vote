package alfa

import (
	"github.com/nebser/crypto-vote/internal/pkg/transaction"

	_blockchain "github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/wallet"
	"github.com/pkg/errors"
)

func Initialize(masterWallet wallet.Wallet, clientWallets wallet.Wallets, blockchain _blockchain.Blockchain, saveNode _blockchain.SaveNodeFn) error {
	genesisTransaction, err := transaction.NewBaseTransaction(masterWallet, masterWallet.Address)
	if err != nil {
		return errors.Wrap(err, "Failed to generate genesis transaction")
	}
	genesisBlock, err := _blockchain.NewBlock(nil, transaction.Transactions{*genesisTransaction})
	if err != nil {
		return errors.Wrap(err, "Failed to create genesis block")
	}
	if err := blockchain.SetGenesis(*genesisBlock); err != nil {
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
	block, err := _blockchain.NewBlock(blockchain.GetTip(), baseTransactions)
	if err != nil {
		return errors.Wrap(err, "Failed to create block of base transactions")
	}
	if err := saveNode(_blockchain.Node{
		Type: _blockchain.AlfaNodeType,
		ID:   "0",
	}); err != nil {
		return errors.Wrap(err, "Failed to create record for alfa node")
	}
	return blockchain.AddBlock(*block)

}
