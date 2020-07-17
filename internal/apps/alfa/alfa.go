package alfa

import (
	"io/ioutil"
	"strings"

	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/transaction"

	_blockchain "github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/keyfiles"
	"github.com/nebser/crypto-vote/internal/pkg/wallet"
	"github.com/pkg/errors"
)

type Options struct {
	New                bool
	PrivateKeyFileName string
	PublicKeyFileName  string
	ClientKeysDir      string
}

func getKeyFiles(keyDirectory string) (keyfiles.KeyFilesList, error) {
	files, err := ioutil.ReadDir(keyDirectory)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read key file directory %s", keyDirectory)
	}

	fileNames := []string{}
	for _, f := range files {
		if !strings.Contains(f.Name(), "address") {
			fileNames = append(fileNames, f.Name())
		}
	}

	result := keyfiles.KeyFilesList{}
	for i := 0; i < len(fileNames); i += 2 {
		result = append(result, keyfiles.KeyFiles{
			PrivateKeyFile: fileNames[i],
			PublicKeyFile:  fileNames[i+1],
		})
	}
	return result, nil
}

func loadWallets(fileList keyfiles.KeyFilesList) (wallet.Wallets, error) {
	result := wallet.Wallets{}
	for _, k := range fileList {
		w, err := wallet.Import(k)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to import keys")
		}
		result = append(result, *w)
	}
	return result, nil
}

func Initialize(getBlockchain blockchain.GetBlockchainFn, initBlockchain blockchain.InitBlockchainFn, addBlock blockchain.AddBlockFn, options Options) (*_blockchain.Blockchain, error) {
	masterWallet, err := wallet.Import(keyfiles.KeyFiles{
		PublicKeyFile:  options.PublicKeyFileName,
		PrivateKeyFile: options.PrivateKeyFileName,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to import keys")
	}
	if !options.New {
		return getBlockchain()
	}
	genesisTransaction, err := transaction.NewBaseTransaction(*masterWallet, masterWallet.Address)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate genesis transaction")
	}
	genesisBlock, err := _blockchain.NewBlock(nil, transaction.Transactions{*genesisTransaction})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create genesis block")
	}
	blockchain, err := initBlockchain(*genesisBlock)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize blockchain")
	}
	keyfiles, err := getKeyFiles(options.ClientKeysDir)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to retrieve key files from directory %s", options.ClientKeysDir)
	}
	wallets, err := loadWallets(keyfiles)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load wallets")
	}
	baseTransactions := transaction.Transactions{}
	for _, w := range wallets {
		t, err := transaction.NewBaseTransaction(*masterWallet, w.Address)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to create transaction to wallet %#v", w)
		}
		baseTransactions = append(baseTransactions, *t)
	}
	block, err := _blockchain.NewBlock(blockchain.Tip, baseTransactions)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create block of base transactions")
	}
	return addBlock(*blockchain, *block)

}
