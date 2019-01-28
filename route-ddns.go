package main

import (
	"container/ring"
	"fmt"
	"log"
	"os"
	"time"

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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./route-ddns.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("route-ddns")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
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

		ticker := time.NewTicker(20 * time.Second)
		for ; true; <-ticker.C {
			var provider = ipProvidersRing.Value.(string)
			ipProvidersRing = ipProvidersRing.Next()
			var resolvedIP, err = resolveIP(provider)

			if err != nil {
				log.Print(fmt.Sprintf("ERROR: Max attempts reached while resolving from %v", provider))
				continue
			}

			if resolvedIP != currentIP {
				log.Print(fmt.Sprintf("IP '%v' resolved from '%v' does not match current IP '%v'", resolvedIP, provider, currentIP))
				err = changeIP(resolvedIP, cfg)
				if err != nil {
					log.Print(fmt.Sprintf("ERROR: Max attempts reached while applying DNS changes"))
				}
				log.Print(fmt.Sprintf("DNS records updated"))
				currentIP = resolvedIP
			} else {
				log.Print(fmt.Sprintf("IP '%v' resolved from '%v' matches current IP '%v'", resolvedIP, provider, currentIP))
			}
		}
	},
}
