package service

import "net/http"

// @sequence guard nil user

// @sequence call
// @func doSomething
// @param order.Total

// @sequence put
// @model Order.Update
// @param order.ID

// @sequence response json
// @var user
// @var order
// @var total
func Test2VarFlow(w http.ResponseWriter, r *http.Request) {}
