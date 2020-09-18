package panaccess

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

//Smartcard class representation from panaccess
type Smartcard struct {
	Alias                   string   `json:"alias"`
	Blacklisted             bool     `json:"blacklisted"`
	CamlibVersion           string   `json:"camlibVersion"`
	CasIDs                  string   `json:"casIds"`
	ConfigID                string   `json:"configId"`
	ConfigProtected         bool     `json:"configProtected"`
	Defect                  bool     `json:"defect"`
	Disabled                bool     `json:"disabled"`
	FirmwareVersion         string   `json:"firmwareVersion"`
	FirstName               string   `json:"firstName"`
	HCID                    string   `json:"hcId"`
	LastActivation          string   `json:"lastActivation"`
	LastName                string   `json:"lastName"`
	LastServiceListDownload string   `json:"lastServiceListDownload"`
	MAC                     string   `json:"mac"`
	MasterSN                string   `json:"masterSn"`
	PackageNames            []string `json:"packageNames"`
	Packages                []int    `json:"packages"`
	PairedBox               string   `json:"pairedBox"`
	PIN                     string   `json:"pin"`
	Products                []string `json:"products"`
	RegionID                int      `json:"regionId"`
	RegionName              string   `json:"regionName"`
	SN                      string   `json:"sn"`
	STBChipset              string   `json:"stbChipset"`
	STBModel                string   `json:"stbModel"`
	STBVendor               string   `json:"stbVendor"`
	SubscriberCode          string   `json:"subscriberCode"`
}

//Smartcards array of smartcard
type Smartcards []Smartcard

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
	rows := GetListOfSmartcardsResponse{}
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
	jsonBody, err := json.Marshal(resp.Answer)
	err = json.Unmarshal(jsonBody, &rows)
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
	var rows Smartcards
	bodyBytes, err := json.Marshal(resp.Answer)
	fmt.Println(string(bodyBytes))
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
