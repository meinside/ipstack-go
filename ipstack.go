package ipstack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	apiBaseURL     = "https://api.ipstack.com"
	apiBaseURLFree = "http://api.ipstack.com" // free plan doesn't support https yet
)

// Client struct
type Client struct {
	free       bool
	accessKey  string
	httpClient *http.Client
}

// ResponseError for response error
//
// https://ipstack.com/documentation#errors
type ResponseError struct {
	Success bool `json:"success"`
	Error   struct {
		Code int    `json:"code"`
		Type string `json:"type"`
		Info string `json:"info"`
	} `json:"error"`
}

// Response for responses
//
// https://ipstack.com/documentation#apiresponse
type Response struct {
	IP            string  `json:"ip"`
	Hostname      string  `json:"hostname"`
	Type          string  `json:"type"`
	ContinentCode string  `json:"continent_code"`
	ContinentName string  `json:"continent_name"`
	CountryCode   string  `json:"country_code"`
	CountryName   string  `json:"country_name"`
	RegionCode    string  `json:"region_code"`
	RegionName    string  `json:"region_name"`
	City          string  `json:"city"`
	Zip           string  `json:"zip"`
	Latitude      float32 `json:"latitude"`
	Longitude     float32 `json:"longitude"`
	Location      struct {
		GeonameID int    `json:"geoname_id"`
		Capital   string `json:"capital"`
		Languages []struct {
			Code   string `json:"code"`
			Name   string `json:"name"`
			Native string `json:"native"`
		} `json:"languages"`
		CountryFlag             string `json:"country_flag"`
		CountryFlagEmoji        string `json:"country_flag_emoji"`
		CountryFlagEmojiUnicode string `json:"country_flag_emoji_unicode"`
		CallingCode             string `json:"calling_code"`
		IsEU                    bool   `json:"is_eu"`
	} `json:"location"`
	Timezone struct {
		ID               string `json:"id"`
		CurrentTime      string `json:"current_time"`
		GMTOffset        int32  `json:"gmt_offset"`
		Code             string `json:"code"`
		IsDaylightSaving bool   `json:"is_daylight_saving"`
	} `json:"time_zone"`
	Currency struct {
		Code         string `json:"code"`
		Name         string `json:"name"`
		Plural       string `json:"plural"`
		Symbol       string `json:"symbol"`
		SymbolNative string `json:"symbol_native"`
	} `json:"currency"`
	Connection struct {
		ASN string `json:"asn"`
		ISP string `json:"isp"`
	} `json:"connection"`
	Security struct {
		IsProxy     bool     `json:"is_proxy"`
		ProxyType   string   `json:"proxy_type"`
		IsCrawler   bool     `json:"is_crawler"`
		CrawlerName string   `json:"crawler_name"`
		CrawlerType string   `json:"crawler_type"`
		IsTor       bool     `json:"is_tor"`
		ThreatLevel string   `json:"threat_level"`
		ThreatTypes []string `json:"threat_types"`
	} `json:"security"`

	ResponseError
}

// NewClient returns a new client
func NewClient(accessKey string, isFreePlan bool) *Client {
	return &Client{
		free:      isFreePlan,
		accessKey: accessKey,
		httpClient: &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 300 * time.Second,
				}).Dial,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
	}
}

// LookupStandard requests standard lookup
//
// https://ipstack.com/documentation#standard
func (c *Client) LookupStandard(ip string) (response Response, err error) {
	var url string
	if c.free {
		url = fmt.Sprintf("%s/%s", apiBaseURLFree, ip)
	} else {
		url = fmt.Sprintf("%s/%s", apiBaseURL, ip)
	}

	var bytes []byte
	if bytes, err = c.get(url, map[string]string{}, map[string]interface{}{}); err == nil {
		if err = json.Unmarshal(bytes, &response); err == nil {
			if response.Error.Code > 0 {
				err = fmt.Errorf("error %d: %s (%s)", response.Error.Code, response.Error.Type, response.Error.Info)
			} else {
				return response, nil
			}
		}
	}

	return Response{}, err
}

// LookupBulk requests bulk lookup
// (not supported for free plan)
//
// https://ipstack.com/documentation#bulk
func (c *Client) LookupBulk(ips []string) (response []Response, err error) {
	if c.free {
		return []Response{}, fmt.Errorf("free plan does not support bulk lookup")
	}

	url := fmt.Sprintf("%s/%s", apiBaseURL, strings.Join(ips, ","))

	var bytes []byte
	if bytes, err = c.get(url, map[string]string{}, map[string]interface{}{}); err == nil {
		if err = json.Unmarshal(bytes, &response); err == nil {
			return response, nil
		}

		var errRes ResponseError
		if err = json.Unmarshal(bytes, &errRes); err == nil {
			err = fmt.Errorf("error %d: %s (%s)", errRes.Error.Code, errRes.Error.Type, errRes.Error.Info)
		}
	}

	return []Response{}, err
}

// LookupRequester requests lookup for this requester
//
// https://ipstack.com/documentation#requester
func (c *Client) LookupRequester() (response Response, err error) {
	var url string
	if c.free {
		url = fmt.Sprintf("%s/check", apiBaseURLFree)
	} else {
		url = fmt.Sprintf("%s/check", apiBaseURL)
	}

	var bytes []byte
	if bytes, err = c.get(url, map[string]string{}, map[string]interface{}{}); err == nil {
		if err = json.Unmarshal(bytes, &response); err == nil {
			if response.Error.Code > 0 {
				err = fmt.Errorf("error %d: %s (%s)", response.Error.Code, response.Error.Type, response.Error.Info)
			} else {
				return response, nil
			}
		}
	}

	return Response{}, err
}

// HTTP GET
func (c *Client) get(apiURL string, headers map[string]string, params map[string]interface{}) ([]byte, error) {
	// set default params
	if params == nil {
		params = map[string]interface{}{}
	}
	params["access_key"] = c.accessKey
	params["hostname"] = 1
	params["output"] = "json"
	if !c.free {
		params["security"] = 1
	}

	var err error
	var req *http.Request
	if req, err = http.NewRequest("GET", apiURL, nil); err == nil {
		// set HTTP headers
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		// set parameters
		queries := req.URL.Query()
		for key, value := range params {
			queries.Add(key, fmt.Sprintf("%v", value))
		}
		req.URL.RawQuery = queries.Encode()

		var resp *http.Response
		resp, err = c.httpClient.Do(req)

		if resp != nil {
			defer resp.Body.Close()
		}

		var bytes []byte
		if bytes, err = ioutil.ReadAll(resp.Body); err == nil {
			if resp.StatusCode == 200 {
				return bytes, nil
			}
		}
	}

	return []byte{}, err
}
