package panaccess

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
)

//GetListOfSmartcardsResponse from panaccess
type GetListOfSmartcardsResponse struct {
	Count            int         `json:"count"`
	SmartcardEntries []Smartcard `json:"smartcardEntries"`
}

//GetUnusedSmartcardsResponse from panaccess
type GetUnusedSmartcardsResponse struct {
	Success bool        `json:"success"`
	Answer  []Smartcard `json:"answer"`
}

//Smartcard class representation from panaccess
type Smartcard struct {
	SN          string   `json:"sn"`
	PIN         string   `json:"pin"`
	Checksum    string   `json:"checksum,omitempty"`
	HCID        string   `json:"hcId"`
	Disabled    bool     `json:"disabled"`
	Defect      bool     `json:"defect"`
	Blacklisted bool     `json:"blacklisted"`
	Products    []string `json:"products"`
}

//GetSmartcardOrdersResponse from panaccess
type GetSmartcardOrdersResponse struct {
	Success bool    `json:"success"`
	Answer  []Order `json:"answer"`
}

//Get smartcard from panaccess
func (card *Smartcard) Get(pan *Panaccess, params *url.Values) ([]Smartcard, error) {
	//Everything has a limit
	if (*params).Get("limit") == "" {
		(*params).Add("limit", "1000")
	}
	//Call Function
	resp, err := pan.Call(
		"getListOfSmartcards",
		params,
	)
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
	var rows GetListOfSmartcardsResponse
	bodyBytes, err = json.Marshal(ret.Answer)
	err = json.Unmarshal(bodyBytes, &rows)
	if err != nil {
		return nil, err
	}
	return rows.SmartcardEntries, nil
}

//GetWithFilter smartcard from panaccess
func (card *Smartcard) GetWithFilter(pan *Panaccess, params *url.Values, groupOp string, filters []Rule) ([]Smartcard, error) {
	//Everything has a limit
	if (*params).Get("limit") == "" {
		(*params).Add("limit", "1000")
	}
	//Call Function
	resp, err := pan.CallWithFilters(
		"getListOfSmartcards",
		params,
		groupOp,
		filters,
	)
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
	var rows GetListOfSmartcardsResponse
	bodyBytes, err = json.Marshal(ret.Answer)
	err = json.Unmarshal(bodyBytes, &rows)
	if err != nil {
		return nil, err
	}
	return rows.SmartcardEntries, nil
}

//GetUnused smartcard from panaccess
func (card *Smartcard) GetUnused(pan *Panaccess, params *url.Values) ([]Smartcard, error) {
	//Everything has a limit
	if (*params).Get("limit") == "" {
		(*params).Add("limit", "1000")
	}
	//Call Function
	resp, err := pan.Call(
		"getUnusedSmartcards",
		params,
	)
	if err != nil {
		return nil, err
	}
	//Decode Response to Struct
	ret := GetUnusedSmartcardsResponse{}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, &ret)
	if err != nil {
		return nil, err
	}
	return ret.Answer, nil
}

//Unlock smartcard from panaccess
func (card *Smartcard) Unlock(pan *Panaccess) error {
	//Params
	params := url.Values{}
	params.Add("smartcardId", card.SN)
	//Call Function
	resp, err := pan.Call(
		"enableSmartcard",
		&params,
	)
	if err != nil {
		return err
	}
	//Decode Response to Struct
	ret := APIResponse{}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyBytes, &ret)
	if err != nil {
		return err
	}
	if !ret.Success {
		return errors.New(ret.ErrorMessage)
	}
	return nil
}

//Lock smartcard from panaccess
func (card *Smartcard) Lock(pan *Panaccess) error {
	//Params
	params := url.Values{}
	params.Add("smartcardId", card.SN)
	//Call Function
	resp, err := pan.Call(
		"disableSmartcard",
		&params,
	)
	if err != nil {
		return err
	}
	//Decode Response to Struct
	ret := APIResponse{}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyBytes, &ret)
	if err != nil {
		return err
	}
	if !ret.Success {
		return errors.New(ret.ErrorMessage)
	}
	return nil
}
