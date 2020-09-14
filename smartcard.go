package panaccess

import (
	"encoding/json"
	"errors"
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
	SN          string   `json:"sn,omitempty"`
	PIN         string   `json:"pin,omitempty"`
	Checksum    string   `json:"checksum,omitempty"`
	HCID        string   `json:"hcId,omitempty"`
	Disabled    bool     `json:"disabled,omitempty"`
	Defect      bool     `json:"defect,omitempty"`
	Blacklisted bool     `json:"blacklisted,omitempty"`
	Products    []string `json:"products,omitempty"`
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
	//Retrieve all rows and parse as a slice of Subscriber
	var rows GetListOfSmartcardsResponse
	bodyBytes, err := json.Marshal(resp.Answer)
	err = json.Unmarshal(bodyBytes, &rows)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.ErrorMessage)
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
	//Retrieve all rows and parse as a slice of Subscriber
	var rows GetListOfSmartcardsResponse
	bodyBytes, err := json.Marshal(resp.Answer)
	err = json.Unmarshal(bodyBytes, &rows)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.ErrorMessage)
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
	if !resp.Success {
		return nil, errors.New(resp.ErrorMessage)
	}
	var rows []Smartcard
	bodyBytes, err := json.Marshal(resp.Answer)
	err = json.Unmarshal(bodyBytes, &rows)
	return rows, nil
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
	if !resp.Success {
		return errors.New(resp.ErrorMessage)
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
	if !resp.Success {
		return errors.New(resp.ErrorMessage)
	}
	return nil
}
