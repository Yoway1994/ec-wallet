package walletservice

import (
	"ec-wallet/internal/domain/wallet"
	"ec-wallet/internal/errors"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// 產生助記詞
func GenMnemonic() (mnemonic string, err error) {
	// 產生亂數
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return
	}
	// 從亂數產生註記詞
	mnemonic, err = bip39.NewMnemonic(entropy)
	if err != nil {
		return
	}
	return
}

// DeriveKeyFromPath generates a key pair from a specific derivation path
func (s *walletService) DeriveKeyFromPath(mnemonic string, hdPath *wallet.HDPath) (*wallet.KeyPair, error) {
	if mnemonic == "" {
		return nil, errors.ErrWalletMnemonicRequired
	}

	// Generate seed
	seed := bip39.NewSeed(mnemonic, "")

	// Generate master keykey
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}

	// Derive child keys according to path
	child := masterKey
	for _, component := range hdPath.Components {
		child, err = child.NewChildKey(component)
		if err != nil {
			return nil, err
		}
	}

	// Generate secp256k1 key pair
	privateKey := secp256k1.PrivKeyFromBytes(child.Key)
	publicKey := privateKey.PubKey()

	// Generate Ethereum address
	address := crypto.PubkeyToAddress(*publicKey.ToECDSA()).Hex()

	return &wallet.KeyPair{
		PrivateKey: child.Key,
		PublicKey:  publicKey.SerializeCompressed(),
		Address:    address,
	}, nil
}
