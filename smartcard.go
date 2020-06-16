package panaccess

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
)

type GetListOfSmartcardsResponse struct {
	Count            int         `json:"count"`
	SmartcardEntries []Smartcard `json:"smartcardEntries"`
}

type GetUnusedSmartcardsResponse struct {
	Success bool        `json:"success"`
	Answer  []Smartcard `json:"answer"`
}

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

type GetSmartcardOrdersResponse struct {
	Success bool    `json:"success"`
	Answer  []Order `json:"answer"`
}

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
	ret := ApiResponse{}
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
	ret := ApiResponse{}
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
	ret := ApiResponse{}
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
	ret := ApiResponse{}
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
