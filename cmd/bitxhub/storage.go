package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"path/filepath"
	"runtime"
	"time"

	"github.com/meshplus/bitxhub-core/agency"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/log"
	"github.com/meshplus/bitxhub-kit/storage"
	"github.com/meshplus/bitxhub-kit/storage/blockfile"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	"github.com/meshplus/bitxhub/internal/executor/contracts"
	"github.com/meshplus/bitxhub/internal/ledger"
	"github.com/meshplus/bitxhub/internal/loggers"
	"github.com/meshplus/bitxhub/internal/profile"
	"github.com/meshplus/bitxhub/internal/repo"
	"github.com/meshplus/bitxhub/internal/storages"
	"github.com/meshplus/bitxhub/pkg/vm"
	"github.com/meshplus/bitxhub/pkg/vm/boltvm"
	"github.com/urfave/cli"
)

var storeLogger = log.NewWithModule("cmd")

func storageCMD() cli.Command {
	return cli.Command{
		Name:  "storage",
		Usage: "Start a storage test",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "passwd",
				Usage:    "bitxhub key password",
				Required: false,
			},
			cli.StringFlag{
				Name:  "config",
				Usage: "bitxhub config path",
			},
			cli.StringFlag{
				Name:  "network",
				Usage: "bitxhub network config path",
			},
			cli.StringFlag{
				Name:  "order",
				Usage: "bitxhub order config path",
			},
		},
		Action: testStorage,
	}
}

func testStorage(ctx *cli.Context) error {
	repoRoot, err := repo.PathRootWithDefault(ctx.GlobalString("repo"))
	if err != nil {
		return fmt.Errorf("get repo path: %w", err)
	}
	passwd := ctx.String("passwd")
	configPath := ctx.String("config")
	networkPath := ctx.String("network")

	repo1, err := repo.Load(repoRoot, passwd, configPath, networkPath)
	if err != nil {
		return fmt.Errorf("repo load: %w", err)
	}

	err = log.Initialize(
		log.WithReportCaller(repo1.Config.Log.ReportCaller),
		log.WithPersist(true),
		log.WithFilePath(filepath.Join(repoRoot, repo1.Config.Log.Dir)),
		log.WithFileName(repo1.Config.Log.Filename),
		log.WithMaxAge(90*24*time.Hour),
		log.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return fmt.Errorf("log initialize: %w", err)
	}

	loggers.Initialize(repo1.Config)

	if err := storages.Initialize(repoRoot); err != nil {
		fmt.Println(fmt.Errorf("storages initialize: %w", err))
		return fmt.Errorf("storages initialize: %w", err)
	}

	bcStorage, err := storages.Get(storages.BlockChain)
	if err != nil {
		fmt.Println(fmt.Errorf("create blockchain storage: %w", err))
		return fmt.Errorf("create blockchain storage: %w", err)
	}

	stateStorage, err := ledger.OpenStateDB(repo.GetStoragePath(repoRoot, "ledger"), repo1.Config.Ledger.Type)
	if err != nil {
		fmt.Println(fmt.Errorf("create tm-leveldb: %w", err))
		return fmt.Errorf("create tm-leveldb: %w", err)
	}

	bf, err := blockfile.NewBlockFile(repoRoot, loggers.Logger(loggers.Storage))
	if err != nil {
		fmt.Println(fmt.Errorf("blockfile initialize: %w", err))
		return fmt.Errorf("blockfile initialize: %w", err)
	}

	rwLdg, err := ledger.New(repo1, bcStorage, stateStorage, bf, nil, loggers.Logger(loggers.Executor))
	if err != nil {
		fmt.Println(fmt.Errorf("create RW ledger: %w", err))
		return fmt.Errorf("create RW ledger: %w", err)
	}

	viewLdg := &ledger.Ledger{
		ChainLedger: rwLdg.ChainLedger,
	}
	viewLdg.StateLedger, err = ledger.NewSimpleLedger(repo1, stateStorage.(storage.Storage), nil, loggers.Logger(loggers.Executor))
	if err != nil {
		fmt.Println(err)
		return err
	}

	pprof, err := profile.NewPprof(repo1.Config)
	if err != nil {
		return err
	}
	if err := pprof.Start(); err != nil {
		return err
	}

	printVersion()
	testLogger := loggers.Logger(loggers.Storage)
	testLogger.Infoln("init several addresses...")
	addresses, err := initSeveralAddress(100)
	if err != nil {
		fmt.Println(err)
		return err
	}
	testLogger.Infoln("init address map...")
	addressMap := initAddressMap(addresses)
	if err != nil {
		fmt.Println(err)
		return err
	}
	testLogger.Infoln("start storage test...")

	contracts := registerBoltContracts()
	height := 2
	for {
		testLogger.Infoln("generate ibtp transactions...")
		txs := genTransactions(addresses, addressMap)
		testLogger.Infoln("begin block", height)
		block := &pb.Block{
			BlockHeader: &pb.BlockHeader{
				Number:      uint64(height),
				StateRoot:   types.NewHashByStr("0xc30B6E0ad5327fc8548f4BaFab3271cA6a5bD92f084095958c84970165bfA6E7"),
				TxRoot:      types.NewHashByStr("0xc30B6E0ad5327fc8548f4BaFab3271cA6a5bD92f084095958c84970165bfA6E7"),
				ReceiptRoot: types.NewHashByStr("0xc30B6E0ad5327fc8548f4BaFab3271cA6a5bD92f084095958c84970165bfA6E7"),
				ParentHash:  types.NewHashByStr("0xc30B6E0ad5327fc8548f4BaFab3271cA6a5bD92f084095958c84970165bfA6E7"),
				Timestamp:   time.Now().UnixNano(),
			},
			Transactions: &pb.Transactions{
				Transactions: txs,
			},
			BlockHash: types.NewHashByStr("0xc30B6E0ad5327fc8548f4BaFab3271cA6a5bD92f084095958c84970165bfA6E7"),
			Signature: []byte("111223123123213211121312312"),
			Extra:     []byte(""),
		}
		counter := make(map[string]*pb.VerifiedIndexSlice)
		for k, _ := range addressMap {
			var v []*pb.VerifiedIndex
			for i := 0; i < 5; i++ {
				v = append(v, &pb.VerifiedIndex{
					Index: uint64(i),
					Valid: true,
				})
			}
			counter[k] = &pb.VerifiedIndexSlice{Slice: v}
		}
		var l2Roots []types.Hash
		for i := 0; i < 200; i++ {
			l2Roots = append(l2Roots, *types.NewHashByStr("0xc30B6E0ad5327fc8548f4BaFab3271cA6a5bD92f084095958c84970165bfA6E7"))
		}
		var timeoutL2Roots []types.Hash
		for i := 0; i < 2; i++ {
			timeoutL2Roots = append(timeoutL2Roots, *types.NewHashByStr("0xc30B6E0ad5327fc8548f4BaFab3271cA6a5bD92f084095958c84970165bfA6E7"))
		}
		interchainMeta := &pb.InterchainMeta{
			Counter:        counter,
			L2Roots:        l2Roots,
			TimeoutCounter: make(map[string]*pb.StringSlice),
			TimeoutL2Roots: timeoutL2Roots,
			MultiTxCounter: make(map[string]*pb.StringSlice),
		}
		var receipts []*pb.Receipt
		var txHashList []*types.Hash
		for _, tx := range txs {
			txHashList = append(txHashList, tx.GetHash())
		}
		time1 := time.Now()
		for i, tx := range txs {
			vmCtx := vm.NewContext(tx, uint64(i), nil, uint64(height), rwLdg, testLogger, false)
			instance := boltvm.New(vmCtx, nil, nil, contracts)
			ret, err := instance.HandleIBTP(tx.GetIBTP())
			// testLogger.Infoln("handel success", i)
			if err != nil {
				fmt.Println(err)
				return err
			}
			receipt := &pb.Receipt{
				Version: tx.GetVersion(),
				TxHash:  tx.GetHash(),
			}
			receipt.Status = pb.Receipt_SUCCESS
			receipt.Ret = ret
			receipts = append(receipts, receipt)
			rwLdg.Finalise(true)
		}
		accounts, journalHash := rwLdg.FlushDirtyData()
		data := &ledger.BlockData{
			Block:          block,
			Receipts:       receipts,
			Accounts:       accounts,
			InterchainMeta: interchainMeta,
			TxHashList:     txHashList,
		}
		data.Block.BlockHeader.StateRoot = journalHash
		rwLdg.PersistBlockData(data)
		rwLdg.Clear()
		time2 := time.Now()
		exeTime := time2.Sub(time1)
		testLogger.Infoln(exeTime.Milliseconds())
		height++
		accounts = nil
		runtime.GC()
	}

	return nil
}

func initSeveralAddress(num int) ([]*types.Address, error) {
	var address []*types.Address
	for i := 0; i < num; i++ {
		privKey, err := asym.GenerateKeyPair(crypto.Secp256k1)
		if err != nil {
			return nil, fmt.Errorf("generate key: %w", err)
		}
		newAddress, err := privKey.PublicKey().Address()
		if err != nil {
			return nil, fmt.Errorf("generate address: %w", err)
		}
		address = append(address, newAddress)
	}

	return address, nil
}

func initAddressMap(addresses []*types.Address) map[string]uint64 {
	addressMap := make(map[string]uint64)
	for _, addr := range addresses {
		addressMap[addr.String()] = 1
	}

	return addressMap
}

func genTransactions(addresses []*types.Address, addressMap map[string]uint64) []pb.Transaction {
	var txs []pb.Transaction
	content := &pb.Content{
		Func: "interchainCharge",
		Args: [][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	}
	bytes, _ := content.Marshal()
	payload := &pb.Payload{
		Encrypted: false,
		Content:   bytes,
	}
	ibtppd, _ := payload.Marshal()
	proof := []byte("")
	for i := 0; i < 200; i++ {
		rand.Seed(time.Now().UnixNano())
		randIndex := rand.Intn(99)
		hash := sha256.Sum256([]byte(randString(20)))
		tx := &pb.BxhTransaction{
			From:            addresses[randIndex],
			To:              constant.InterchainContractAddr.Address(),
			TransactionHash: types.NewHash(hash[:]),
			IBTP: &pb.IBTP{
				From:          fmt.Sprintf("%s:%s:%s", "1356", addresses[randIndex], "transfer"),
				To:            fmt.Sprintf("%s:%s:%s", "1356", addresses[randIndex+1], "transfer"),
				Index:         addressMap[addresses[randIndex].String()],
				TimeoutHeight: 10,
				Payload:       ibtppd,
				Proof:         proof,
			},
		}
		txs = append(txs, tx)
		addressMap[addresses[randIndex].String()] = addressMap[addresses[randIndex].String()] + 1
	}
	return txs
}

func randString(len int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func registerBoltContracts() map[string]agency.Contract {
	boltContracts := []*boltvm.BoltContract{
		{
			Enabled:  true,
			Name:     "interchain manager contract",
			Address:  constant.InterchainContractAddr.Address().String(),
			Contract: &contracts.InterchainManager{},
		},
		{
			Enabled:  true,
			Name:     "store service",
			Address:  constant.StoreContractAddr.Address().String(),
			Contract: &contracts.Store{},
		},
		{
			Enabled:  true,
			Name:     "rule manager service",
			Address:  constant.RuleManagerContractAddr.Address().String(),
			Contract: &contracts.RuleManager{},
		},
		{
			Enabled:  true,
			Name:     "role manager service",
			Address:  constant.RoleContractAddr.Address().String(),
			Contract: &contracts.RoleManager{},
		},
		{
			Enabled:  true,
			Name:     "appchain manager service",
			Address:  constant.AppchainMgrContractAddr.Address().String(),
			Contract: &contracts.AppchainManager{},
		},
		{
			Enabled:  true,
			Name:     "transaction manager service",
			Address:  constant.TransactionMgrContractAddr.Address().String(),
			Contract: &contracts.TransactionManager{},
		},
		{
			Enabled:  true,
			Name:     "governance service",
			Address:  constant.GovernanceContractAddr.Address().String(),
			Contract: &contracts.Governance{},
		},
		{
			Enabled:  true,
			Name:     "node manager service",
			Address:  constant.NodeManagerContractAddr.Address().String(),
			Contract: &contracts.NodeManager{},
		},
		{
			Enabled:  true,
			Name:     "inter broker service",
			Address:  constant.InterBrokerContractAddr.Address().String(),
			Contract: &contracts.InterBroker{},
		},
		{
			Enabled:  true,
			Name:     "service manager service",
			Address:  constant.ServiceMgrContractAddr.Address().String(),
			Contract: &contracts.ServiceManager{},
		},
		{
			Enabled:  true,
			Name:     "dapp manager service",
			Address:  constant.DappMgrContractAddr.Address().String(),
			Contract: &contracts.DappManager{},
		},
	}

	ContractsInfo := agency.GetRegisteredContractInfo()
	for addr, info := range ContractsInfo {
		boltContracts = append(boltContracts, &boltvm.BoltContract{
			Enabled:  true,
			Name:     info.Name,
			Address:  addr,
			Contract: info.Constructor(),
		})
	}

	return boltvm.Register(boltContracts)
}
