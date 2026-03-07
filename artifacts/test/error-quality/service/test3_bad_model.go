package service

import "net/http"

// @sequence get
// @model User
// @result user User

// @sequence get
// @model .FindByID
// @result x X

// @sequence get
// @model Project.FindByID
// @param ProjectID request
// @result project Project

// @sequence post
// @model Session.Create
// @param project.ID
// @result session Session

// @sequence response json
// @var session
func Test3BadModel(w http.ResponseWriter, r *http.Request) {}
