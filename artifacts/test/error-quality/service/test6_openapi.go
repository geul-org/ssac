package service

import "net/http"

// @sequence get
// @model User.FindByEmail
// @param Phone request
// @param Address request
// @result user User

// @sequence response json
// @var user
// @var token
func Login(w http.ResponseWriter, r *http.Request) {}
