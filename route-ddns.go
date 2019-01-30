package main

import (
	"container/ring"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cfg Config

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// CLI
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ./route-ddns.yaml)")

	// Set config defaults
	viper.SetDefault("cycleTime", 300)
	viper.SetDefault("providers", []string{"https://api.ipify.org?format=text"})
}

func initConfig() {
	// Setup config file location
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("route-ddns")
	}
	// Setup ENV configuration
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err == nil {
		log.Print("Using config file: ", viper.ConfigFileUsed())
	}

	// Marshall config file to Config struct
	viper.Unmarshal(&cfg)
}

var rootCmd = &cobra.Command{
	Use:   "route-ddns",
	Short: "A ddns client that updates AWS route53 records.",
	Long:  `route-ddns is a ddns client that updates record sets in AWS route53 when an IP change is detected.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create ring of providers to distibute load across these free services
		var numberOfProviders = len(cfg.Providers)
		ipProvidersRing := ring.New(numberOfProviders)
		for i := 0; i < numberOfProviders; i++ {
			ipProvidersRing.Value = cfg.Providers[i]
			ipProvidersRing = ipProvidersRing.Next()
		}

		var currentIP = ""

		// Main logic loop
		ticker := time.NewTicker(time.Duration(cfg.CycleTime) * time.Second)
		for ; true; <-ticker.C {
			// Load provider
			var provider = ipProvidersRing.Value.(string)
			ipProvidersRing = ipProvidersRing.Next()

			// Resolve IP
			var resolvedIP, err = resolveIP(provider)
			if err != nil {
				continue
			}

			// Change DNS record if necessary
			if resolvedIP != currentIP {
				log.Info(fmt.Sprintf("IP '%v' resolved from '%v' does not match current IP '%v'", resolvedIP, provider, currentIP))
				changeIP(resolvedIP, cfg)
				currentIP = resolvedIP
			} else {
				log.Info(fmt.Sprintf("IP '%v' resolved from '%v' matches current IP '%v'", resolvedIP, provider, currentIP))
			}
		}
	},
}
