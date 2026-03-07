package service

import "net/http"

// @sequence get
// @model Account.FindByID
// @param AccountID request
// @result account Account

// @sequence call
// @component emailer
// @param account

// @sequence call
// @func processPayment
// @param account
// @result payment Payment

// @sequence response json
// @var account
// @var payment
func Test5External(w http.ResponseWriter, r *http.Request) {}
