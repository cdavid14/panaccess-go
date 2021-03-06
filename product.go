package panaccess

import (
	"encoding/json"
	"errors"
	"net/url"
)

//GetListOfProductsReponse from panaccess
type GetListOfProductsReponse struct {
	Count          int       `json:"count"`
	ProductEntries []Product `json:"productEntries"`
}

//Product class representation from panaccess
type Product struct {
	ID      int    `json:"productId"`
	Name    string `json:"name"`
	Deleted bool   `json:"deleted"`
}

//Get product from panaccess
func (prod *Product) Get(pan *Panaccess, params *url.Values) ([]Product, error) {
	//Everything has a limit
	if (*params).Get("limit") == "" {
		(*params).Add("limit", "1000")
	}
	//Call Function
	resp, err := pan.Call(
		"getListOfProducts",
		params)
	if err != nil {
		return nil, err
	}
	//Retrieve all rows and parse as a slice of Subscriber
	var rows GetListOfProductsReponse
	bodyBytes, err := json.Marshal(resp.Answer)
	err = json.Unmarshal(bodyBytes, &rows)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.ErrorMessage)
	}
	return rows.ProductEntries, nil
}

//GetWithFilter product from panaccess
func (prod *Product) GetWithFilter(pan *Panaccess, params *url.Values, groupOp string, filters []Rule) ([]Product, error) {
	//Everything has a limit
	if (*params).Get("limit") == "" {
		(*params).Add("limit", "1000")
	}
	//Call Function
	resp, err := pan.CallWithFilters(
		"getListOfProducts",
		params,
		groupOp,
		filters)
	if err != nil {
		return nil, err
	}
	//Retrieve all rows and parse as a slice of Subscriber
	var rows GetListOfProductsReponse
	bodyBytes, err := json.Marshal(resp.Answer)
	err = json.Unmarshal(bodyBytes, &rows)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.ErrorMessage)
	}
	return rows.ProductEntries, nil
}
