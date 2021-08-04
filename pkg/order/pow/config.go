package pow

import (
	"github.com/spf13/viper"
	"math/big"
	"path/filepath"
	"time"
)

const (
	Difficulty = 1
	BlockNonce = 0
)

var (
	// two256 is a big integer representing 2^256
	two256 = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0))
)

type POWConfig struct {
	POW POW
}

type POW struct {
	BatchTimeout  time.Duration `mapstructure:"batch_timeout"`
	MempoolConfig MempoolConfig `mapstructure:"mempool"`
}

type MempoolConfig struct {
	BatchSize      uint64        `mapstructure:"batch_size"`
	PoolSize       uint64        `mapstructure:"pool_size"`
	TxSliceSize    uint64        `mapstructure:"tx_slice_size"`
	TxSliceTimeout time.Duration `mapstructure:"tx_slice_timeout"`
}

func generatePowConfig(repoRoot string) (time.Duration, MempoolConfig, error) {
	readConfig, err := readConfig(repoRoot)
	if err != nil {
		return 0, MempoolConfig{}, err
	}
	mempoolConf := MempoolConfig{}
	mempoolConf.BatchSize = readConfig.POW.MempoolConfig.BatchSize
	mempoolConf.PoolSize = readConfig.POW.MempoolConfig.PoolSize
	mempoolConf.TxSliceSize = readConfig.POW.MempoolConfig.TxSliceSize
	mempoolConf.TxSliceTimeout = readConfig.POW.MempoolConfig.TxSliceTimeout
	return readConfig.POW.BatchTimeout, mempoolConf, nil
}

func readConfig(repoRoot string) (*POWConfig, error) {
	v := viper.New()
	v.SetConfigFile(filepath.Join(repoRoot, "order.toml"))
	v.SetConfigType("toml")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &POWConfig{}

	if err := v.Unmarshal(config); err != nil {
		return nil, err
	}
	return config, nil
}
