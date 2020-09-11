package panaccess

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
)

//GetOrdersFilterResponse from panaccess
type GetOrdersFilterResponse struct {
	Count        int     `json:"count"`
	OrderEntries []Order `json:"orderEntries"`
}

//Order class representation from panaccess
type Order struct {
	ID                 int      `json:"orderId"`
	OrderTime          string   `json:"orderTime,omitempty"`
	ProductID          int      `json:"productId,omitempty"`
	ProductName        string   `json:"productName,omitempty"`
	ActivationTime     string   `json:"activationTime"`
	SubscriberCode     int      `json:"code,omitempty"`
	ExpiryTime         string   `json:"expiryTime"`
	Smartcards         []string `json:"smartcards"`
	Disabled           bool     `json:"disabled"`
	DisabledBySystem   bool     `json:"disabledBySystem"`
	DisabledByOperator bool     `json:"disabledByOperator"`
}

//Get order from panaccess
func (order *Order) Get(pan *Panaccess, params *url.Values) ([]Order, error) {
	//Everything has a limit
	if (*params).Get("limit") == "" {
		(*params).Add("limit", "1000")
	}
	//Call Function
	resp, err := pan.Call(
		"getListOfOrders",
		params)
	if err != nil {
		return nil, err
	}
	//Retrieve all rows and parse as a slice of Subscriber
	var rows GetOrdersFilterResponse
	bodyBytes, err := json.Marshal(resp.Answer)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, &rows)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.ErrorMessage)
	}
	return rows.OrderEntries, nil
}

//GetWithFilters order from panaccess
func (order *Order) GetWithFilters(pan *Panaccess, params *url.Values, groupOp string, filters []Rule) ([]Order, error) {
	//Everything has a limit
	if (*params).Get("limit") == "" {
		(*params).Add("limit", "1000")
	}
	//Call Function
	resp, err := pan.CallWithFilters(
		"getListOfOrders",
		params,
		groupOp,
		filters,
	)
	if err != nil {
		return nil, err
	}
	//Retrieve all rows and parse as a slice of Subscriber
	var rows GetOrdersFilterResponse
	bodyBytes, err := json.Marshal(resp.Answer)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, &rows)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.ErrorMessage)
	}
	return rows.OrderEntries, nil
}

//AddToSubscriber a order from panaccess
func (order *Order) AddToSubscriber(pan *Panaccess, params *url.Values) error {
	//Verify Fields
	if params.Get("productId") == "" || params.Get("subscriberCode") == "" || params.Get("activationTime") == "" || params.Get("expiryTime") == "" {
		return errors.New("Please fill all required fields")
	}
	(*params).Set("onlySpecifiedSmartcards", "true")
	//Verify if user exists
	if params.Get("subscriberCode") != "" {
		resp, err := pan.Call(
			"subscriberExists",
			params,
		)
		if err != nil {
			return err
		}
		if !resp.Success {
			return errors.New(resp.ErrorMessage)
		}
	}
	//Get Subscriber smartcards
	sub := Subscriber{
		SubscriberCode: params.Get("subscriberCode"),
	}
	cards, err := sub.GetSmartcards(pan)
	if err != nil {
		return err
	}
	//Get Product Name
	prod := Product{}
	prods, err := prod.GetWithFilter(pan, &url.Values{}, "AND", []Rule{
		{
			Field: "productId",
			OP:    "eq",
			Data:  params.Get("productId"),
		},
	})
	if err != nil {
		return err
	}
	if len(prods) == 0 {
		return errors.New("ProductId not found")
	}
	//Add card to product if hasn't
	for _, card := range cards {
		if sort.SearchStrings(card.Products, prods[0].Name) == len(card.Products) {
			params.Add("smartcards", card.SN)
		}
	}
	//Send data to make new subscriber
	resp, err := pan.Call(
		"addFlexibleOrderToSubscriber",
		params)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.ErrorMessage)
	}
	return nil
}

//RemoveFromSubscriber order from panaccess
func (order *Order) RemoveFromSubscriber(pan *Panaccess, sub *Subscriber) error {
	params := url.Values{}
	params.Add("orderId", fmt.Sprint(order.ID))
	params.Add("subscriberCode", sub.SubscriberCode)
	resp, err := pan.Call(
		"cancelOrderOfSubscriber",
		&params)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.ErrorMessage)
	}
	return nil
}
