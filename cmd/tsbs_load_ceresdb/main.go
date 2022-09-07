package main

import (
	"fmt"
	"log"

	"github.com/blagojts/viper"
	"github.com/spf13/pflag"
	"github.com/timescale/tsbs/internal/utils"
	"github.com/timescale/tsbs/load"
	"github.com/timescale/tsbs/pkg/data/source"
	"github.com/timescale/tsbs/pkg/targets/ceresdb"
)

func initProgramOptions() (*ceresdb.SpecificConfig, load.BenchmarkRunner, *load.BenchmarkRunnerConfig) {
	target := ceresdb.NewTarget()

	loaderConf := load.BenchmarkRunnerConfig{}
	loaderConf.AddToFlagSet(pflag.CommandLine)
	target.TargetSpecificFlags("", pflag.CommandLine)
	pflag.Parse()

	if err := utils.SetupConfigFile(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	if err := viper.Unmarshal(&loaderConf); err != nil {
		panic(fmt.Errorf("unable to decode config: %s", err))
	}
	ceresdbAddr := viper.GetString("ceresdbAddr")
	if len(ceresdbAddr) == 0 {
		log.Fatalf("missing `ceresdbAddr` flag")
	}

	loader := load.GetBenchmarkRunner(loaderConf)
	return &ceresdb.SpecificConfig{CeresdbAddr: ceresdbAddr}, loader, &loaderConf
}

func main() {
	vmConf, loader, loaderConf := initProgramOptions()
	benchmark, err := ceresdb.NewBenchmark(vmConf, &source.DataSourceConfig{
		Type: source.FileDataSourceType,
		File: &source.FileDataSourceConfig{Location: loaderConf.FileName},
	})

	if err != nil {
		panic(err)
	}
	loader.RunBenchmark(benchmark)
}