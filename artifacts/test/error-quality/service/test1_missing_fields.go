package service

import "net/http"

// @sequence authorize

// @sequence get
// @param ID request

// @sequence guard nil

// @sequence password
// @param hash

// @sequence call

// @sequence response json
// @var result
func Test1MissingFields(w http.ResponseWriter, r *http.Request) {}
