package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "embed"

	log "github.com/sirupsen/logrus"

	"dario.cat/mergo"
	"github.com/santhosh-tekuri/jsonschema/v5"

	"gopkg.in/yaml.v3"
)

type TaxCategoryType string

const (
	Debt           TaxCategoryType = "debt"
	Equity         TaxCategoryType = "equity"
	Equity65       TaxCategoryType = "equity65"
	Equity35       TaxCategoryType = "equity35"
	UnlistedEquity TaxCategoryType = "unlisted_equity"
)

type CommodityType string

const (
	MutualFund CommodityType = "mutualfund"
	NPS        CommodityType = "nps"
	Stock      CommodityType = "stock"
	Unknown    CommodityType = "unknown"
)

type Price struct {
	Provider string `json:"provider" yaml:"provider"`
	Code     string `json:"code" yaml:"code"`
}

type Commodity struct {
	Name        string          `json:"name" yaml:"name"`
	Type        CommodityType   `json:"type" yaml:"type"`
	Price       Price           `json:"price" yaml:"price"`
	Harvest     int             `json:"harvest" yaml:"harvest"`
	TaxCategory TaxCategoryType `json:"tax_category" yaml:"tax_category"`
}

type Retirement struct {
	SWR            float64  `json:"swr" yaml:"swr"`
	Expenses       []string `json:"expenses" yaml:"expenses"`
	Savings        []string `json:"savings" yaml:"savings"`
	YearlyExpenses float64  `json:"yearly_expenses" yaml:"yearly_expenses"`
}

type ScheduleAL struct {
	Code     string   `json:"code" yaml:"code"`
	Accounts []string `json:"accounts" yaml:"accounts"`
}

type AllocationTarget struct {
	Name     string   `json:"name" yaml:"name"`
	Target   float64  `json:"target" yaml:"target"`
	Accounts []string `json:"accounts" yaml:"accounts"`
}

type Config struct {
	JournalPath                string     `json:"journal_path" yaml:"journal_path"`
	DBPath                     string     `json:"db_path" yaml:"db_path"`
	LedgerCli                  string     `json:"ledger_cli" yaml:"ledger_cli"`
	DefaultCurrency            string     `json:"default_currency" yaml:"default_currency"`
	Locale                     string     `json:"locale" yaml:"locale"`
	FinancialYearStartingMonth time.Month `json:"financial_year_starting_month" yaml:"financial_year_starting_month"`

	Retirement Retirement `json:"retirement" yaml:"retirement"`

	ScheduleALs []ScheduleAL `json:"schedule_al" yaml:"schedule_al"`

	AllocationTargets []AllocationTarget `json:"allocation_targets" yaml:"allocation_targets"`

	Commodities []Commodity `json:"commodities" yaml:"commodities"`
}

var config Config
var configPath string

var defaultConfig = Config{
	LedgerCli:                  "ledger",
	DefaultCurrency:            "INR",
	Locale:                     "en-IN",
	Retirement:                 Retirement{SWR: 4, Savings: []string{"Assets:*"}, Expenses: []string{"Expenses:*"}, YearlyExpenses: 0},
	FinancialYearStartingMonth: 4,
	ScheduleALs:                []ScheduleAL{},
	AllocationTargets:          []AllocationTarget{},
	Commodities:                []Commodity{},
}

//go:embed schema.json
var SchemaJson string
var schema *jsonschema.Schema

func init() {
	schema = jsonschema.MustCompileString("", SchemaJson)
}

func SaveConfig(content []byte) error {
	err := LoadConfig(content, "")
	if err != nil {
		return err
	}

	yamlContent, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, yamlContent, 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadConfigFile(path string) {
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		log.Warn("Failed to read config file: ", path)
		log.Fatal(err)
	}

	err = LoadConfig(content, path)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Using config file: ", path)
}

func LoadConfig(content []byte, cp string) error {
	var configJson interface{}
	err := yaml.Unmarshal(content, &configJson)
	if err != nil {
		return err
	}

	err = schema.Validate(configJson)
	if err != nil {
		return errors.New(fmt.Sprintf("Invalid configuration\n%#v", err))
	}

	config = Config{}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return err
	}

	err = mergo.Merge(&config, defaultConfig, mergo.WithOverrideEmptySlice)

	if err != nil {
		return err
	}

	journalDir := filepath.Dir(configPath)

	if !filepath.IsAbs(config.JournalPath) {
		config.JournalPath = filepath.Join(journalDir, config.JournalPath)
	}

	if !filepath.IsAbs(config.DBPath) {
		config.DBPath = filepath.Join(journalDir, config.DBPath)
	}

	if cp != "" && configPath == "" {
		configPath = cp
	}

	return nil
}

func GetConfig() Config {
	return config
}

func GetConfigDir() string {
	return filepath.Dir(configPath)
}

func GetConfigPath() string {
	return configPath
}

func GetSchema() any {
	var schemaObject any
	err := json.Unmarshal([]byte(SchemaJson), &schemaObject)
	if err != nil {
		log.Fatal(err)
	}
	return schemaObject
}

func EnsureLogFilePath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(cacheDir, "paisa", "paisa.log")

	err = os.MkdirAll(filepath.Dir(path), 0750)
	if err != nil {
		return "", err
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return "", err
	}

	err = file.Close()
	if err != nil {
		return "", err
	}

	return path, err
}

func DefaultCurrency() string {
	return config.DefaultCurrency
}
