package main

import (
	"backend/routes/backend/database"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
	"github.com/plaid/plaid-go/v21/plaid"
)

var plaidClient *plaid.APIClient

func main() {
	database.Connect()
	database.Migrate()

	r := mux.NewRouter()
	r.HandleFunc("/api/create-user", createUser)
	r.HandleFunc("/api/get-net-worth", getNetWorth).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/create-link-token", createLinkToken).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/create-item", createItem).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/refresh-user-items", refreshUserItems).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/get-transactions", getTransactions).Methods("GET", "OPTIONS")

	corsOpts := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // change me!
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
	)

	httpHandler := corsOpts(r)

	addr := ":8080"
	log.Printf("Starting server on 	%s", addr)
	if err := http.ListenAndServe(addr, httpHandler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

type linkTokenRequest struct {
	UserId string `json:"userId"`
}

type linkTokenResponse struct {
	LinkToken string `json:"linkToken"`
}

func createLinkToken(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user linkTokenRequest
	decoder.Decode(&user)

	var result database.User
	err := database.Instance.Where("id = ?", user.UserId).First(&result).Error
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	plaidUser := plaid.LinkTokenCreateRequestUser{
		ClientUserId: result.Id,
	}

	request := plaid.NewLinkTokenCreateRequest(
		"Peppermint",
		"en",
		[]plaid.CountryCode{plaid.COUNTRYCODE_US},
		plaidUser,
	)

	request.SetProducts([]plaid.Products{plaid.PRODUCTS_TRANSACTIONS})
	// request.SetWebhook("https://sample-web-hook.com")
	// request.SetRedirectUri("https://domainname.com/oauth-page.html")

	ctx := context.Background()
	plaidClient := createPlaidClient()
	linkTokenCreateResp, _, err := plaidClient.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		panic(err)
	}
	linkToken := linkTokenCreateResp.GetLinkToken()

	resp := linkTokenResponse{
		LinkToken: linkToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

type plaidCreateItemRequest struct {
	UserId      string
	PublicToken string
	Metadata    plaidMetadata
}

type plaidMetadata struct {
	Institution   institutionData `json:"institution"`
	Accounts      []accountData   `json:"accounts"`
	LinkSessionID string          `json:"link_session_id"`
}

type institutionData struct {
	Name          string `json:"name"`
	InstitutionID string `json:"institution_id"`
}

type accountData struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Mask               string `json:"mask"`
	Type               string `json:"type"`
	Subtype            string `json:"subtype"`
	VerificationStatus string `json:"verification_status,omitempty"`
}

// Creates new Plaid item, given access_token & metadata
// see https://plaid.com/docs/link/web/#onsuccess
func createItem(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var itemRequest plaidCreateItemRequest
	decoder.Decode(&itemRequest)

	plaidClient := createPlaidClient()
	ctx := context.Background()

	exchangePublicTokenReq := plaid.NewItemPublicTokenExchangeRequest(itemRequest.PublicToken)
	exchangePublicTokenResp, _, err := plaidClient.PlaidApi.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(
		*exchangePublicTokenReq,
	).Execute()

	if err != nil {
		panic(err)
	}

	accessToken := exchangePublicTokenResp.GetAccessToken()

	existingAccounts := database.GetAccountsForUser(itemRequest.UserId)

	for _, account := range itemRequest.Metadata.Accounts {
		for _, existingAccount := range existingAccounts {
			if existingAccount.UserID == itemRequest.UserId && account.Subtype == existingAccount.AccountSubtype && itemRequest.Metadata.Institution.InstitutionID == existingAccount.InstitutionID {
				database.Instance.Delete(&existingAccount)
				break
			}
		}

		newAccount := database.Account{
			UserID:             itemRequest.UserId,
			AccountType:        account.Type,
			AccountSubtype:     account.Subtype,
			InstitutionID:      itemRequest.Metadata.Institution.InstitutionID,
			Balance:            0,
			CreditLimit:        0,
			PlaidAccessToken:   accessToken,
			AccountID:          account.ID,
			TransactionsCursor: "",
		}
		database.Instance.Create(&newAccount)
	}
}

type plaidRefreshItemRequest struct {
	UserId string `json:"userId"`
}

// Refresh plaid item transactions associated with user given user id
func refreshUserItems(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var itemRefreshRequest plaidRefreshItemRequest
	decoder.Decode(&itemRefreshRequest)

	plaidClient := createPlaidClient()
	ctx := context.Background()

	accounts := database.GetAccountsForUser(itemRefreshRequest.UserId)

	for _, account := range accounts {
		cursor := database.GetLatestCursorOrNil(account.AccountID)

		var ptr bool = true
		var added []plaid.Transaction
		var modified []plaid.Transaction
		var removed []plaid.RemovedTransaction // Removed transaction ids
		options := plaid.TransactionsSyncRequestOptions{
			IncludePersonalFinanceCategory: &ptr,
		}
		hasMore := true

		// Get transactions from Plaid
		for hasMore {
			request := plaid.NewTransactionsSyncRequest(account.PlaidAccessToken)
			request.SetOptions(options)
			if cursor != "" {
				request.SetCursor(cursor)
			}
			resp, _, err := plaidClient.PlaidApi.TransactionsSync(
				ctx,
			).TransactionsSyncRequest(*request).Execute()

			if err != nil {
				panic(err)
			}

			added = append(added, resp.GetAdded()...)
			modified = append(modified, resp.GetModified()...)
			removed = append(removed, resp.GetRemoved()...)

			hasMore = resp.GetHasMore()
			cursor = resp.GetNextCursor()

			// save cursor to db
			account.TransactionsCursor = cursor
			database.Instance.Save(&account)
		}

		// Update db with transactions
		for _, txn := range added {
			dbTxn := database.ConvertPlaidTransaction(txn, account.AccountID)
			database.Instance.Create(&dbTxn)
		}

		for _, txn := range modified {
			dbTxn := database.ConvertPlaidTransaction(txn, account.AccountID)
			existingTxn := database.GetTransactionByID(txn.TransactionId)
			dbTxn.TransactionID = existingTxn.TransactionID
			database.Instance.Save(&dbTxn)
		}

		for _, txn := range removed {
			database.Instance.Where("transaction_id = ?", txn.TransactionId).Delete(&database.Transaction{})
		}

		// Update account balance
		accountsGetRequest := plaid.NewAccountsGetRequest(account.PlaidAccessToken)

		accountsGetResp, _, err := plaidClient.PlaidApi.AccountsGet(ctx).AccountsGetRequest(
			*accountsGetRequest,
		).Execute()

		if err != nil {
			panic(err)
		}

		for _, accountResp := range accountsGetResp.GetAccounts() {
			if account.AccountID == accountResp.AccountId {
				account.Balance = *accountResp.GetBalances().Current.Get()
				if accountResp.Type == "credit" {
					account.CreditLimit = *accountResp.GetBalances().Limit.Get()
				}
				database.Instance.Save(&account)
			}
		}
	}
}

type netWorthRequest struct {
	userId string
}

func getNetWorth(w http.ResponseWriter, r *http.Request) {
	var netWorth float64 = 0

	userId := r.URL.Query().Get("userId")

	accounts := database.GetAccountsForUser(userId)

	for _, account := range accounts {
		if account.AccountType == "credit" {
			netWorth -= account.Balance
		} else {
			netWorth += account.Balance
		}
	}

	netWorthResponse := database.NetWorthResponse{NetWorth: netWorth, Datetime: time.Now()}

	resp, err := json.Marshal(netWorthResponse)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result := database.Instance.Create(&netWorthResponse); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

type transactionsResponse struct {
	Transactions []database.Transaction
}

func getTransactions(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")

	transactions := database.GetTransactionsForUser(userId)

	resp, err := json.Marshal(transactions)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

// Creates new user & adds to db
func createUser(w http.ResponseWriter, r *http.Request) {
	var user database.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// id := uuid.New().String()
	user.Id = "ladur"

	if result := database.Instance.Create(&user); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Success")
}
