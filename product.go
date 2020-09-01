package panaccess

import (
	"encoding/json"
	"io/ioutil"
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
	//Decode Response to Struct
	ret := APIResponse{}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, &ret)
	if err != nil {
		return nil, err
	}
	//Retrieve all rows and parse as a slice of Subscriber
	var rows GetListOfProductsReponse
	bodyBytes, err = json.Marshal(ret.Answer)
	err = json.Unmarshal(bodyBytes, &rows)
	if err != nil {
		return nil, err
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
	//Decode Response to Struct
	ret := APIResponse{}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, &ret)
	if err != nil {
		return nil, err
	}
	//Retrieve all rows and parse as a slice of Subscriber
	var rows GetListOfProductsReponse
	bodyBytes, err = json.Marshal(ret.Answer)
	err = json.Unmarshal(bodyBytes, &rows)
	if err != nil {
		return nil, err
	}
	return rows.ProductEntries, nil
}
