package Service

type Service interface {
	QueryById() (interface{}, error)
	Query() (interface{}, error)
	UpdateService() error
	Add() (int error)
}
