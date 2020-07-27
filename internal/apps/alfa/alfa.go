package alfa

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/nebser/crypto-vote/internal/pkg/transaction"
	_websocket "github.com/nebser/crypto-vote/internal/pkg/websocket"

	_blockchain "github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/nebser/crypto-vote/internal/pkg/keyfiles"
	"github.com/nebser/crypto-vote/internal/pkg/wallet"
	"github.com/pkg/errors"
)

type Options struct {
	PrivateKeyFileName string
	PublicKeyFileName  string
	ClientKeysDir      string
}

func getKeyFiles(keyDirectory string) (keyfiles.KeyFilesList, error) {
	files, err := ioutil.ReadDir(keyDirectory)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read key file directory %s", keyDirectory)
	}

	fileGroups := map[string]keyfiles.KeyFiles{}
	for _, f := range files {
		if strings.Contains(f.Name(), "address") {
			continue
		}
		name := strings.Replace(f.Name(), "_pub", "", 1)
		group := fileGroups[name]
		if strings.Contains(f.Name(), "pub") {
			group.PublicKeyFile = fmt.Sprintf("%s/%s", keyDirectory, f.Name())
		} else {
			group.PrivateKeyFile = fmt.Sprintf("%s/%s", keyDirectory, f.Name())
		}
		fileGroups[name] = group
	}

	result := keyfiles.KeyFilesList{}
	for _, keyFiles := range fileGroups {
		result = append(result, keyFiles)
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

func Initialize(blockchain _blockchain.Blockchain, options Options) error {
	masterWallet, err := wallet.Import(keyfiles.KeyFiles{
		PublicKeyFile:  options.PublicKeyFileName,
		PrivateKeyFile: options.PrivateKeyFileName,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to import keys")
	}
	genesisTransaction, err := transaction.NewBaseTransaction(*masterWallet, masterWallet.Address)
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
	keyfiles, err := getKeyFiles(options.ClientKeysDir)
	if err != nil {
		return errors.Wrapf(err, "Failed to retrieve key files from directory %s", options.ClientKeysDir)
	}
	wallets, err := loadWallets(keyfiles)
	if err != nil {
		return errors.Wrap(err, "Failed to load wallets")
	}
	baseTransactions := transaction.Transactions{}
	for _, w := range wallets {
		t, err := transaction.NewBaseTransaction(*masterWallet, w.Address)
		if err != nil {
			return errors.Wrapf(err, "Failed to create transaction to wallet %#v", w)
		}
		baseTransactions = append(baseTransactions, *t)
	}
	block, err := _blockchain.NewBlock(blockchain.GetTip(), baseTransactions)
	if err != nil {
		return errors.Wrap(err, "Failed to create block of base transactions")
	}
	return blockchain.AddBlock(*block)

}

func Handler(b _blockchain.Blockchain) _websocket.Handler {
	return func(resp http.ResponseWriter, request *http.Request) error {
		// upgrader := websocket.Upgrader{}
		// conn, err := upgrader.Upgrade(resp, request, nil)
		// if err != nil {
		// 	return errors.Wrap(err, "Failed to open websocket")
		// }
		// defer conn.Close()

		// for {
		// 	var command _websocket.Command
		// 	if err := conn.ReadJSON(&command); err != nil {
		// 		return errors.Wrap(err, "Failed to parse json into command structure")
		// 	}

		// }
		return nil
	}
}
