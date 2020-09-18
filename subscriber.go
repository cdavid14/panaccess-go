package panaccess

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

//GetListOfSubscribersResponse from panaccess
type GetListOfSubscribersResponse struct {
	Count             int          `json:"count"`
	SubscriberEntries []Subscriber `json:"extendedSubscriberEntries"`
}

//Subscriber class representation from panaccess
type Subscriber struct {
	SubscriberCode string      `json:"subscriberCode"`
	RegionID       int         `json:"regionId"`
	FirstName      string      `json:"firstName"`
	LastName       string      `json:"lastName"`
	CountryCode    string      `json:"countryCode"`
	CAF            interface{} `json:"caf"`
	Smartcards     []string    `json:"smartcards"`
	Comment        string      `json:"comment"`
	Supervisor     string      `json:"supervisor"`
	TechNotes      string      `json:"technicalNotes"`
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
	//Retrieve all rows and parse as a slice of Subscriber
	var rows GetListOfSubscribersResponse
	bodyBytes, err := json.Marshal(resp.Answer)
	err = json.Unmarshal(bodyBytes, &rows)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.ErrorMessage)
	}
	return rows.SubscriberEntries, nil
}

//Delete a subscriber
func (sub *Subscriber) Delete(pan *Panaccess) error {
	params := url.Values{}
	params.Add("code", sub.SubscriberCode)
	//Call Function
	resp, err := pan.Call(
		"deleteSubscriber",
		&params)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.ErrorMessage)
	}
	return nil
}

//GetWithFilters a list of subscribers with specific filters
func (sub *Subscriber) GetWithFilters(pan *Panaccess, params *url.Values, groupOp string, filters []Rule) ([]Subscriber, error) {
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
	//Retrieve all rows and parse as a slice of Subscriber
	var rows GetListOfSubscribersResponse
	bodyBytes, err := json.Marshal(resp.Answer)
	if err != nil {
		return nil, err
	}
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
			OP:    "eq",
			Data:  sub.SubscriberCode,
		},
	})
}

//GetSmartcardsWithFilter of Subscriber
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
	var ordersResponse []Order
	bodyBytes, err := json.MarshalIndent(resp.Answer, "", "  ")
	if err != nil {
		return nil, err
	}
	fmt.Println(string(bodyBytes))
	err = json.Unmarshal(bodyBytes, &ordersResponse)
	if err != nil {
		return nil, err
	}
	return ordersResponse, nil
}

//LockOrder from subscriber at panaccess
func (sub *Subscriber) LockOrder(pan *Panaccess, order *Order) error {
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
	if !resp.Success {
		return errors.New(resp.ErrorMessage)
	}
	return nil
}

//UnlockOrder from subscriber at panaccess
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
	if !resp.Success {
		return errors.New(resp.ErrorMessage)
	}
	return nil
}
