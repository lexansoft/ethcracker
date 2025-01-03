package keystore

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	//"io/ioutil"
	//    "github.com/pborman/uuid"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"sync"

	//    "strconv"
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

type CrackerParams struct {
	key_version string
	key_v1      *encryptedKeyJSONV1
	key_v3      *encryptedKeyJSONV3

	// for presale
	iv         []byte
	cipherText []byte
	EthAddr    string

	V          int // Verbosity
	Start_from int

	N       int
	Total   int
	RE      int
	Skipped int

	StartTime time.Time
}

var mutex = &sync.Mutex{}

func LoadPresaleFile(params *CrackerParams, path string) error {

	preSaleKeyStruct := struct {
		EncSeed string
		EthAddr string
		Email   string
		BtcAddr string
	}{}

	keyFileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	params.key_version = "presale"

	err = json.Unmarshal(keyFileContent, &preSaleKeyStruct)
	if err != nil {
		return err
	}

	params.EthAddr = preSaleKeyStruct.EthAddr

	encSeedBytes, err := hex.DecodeString(preSaleKeyStruct.EncSeed)
	params.iv = encSeedBytes[:16]
	params.cipherText = encSeedBytes[16:]

	println("Key file version:", params.key_version)
	return nil
}

func LoadKeyFile(params *CrackerParams, path string, log_level int) error {

	keyFileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	//    println( "Private key file content:", string( keyFileContent ) )

	params.key_version = "v3"
	params.key_v3, err = LoadKeyVersion3(keyFileContent)

	if err != nil {
		params.key_version = "v1"
		params.key_v1, err = LoadKeyVersion1(keyFileContent)
		if err != nil {
			return err
		}
		//println( "Private key JSON (version 1):", string( pk_log ) )
	} else {
		//println( "Private key JSON (version 3):", string( pk_log ) )
	}

	if log_level > 0 {
		println("Key file version:", params.key_version)
	}
	return nil
}

func LoadKeyVersion1(fileContent []byte) (*encryptedKeyJSONV1, error) {
	key := new(encryptedKeyJSONV1)
	err := json.Unmarshal(fileContent, key)
	return key, err
}

func LoadKeyVersion3(fileContent []byte) (*encryptedKeyJSONV3, error) {
	key := new(encryptedKeyJSONV3)
	err := json.Unmarshal(fileContent, key)
	return key, err
}

func Test_pass(params *CrackerParams, s string, thread int) error {
	var err error

	mutex.Lock()
	params.N++
	if params.V > 0 {
		if params.N+params.Skipped > params.Start_from {
			h := time.Since(params.StartTime).Hours() *
				float64(params.Total-(params.N+params.Skipped)) / float64(params.N+params.Skipped-params.Start_from)

			if params.N%params.RE == 0 {
				fmt.Printf("TH%d-> %d/%d %d%% Skipped: %d Left: %d years %d days %d hours %d minutes %v\n",
					thread,
					params.N+params.Skipped,
					params.Total,
					(params.N+params.Skipped)*100/params.Total,
					params.Skipped,
					int64(h)/(24*365), (int64(h)%(24*365))/24, int64(h)%24, int64(h*60)%60,
					s)
			}

		} else {
			if params.Start_from > 0 && (params.N+params.Skipped)%(100000) == 0 {

				h := time.Since(params.StartTime).Hours() *
					float64(params.Start_from-(params.N+params.Skipped)) / float64(params.Start_from)

				fmt.Printf("Skipping first %d -> %d %d%% Skipped: %d Left: %d years %d days %d hours %d minutes %v\n",
					params.Start_from,
					params.N+params.Skipped,
					(params.N+params.Skipped)*100/params.Start_from,
					params.Skipped,
					int64(h)/(24*365), (int64(h)%(24*365))/24, int64(h)%24, int64(h*60)%60,
					s)

			}
		}
	}
	mutex.Unlock()
	if params.N+params.Skipped < params.Start_from {
		return errors.New("skipped")
	}

	var addr string
	var pk []byte
	switch params.key_version {
	case "v3":
		err, addr, pk = Test_pass_v3(params.key_v3, s)
	case "v1":
		err, addr, pk = Test_pass_v1(params.key_v1, s)
	case "presale":
		err, addr, pk = Test_pass_presale(params, s)
	}

	if err == nil {
		println("")
		println("")
		println("-------------------------------------------------------------------------")
		println("              CONGRATULATIONS !!! WE FOUND YOUR PASSWORD !!!")
		println("-------------------------------------------------------------------------")
		println("")
		println("   Password:", s)
		println("    Address:", addr)
		println("Private Key:", hex.EncodeToString(pk))
		println("")
		println("          Do not forget to donate some ETH to the developer:")
		println("     Ethereum Address: 0x281694Fabfdd9735e01bB59942B18c469b6e3df6")
		println("-------------------------------------------------------------------------")
		println("")
		println("")
		os.Exit(0)
	}

	return err
}

func Test_pass_v1(k *encryptedKeyJSONV1, auth string) (error, string, []byte) {
	pk, _, err := decryptKeyV1(k, auth)
	return err, k.Address, pk
}

func Test_pass_v3(k *encryptedKeyJSONV3, auth string) (error, string, []byte) {
	pk, _, err := decryptKeyV3(k, auth)
	return err, k.Address, pk
}

func Test_pass_presale(params *CrackerParams, password string) (error, string, []byte) {
	passBytes := []byte(password)
	derivedKey := pbkdf2.Key(passBytes, passBytes, 2000, 16, sha256.New)
	//	_, err := aesCBCDecrypt(derivedKey, params.cipherText, params.iv)

	plainText, err := aesCBCDecrypt(derivedKey, params.cipherText, params.iv)
	if err != nil {
		return err, "", nil
	}
	ethPriv := crypto.Keccak256(plainText) // Sha3(plainText)
	ecKey, err := crypto.ToECDSA(ethPriv)
	if err != nil {
		return err, "", nil
	}
	//	key = &Key {
	//		Id:         nil,
	//		Address:    PubkeyToAddress( ecKey.PublicKey),
	//		PrivateKey: ecKey,
	//	}

	Address := crypto.PubkeyToAddress(ecKey.PublicKey)

	//	derivedAddr := hex.EncodeToString(key.Address.Bytes()) // needed because .Hex() gives leading "0x"
	derivedAddr := hex.EncodeToString(Address.Bytes()) // needed because .Hex() gives leading "0x"
	expectedAddr := params.EthAddr
	if derivedAddr != expectedAddr {
		err = fmt.Errorf("decrypted addr '%s' not equal to expected addr '%s'", derivedAddr, expectedAddr)
	}

	return err, derivedAddr, ethPriv
}
