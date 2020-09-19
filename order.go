package panaccess

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

//Order class representation from panaccess
type Order struct {
	ActivationTime string   `json:"activationTime"`
	Alias          string   `json:"alias"`
	SubscriberCode string   `json:"code"`
	Created        string   `json:"created"`
	ExpiryTime     string   `json:"expiryTime"`
	FirstName      string   `json:"firstName"`
	LastName       string   `json:"lastName"`
	Modified       string   `json:"modified"`
	ID             int      `json:"orderId"`
	OrderTime      string   `json:"orderTime"`
	ProductID      int      `json:"productId"`
	ProductName    string   `json:"productName"`
	ScDefect       bool     `json:"scDefect"`
	ScDisabled     bool     `json:"scDisabled"`
	Smartcards     []string `json:"smartcards"`
	SN             string   `json:"sn"`
}

//GetOrdersFilterResponse from panaccess
type GetOrdersFilterResponse struct {
	Count        int     `json:"count"`
	OrderEntries []Order `json:"orderEntries"`
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
	bodyBytes, err := json.MarshalIndent(resp.Answer, "", "  ")
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
	fmt.Printf("Cards: %v\n", cards)
	for _, card := range cards {
		fmt.Printf("Products: %v | Name: %v\n", card.Products, prods[0].Name)
		fmt.Printf("Len1: %v | Len2: %v\n", sort.SearchStrings(card.Products, prods[0].Name), len(card.Products))
		found := false
		for _, v := range card.Products {
			if strings.Compare(v, prods[0].Name) == 0 {
				found = true
			}
		}
		if !found {
			(*params).Add("smartcards[]", card.SN)
		}
	}
	//Send data to make new subscriber
	fmt.Println(*params)
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
		"terminateOrderOfSubscriber",
		&params)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.ErrorMessage)
	}
	return nil
}
