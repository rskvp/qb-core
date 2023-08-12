package qb_net_http

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/rskvp/qb-core/qb_utils"
)

type HttpClientOptions struct {
	InsecureSkipVerify bool
	Timeout            time.Duration
	Header             map[string]string
}

type HttpClient struct {
	options *HttpClientOptions
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewHttpClient(options ...*HttpClientOptions) *HttpClient {
	instance := new(HttpClient)
	if len(options) == 1 {
		instance.options = options[0]
	} else {
		instance.options = new(HttpClientOptions)
		instance.options.InsecureSkipVerify = true
	}
	if nil == instance.options.Header {
		instance.options.Header = make(map[string]string)
	}
	if instance.options.Timeout == 0 {
		instance.options.Timeout = time.Second * 120
	}
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpClient) AddHeader(key, value string) {
	if nil != instance && nil != instance.options {
		if nil == instance.options.Header {
			instance.options.Header = make(map[string]string)
		}
		instance.options.Header[key] = value
	}
}

func (instance *HttpClient) RemoveHeader(key string) {
	if nil != instance && nil != instance.options && nil != instance.options.Header {
		delete(instance.options.Header, key)
	}
}

func (instance *HttpClient) Get(url string) (response *ResponseData, err error) {
	if nil != instance {
		client := instance.client()
		response, err = do(client, MethodGet, url, instance.options.Header, nil)
	}
	return
}

func (instance *HttpClient) Post(url string, body interface{}) (response *ResponseData, err error) {
	if nil != instance {
		client := instance.client()
		response, err = do(client, MethodPost, url, instance.options.Header, body)
	}
	return
}

func (instance *HttpClient) Delete(url string, body interface{}) (response *ResponseData, err error) {
	if nil != instance {
		client := instance.client()
		response, err = do(client, MethodDelete, url, instance.options.Header, body)
	}
	return
}

func (instance *HttpClient) Put(url string, body interface{}) (response *ResponseData, err error) {
	if nil != instance {
		client := instance.client()
		response, err = do(client, MethodPut, url, instance.options.Header, body)
	}
	return
}

func (instance *HttpClient) Upload(url string, filename, optParamName string, optParams map[string]interface{}) (*ResponseData, error) {
	return instance.UploadTimeout(url, filename, optParamName, optParams, time.Second*120)
}

func (instance *HttpClient) UploadTimeout(url string, filename, optParamName string, optParams map[string]interface{}, timeout time.Duration) (*ResponseData, error) {
	if nil == optParams {
		optParams = make(map[string]interface{})
	}
	if len(optParamName) == 0 {
		optParamName = "file"
	}

	// open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	part, err := writer.CreateFormFile(optParamName, filepath.Base(filename))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	for key, val := range optParams {
		var sval string
		sval = qb_utils.Convert.ToString(val)
		_ = writer.WriteField(key, sval)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if nil != instance.options.Header {
		// fmt.Println(req.Header.String())
		for k, v := range instance.options.Header {
			req.Header.Set(k, v)
		}
	}

	client := instance.client()
	client.Timeout = timeout
	resp, err := client.Do(req)
	if nil != err {
		return nil, err
	}
	defer resp.Body.Close()

	response := NewResponseDataEmpty()
	response.Body, _ = ioutil.ReadAll(resp.Body)
	response.StatusCode = resp.StatusCode
	response.Header = HttpHeaderToMap(resp.Header)
	return response, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpClient) client() (client *http.Client) {
	if nil != instance {
		client = new(http.Client)
		client.Timeout = instance.options.Timeout

		if instance.options.InsecureSkipVerify {
			transCfg := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
			}
			client.Transport = transCfg
		}
	}
	return
}

func do(client *http.Client, method string, uri string, header map[string]string, payload interface{}) (response *ResponseData, err error) {
	if nil != client {
		var req *http.Request
		req, err = http.NewRequest(method, uri, ReadPayload(payload))
		if nil == err {
			// headers
			if nil != header {
				for k, v := range header {
					req.Header.Add(k, v)
				}
			}

			// client
			var res *http.Response
			res, err = client.Do(req)
			if nil == err {
				defer res.Body.Close()
				response, err = NewResponseData(res)
			}
		}
	}
	return
}
