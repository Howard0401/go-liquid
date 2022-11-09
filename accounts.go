package liquid

import (
	"fmt"
	"math/big"
	"strings"
)

type Account struct {
	ID                       int       `json:"id"`
	Currency                 string    `json:"currency"`
	CurrencySymbol           string    `json:"currency_symbol"`
	Balance                  float64   `json:"balance,string"`
	PusherChannel            string    `json:"pusher_channel"`
	LowestOfferInterestRate  float64   `json:"lowest_offer_interest_rate,string"`
	HighestOfferInterestRate float64   `json:"highest_offer_interest_rate,string"`
	ExchangeRate             big.Float `json:"exchange_rate,string"`
	CurrencyType             string    `json:"currency_type"`
}

func (c *Client) GetFiatAccounts() ([]Account, error) {
	res, err := c.sendRequest("GET", "/fiat_accounts", nil, nil)
	if err != nil {
		return nil, err
	}

	var accounts []Account
	if err := decode(res, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (c *Client) CreateFiatAccount(currency string) (Account, error) {
	var account Account
	body := fmt.Sprintf(`{"currency":"%s"}`, currency)
	res, err := c.sendRequest("POST", "/fiat_accounts", strings.NewReader(body), nil)
	if err != nil {
		return account, err
	}

	if err := decode(res, &account); err != nil {
		return account, err
	}

	return account, nil
}

type CryptoAccount struct {
	ID                       int     `json:"id"`
	Balance                  string  `json:"balance"`
	Address                  string  `json:"address"`
	Currency                 string  `json:"currency"`
	CurrencySymbol           string  `json:"currency_symbol"`
	PusherChannel            string  `json:"pusher_channel"`
	MinimumWithdraw          float64 `json:"minimum_withdraw"`
	LowestOfferInterestRate  string  `json:"lowest_offer_interest_rate"`
	HighestOfferInterestRate string  `json:"highest_offer_interest_rate"`
	CurrencyType             string  `json:"currency_type"`
}

func (c *Client) GetCryptoAccounts() ([]CryptoAccount, error) {
	res, err := c.sendRequest("GET", "/crypto_accounts", nil, nil)
	if err != nil {
		return nil, err
	}

	var accounts []CryptoAccount
	if err := decode(res, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

type AccountBalance struct {
	Currency string `json:"currency"`
	Balance  string `json:"balance"`
}

type AccountBalanceV2 struct {
	CryptoAccounts []struct {
		ID                       int         `json:"id,omitempty"`
		Currency                 string      `json:"currency,omitempty"`
		Balance                  string      `json:"balance,omitempty"`
		ReservedBalance          string      `json:"reserved_balance,omitempty"`
		PusherChannel            string      `json:"pusher_channel,omitempty"`
		LowestOfferInterestRate  interface{} `json:"lowest_offer_interest_rate,omitempty"`
		HighestOfferInterestRate interface{} `json:"highest_offer_interest_rate,omitempty"`
		Address                  interface{} `json:"address,omitempty"`
		CurrencySymbol           string      `json:"currency_symbol,omitempty"`
		MinimumWithdraw          interface{} `json:"minimum_withdraw,omitempty"`
		CurrencyType             string      `json:"currency_type,omitempty"`
	} `json:"crypto_accounts,omitempty"`
	FiatAccounts []struct {
		ID                       int         `json:"id,omitempty"`
		Currency                 string      `json:"currency,omitempty"`
		Balance                  string      `json:"balance,omitempty"`
		ReservedBalance          string      `json:"reserved_balance,omitempty"`
		PusherChannel            string      `json:"pusher_channel,omitempty"`
		LowestOfferInterestRate  interface{} `json:"lowest_offer_interest_rate,omitempty"`
		HighestOfferInterestRate interface{} `json:"highest_offer_interest_rate,omitempty"`
		CurrencySymbol           string      `json:"currency_symbol,omitempty"`
		SendToBtcAddress         interface{} `json:"send_to_btc_address,omitempty"`
		ExchangeRate             string      `json:"exchange_rate,omitempty"`
		CurrencyType             string      `json:"currency_type,omitempty"`
	} `json:"fiat_accounts,omitempty"`
}

func (c *Client) GetAllAccountBalances() ([]AccountBalance, error) {
	res, err := c.sendRequest("GET", "/accounts/balance", nil, nil)
	if err != nil {
		return nil, err
	}

	var balances []AccountBalance
	if err := decode(res, &balances); err != nil {
		return nil, err
	}

	return balances, nil
}

func (c *Client) GetAllAccountBalancesV2() (AccountBalanceV2, error) {
	res, err := c.sendRequest("GET", "/accounts/balance", nil, nil)
	if err != nil {
		return AccountBalanceV2{}, err
	}

	var balances AccountBalanceV2
	if err := decode(res, &balances); err != nil {
		return balances, err
	}

	return balances, nil
}
