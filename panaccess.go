package panaccess

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

//Panaccess credentials to login
type Panaccess struct {
	Servers   []string
	User      string
	Password  string
	Token     string
	SessionID string
	HTTP      *http.Client
}

//Rule of a query
type Rule struct {
	Field string `json:"field"`
	//OP
	// eq = equal|ne=not equal|lt=less than
	// le = less or equal|gt = grater then
	// ge = greater or equal|bw=begins with
	// bn = not begins with|ew = ends with
	// en = not ends with|cn = contains
	// nc = not contains
	OP   string `json:"op"`
	Data string `json:"data"`
}

//Filters of a query
type Filters struct {
	GroupOP string `json:"groupOp"`
	Rules   []Rule `json:"rules"`
}

const (
	salt = "_panaccess" //appended to password
)

//APIResponse marshal JSON output to struct
type APIResponse struct {
	Success          bool        `json:"success"`
	ErrorCode        string      `json:"errorCode,omitempty"`
	ErrorTag         string      `json:"errorTag,omitempty"`
	ErrorMessage     string      `json:"errorMessage,omitempty"`
	ShowErrorMessage bool        `json:"showErrorMessage,omitempty"`
	ShowErrorTag     bool        `json:"showErrorTag,omitempty"`
	Answer           interface{} `json:"answer,omitempty"`
}

//LoggedInResponse from panaccess
type LoggedInResponse struct {
	Success bool `json:"success"`
	Answer  bool `json:"answer"`
}

//KeyValuePair necessary to do an set function
type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//Login in system
func (p *Panaccess) Login() error {
	//Add password salt
	p.Password += salt
	//Encrypt password with MD5
	hasher := md5.New()
	hasher.Write([]byte(p.Password))
	p.Password = hex.EncodeToString(hasher.Sum(nil))
	//Call Panaccess login
	form := url.Values{}
	form.Add("apiToken", p.Token)
	form.Add("username", p.User)
	form.Add("password", p.Password)
	resp, err := p.Call("login", &form)
	if err != nil {
		return err
	}
	//Set SessionID
	p.SessionID = resp.Answer.(string)
	return nil
}

//Loggedin in system
func (p *Panaccess) Loggedin() (bool, error) {
	//Function Call
	var resp *http.Response
	var err error
	serverOk := false
	params := url.Values{}
	params.Add("sessionId", p.SessionID)
	for _, server := range p.Servers {
		resp, err = p.HTTP.PostForm(
			fmt.Sprintf("%s?f=loggedIn&requestMode=function", server),
			params,
		)
		if err == nil {
			serverOk = true
			break
		}
	}
	if !serverOk {
		return false, errors.New("Connection Timeout")
	}

	ret := APIResponse{}
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return false, err
	}
	return ret.Answer.(bool), nil
}

//Logout panaccess system
func (p *Panaccess) Logout() error {
	//Not logged yet
	if p.SessionID == "" {
		return nil
	}
	//Call Logout function
	_, err := p.Call("logout", &url.Values{})
	if err != nil {
		return err
	}
	return nil
}

//Call panaccess function
func (p *Panaccess) Call(funcName string, parameters *url.Values) (*APIResponse, error) {
	//Prevent ADD SessionID when logging in or if hasn't logged yet
	if p.SessionID != "" && funcName != "login" {
		(*parameters).Add("sessionId", p.SessionID)
	}
	//Function Call
	var resp *http.Response
	var err error
	serverOk := false
	for _, server := range p.Servers {
		resp, err = p.HTTP.PostForm(
			fmt.Sprintf("%s?f=%s&requestMode=function", server, funcName),
			(*parameters),
		)
		if err == nil {
			serverOk = true
			break
		}
	}
	if !serverOk {
		return nil, errors.New("Connection Timeout")
	}
	//Decode response to struct
	apiResponse := APIResponse{}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, &apiResponse)
	if err != nil {
		return nil, err
	}
	//Verify if the response has generated errors
	if apiResponse.ErrorCode != "" {
		//Ignore login function
		if funcName != "login" {
			//Check if user is logged-in
			loggedIn, err := p.Loggedin()
			if err != nil {
				return nil, err
			}
			//If not
			if !loggedIn {
				//Do login again
				err = p.Login()
				if err != nil {
					return nil, err
				}
				//Call function again
				return p.Call(funcName, parameters)
			}
			return nil, errors.New(apiResponse.ErrorCode)
		}
	}
	return &apiResponse, nil
}

//CallWithFilters panaccess function
func (p *Panaccess) CallWithFilters(funcName string, parameters *url.Values, filterGroupOP string, filters []Rule) (*APIResponse, error) {
	//Prevent ADD SessionID when logging in or if hasn't logged yet
	if p.SessionID != "" && funcName != "login" {
		(*parameters).Add("sessionId", p.SessionID)
	}
	//Filters generator
	filter := Filters{
		GroupOP: filterGroupOP,
		Rules:   filters,
	}
	filtersText, err := json.Marshal(filter)
	if err != nil {
		fmt.Println(err)
	}
	(*parameters).Add("filters", string(filtersText))
	//Function Call
	var resp *http.Response
	serverOk := false
	for _, server := range p.Servers {
		resp, err = p.HTTP.PostForm(
			fmt.Sprintf("%s?f=%s&requestMode=function", server, funcName),
			(*parameters),
		)
		if err == nil {
			serverOk = true
			break
		}
	}
	if !serverOk {
		return nil, errors.New("Connection Timeout")
	}
	//Decode response to struct
	apiResponse := APIResponse{}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, &apiResponse)
	if err != nil {
		return nil, err
	}
	//Verify if the response has generated errors
	if apiResponse.ErrorCode != "" {
		//Ignore login function
		if funcName != "login" {
			//Check if user is logged-in
			loggedIn, err := p.Loggedin()
			if err != nil {
				return nil, err
			}
			//If not
			if !loggedIn {
				//Do login again
				err = p.Login()
				if err != nil {
					return nil, err
				}
				//Call function again
				return p.Call(funcName, parameters)
			} else {
				return nil, errors.New(apiResponse.ErrorCode)
			}
		}
	}
	return &apiResponse, nil
}
