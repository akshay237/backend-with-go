package api

import (
	util "github.com/akshay237/backend-with-go/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if curr, isok := fl.Field().Interface().(string); isok {
		return util.IsSupportedCurrency(curr)
	}
	return false
}
