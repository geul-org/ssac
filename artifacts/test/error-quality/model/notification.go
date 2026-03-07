package model

type Notification interface {
	Execute(account interface{}) error
}
