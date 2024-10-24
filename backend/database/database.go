package database

import (
	"log"
	"time"

	"github.com/plaid/plaid-go/v21/plaid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id       string
	Username string `gorm:"username"`
}

type Account struct {
	gorm.Model
	UserID             string `gorm:"user_id"`
	AccountType        string // one of: depository, credit, loan, investment
	AccountSubtype     string
	InstitutionID      string
	Balance            float64
	CreditLimit        float64
	CashBalance        float64
	PlaidAccessToken   string
	AccountID          string
	TransactionsCursor string
}

type Holding struct {
	gorm.Model
	AccountID string
	Symbol    string
	Quantity  float64
}

type Transaction struct {
	gorm.Model
	AccountID       string  `gorm:"account_id"`
	Amount          float64 `gorm:"amount"`
	IsoCurrencyCode string  `gorm:"iso_currency_code"`
	Category        string  `gorm:"category"`
	Name            string  `gorm:"name"`
	PaymentChannel  string  `gorm:"payment_channel"`
	MerchantName    string  `gorm:"merchant_name"`
	Pending         bool    `gorm:"pending"`
	TransactionID   string  `gorm:"transaction_id",unique`
	Date            string  `gorm:"date"`
	UserID          string
}

type NetWorthResponse struct {
	gorm.Model
	NetWorth float64 `gorm:"netWorth"`
	Datetime time.Time
}

var Instance *gorm.DB
var err error

func Connect() {
	Instance, err = gorm.Open(sqlite.Open("../db/gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to Database")
}

func Migrate() {
	Instance.AutoMigrate(&User{})
	Instance.AutoMigrate(&Account{})
	Instance.AutoMigrate(&Transaction{})
	Instance.AutoMigrate(&Holding{})
	Instance.AutoMigrate(&NetWorthResponse{})
}

func GetAccountsForUser(userId string) []Account {
	var accounts []Account
	Instance.Find(&accounts, "user_id = ?", userId)
	return accounts
}

func GetTransactionsForUser(userId string) []Transaction {
	var transactions []Transaction
	Instance.Find(&transactions, "user_id = ?", userId)
	return transactions

}

func GetAccountById(accountId string) Account {
	var acc Account
	Instance.First(&acc, Account{AccountID: accountId})
	return acc
}

func GetTransactionByID(transactionId string) Transaction {
	var txn Transaction
	Instance.First(&txn, Transaction{TransactionID: transactionId})
	return txn
}

func UpdateCursor(accountID string, cursor string) {
	var acc Account
	if err := Instance.Where("account_id = ?", accountID).First(&acc).Error; err != nil {
		log.Println("Error while updating cursor:", err)
		return
	}
	acc.TransactionsCursor = cursor
	if err := Instance.Save(&acc).Error; err != nil {
		log.Println("Error while saving cursor:", err)
	}
}

func GetLatestCursorOrNil(accountID string) string {
	var acc Account
	if err := Instance.Where("account_id = ?", accountID).First(&acc).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("Account not found while retrieving cursor")
			return ""
		}
		log.Println("Error while retrieving cursor:", err)
		return ""
	}
	return acc.TransactionsCursor
}

func ConvertPlaidTransaction(plaidTransaction plaid.Transaction, accountId string) Transaction {
	account := GetAccountById(accountId)

	var isoCurrencyCode string
	if plaidTransaction.IsoCurrencyCode.IsSet() {
		isoCurrencyCode = *plaidTransaction.IsoCurrencyCode.Get()
	}

	var merchantName string
	if plaidTransaction.MerchantName.IsSet() {
		merchantName = *plaidTransaction.MerchantName.Get()
	}

	var category string
	if plaidTransaction.PersonalFinanceCategory.IsSet() {
		category = plaidTransaction.PersonalFinanceCategory.Get().Primary
	}

	return Transaction{
		AccountID:       account.AccountID,
		Amount:          plaidTransaction.Amount,
		IsoCurrencyCode: isoCurrencyCode,
		Category:        category,
		Name:            plaidTransaction.Name,
		PaymentChannel:  plaidTransaction.PaymentChannel,
		MerchantName:    merchantName,
		Pending:         plaidTransaction.Pending,
		TransactionID:   plaidTransaction.TransactionId,
		Date:            plaidTransaction.Date,
		UserID:          account.UserID,
	}
}
