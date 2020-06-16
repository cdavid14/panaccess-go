package panaccess

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
)

type GetListOfSubscribersResponse struct {
	Count             int          `json:"count"`
	SubscriberEntries []Subscriber `json:"extendedSubscriberEntries"`
}

//Subscriber representation
type Subscriber struct {
	SubscriberCode string      `json:"subscriberCode"`
	RegionID       int         `json:"regionId"`
	FirstName      string      `json:"firstName"`
	LastName       string      `json:"lastName"`
	CountryCode    string      `json:"countryCode"`
	CAF            interface{} `json:"caf"`
	Smartcards     []string    `json:"smartcards"`
	Comment        string      `json:"comment"`
	LastExpiryTime string      `json:"lastExpiryTime"`
	CreatedAt      string      `json:"created"`
}

//Get a list of subscribers
func (sub *Subscriber) Get(pan *Panaccess, params *url.Values) ([]Subscriber, error) {
	//Everything has a limit
	if (*params).Get("limit") == "" {
		(*params).Add("limit", "1000")
	}
	//Call Function
	resp, err := pan.Call(
		"getListOfExtendedSubscribers",
		params)
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
	var rows GetListOfSubscribersResponse
	bodyBytes, err = json.Marshal(ret.Answer)
	err = json.Unmarshal(bodyBytes, &rows)
	if err != nil {
		return nil, err
	}
	return rows.SubscriberEntries, nil
}

//GetWithFilters a list of subscribers with specific filters
func (sub *Subscriber) GetWithFilter(pan *Panaccess, params *url.Values, groupOp string, filters []Rule) ([]Subscriber, error) {
	//Everything has a limit
	if params.Get("limit") == "" {
		params.Add("limit", "1000")
	}
	//Call Function
	resp, err := pan.CallWithFilters(
		"getListOfExtendedSubscribers",
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
	var rows GetListOfSubscribersResponse
	bodyBytes, err = json.Marshal(ret.Answer)
	err = json.Unmarshal(bodyBytes, &rows)
	if err != nil {
		return nil, err
	}
	return rows.SubscriberEntries, nil
}

//GetSmartcards of Subscriber
func (sub *Subscriber) GetSmartcards(pan *Panaccess) ([]Smartcard, error) {
	cards := Smartcard{}
	return cards.GetWithFilter(pan, &url.Values{}, "AND", []Rule{
		{
			Field: "subscriberCode",
			OP:    "cn",
			Data:  sub.SubscriberCode,
		},
	})
}

//GetSmartcards of Subscriber
func (sub *Subscriber) GetSmartcardsWithFilter(pan *Panaccess, filter []Rule) ([]Smartcard, error) {
	cards := Smartcard{}
	filter = append(filter, Rule{
		Field: "subscriberCode",
		OP:    "cn",
		Data:  sub.SubscriberCode,
	})
	return cards.GetWithFilter(pan, &url.Values{}, "AND", filter)
}

//GetOrders of Subscriber
func (sub *Subscriber) GetOrders(pan *Panaccess, params *url.Values) ([]Order, error) {
	if (*params).Get("limit") == "" {
		(*params).Set("limit", "1000")
	}
	(*params).Set("subscriberCode", sub.SubscriberCode)
	//Call Function
	resp, err := pan.Call(
		"getOrdersOfSubscriber",
		params,
	)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	test := GetSmartcardOrdersResponse{}
	err = json.Unmarshal(bodyBytes, &test)
	if err != nil {
		return nil, err
	}
	return test.Answer, nil
}

func (sub *Subscriber) LockOrder(pan *Panaccess, order *Order) error {
	loggedIn, _ := pan.Loggedin()
	if !loggedIn {
		err := pan.Login()
		if err != nil {
			return err
		}
	}
	//Verify Fields
	if order == nil {
		return errors.New("Please fill all required fields")
	}
	params := url.Values{}
	params.Add("orderId", fmt.Sprint(order.ID))
	params.Add("subscriberCode", sub.SubscriberCode)
	//Send data to make new subscriber
	resp, err := pan.Call(
		"disableOrderOfSubscriber",
		&params)
	if err != nil {
		return err
	}
	ret := ApiResponse{}
	json.NewDecoder(resp.Body).Decode(&ret)
	if !ret.Success {
		return errors.New(ret.ErrorMessage)
	}
	return nil
}

func (sub *Subscriber) UnlockOrder(pan *Panaccess, order *Order) error {
	loggedIn, _ := pan.Loggedin()
	if !loggedIn {
		err := pan.Login()
		if err != nil {
			return err
		}
	}
	//Verify Fields
	if order == nil {
		return errors.New("Please fill all required fields")
	}
	params := url.Values{}
	params.Add("orderId", fmt.Sprint(order.ID))
	params.Add("subscriberCode", sub.SubscriberCode)
	params.Add("until", "")
	//Send data to make new subscriber
	resp, err := pan.Call(
		"enableOrderOfSubscriber",
		&params)
	if err != nil {
		return err
	}
	ret := ApiResponse{}
	json.NewDecoder(resp.Body).Decode(&ret)
	if !ret.Success {
		return errors.New(ret.ErrorMessage)
	}
	return nil
}
