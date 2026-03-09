package service

import "net/http"

// @sequence get
// @model User.FindByEmail
// @param Email request
// @result user User

// @sequence response json
// @var user
func Login(w http.ResponseWriter, r *http.Request) {}
