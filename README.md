<p align="center">
  <a href="https://www.panaccess.com">
    <img src="https://www.panaccess.com/wp-content/uploads/panaccess_logo1.png" alt="Panaccess" style="display: block;margin-left: auto;margin-right: auto;width: 40%;" />
  </a>

  <h3 align="center">Panaccess-GO</h3>

  <p align="center">
    A useful library for using panaccess API using the same javascripts requests as the CableView from panaccess, developed in go to make an easy interaction with all panaccess functions.
    <a href="https://github.com/cdavid14/panaccess-go/issues/new?template=feature.md&labels=feature">Request feature</a>
  </p>
</p>


## Table of contents

- [Status](#status)
- [Good Practices](#good-practices)
- [Example Code](#example-code)
- [Bugs and feature requests](#bugs-and-feature-requests)
- [Contributing](#contributing)
- [Creators](#creators)


## Status

- There are some classes implemented and some limited functions, there are more but i have developed only the functions that are useful for me at this moment, new improvements are welcome, please use issues and PR as much as you can
- Tests implement
- A logo for the repository

## Good practices

- Requests are expensives for panaccess, so try to make only the necessaries requests, if you don't need to deal with variable data, try to maintain a local cache,a request example:
```
login
getSmartcards?
ERROR!
loggedIn?
yes!
return ERROR
```
And
```
login
getSmartcards?
ERROR!
loggedIn?
NO!
login
getSmartcards?
return smartcards or ERROR
```

## Example code

```golang
package main

import (
	"log"
	"net/http"
	"net/url"
	"time"
	"github.com/cdavid14/panaccess-go"
)

func main() {
  pan := panaccess.Panaccess{
		Servers:  []string{"https://cv01.panaccess.com", "https://cv01a.panaccess.com", "https://cv01b.panaccess.com"},
		User:     "demo",
		Password: "demo2010",
		Token:    "nziyNTEQsBbwvRWxLXzo",
		HTTP:     &http.Client{Timeout: time.Duration(30) * time.Second},
	}
	err := pan.Login()
	if err != nil {
		log.Fatalf("Panaccess login incorrect or servers not available: %v", err)
  }
  
  subObj := panaccess.Subscriber{}
  params := url.Values{}
	
	params.Add("limit", "1000")
	subs, err := subObj.Get(&pan, &params)
	if err != nil {
		log.Fatalf("Failed Get subs: %v", err)
  }
  ...
}
```
## Bugs and feature requests

Have a bug or a feature request? Please first read the [issue guidelines](https://github.com/cdavid14/blob/master/CONTRIBUTING.md) and search for existing and closed issues. If your problem or idea is not addressed yet, [please open a new issue](https://github.com/cdavid14/issues/new).

## Contributing

Please read through our [contributing guidelines](https://github.com/cdavid14/blob/master/CONTRIBUTING.md). Included are directions for opening issues, coding standards, and notes on development.

## Creators

**Christian David - Telecab**

- <https://github.com/cdavid14>

