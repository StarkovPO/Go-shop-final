package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

const (
	RunAddressKey           = "RUN_ADDRESS"
	DatabaseURIKey          = "DATABASE_URI"
	AccrualSystemAddressKey = "ACCRUAL_SYSTEM_ADDRESS"
)

type Config struct {
	RunAddressValue           string
	DatabaseURIValue          string
	AccrualSystemAddressValue string
}

func NewConfig() *Config {
	return &Config{}
}

func ReadFlags(c Config) error {
	rootCmd := &cobra.Command{
		Use:   "go-shop",
		Short: "[-a address], [-d database], [-r accrual system]",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("address: %s, database: %s, accrual system: %s\n", c.RunAddressValue, c.DatabaseURIValue, c.AccrualSystemAddressValue)
		},
	}
	rootCmd.Flags().StringVarP(&c.RunAddressValue, "Port for service", "a", ":8080", "Server address")
	rootCmd.Flags().StringVarP(&c.DatabaseURIValue, "URI for Postgres DB", "d", "postgres://admin:admin@localhost/go-shop?sslmode=disable", "Postgres URI")
	rootCmd.Flags().StringVarP(&c.AccrualSystemAddressValue, "ACCRUAL SYSTEM ADDRESS", "r", ":8000", "ACCRUAL_SYSTEM_ADDRESS")

	err := rootCmd.Execute()
	if err != nil {
		logrus.Fatalf("Unsuccessful attempt to read flags: %v", err)
		return err
	}

	return nil
}

func Init() (Config, error) {

	c := NewConfig()
	c.RunAddressValue = os.Getenv(RunAddressKey)
	c.DatabaseURIValue = os.Getenv(DatabaseURIKey)
	c.AccrualSystemAddressValue = os.Getenv(AccrualSystemAddressKey)

	if c.RunAddressValue == "" || c.DatabaseURIValue == "" || c.AccrualSystemAddressValue == "" {
		if err := ReadFlags(*c); err != nil {
			return Config{}, err
		}
	}
	return *c, nil
}
