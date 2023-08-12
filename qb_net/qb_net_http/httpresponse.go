package qb_net_http

import (
	"io/ioutil"
	"net/http"

	"github.com/rskvp/qb-core/qb_utils"
)

type ResponseData struct {
	StatusCode int
	Body       []byte
	Header     map[string]string
}

func NewResponseDataEmpty() *ResponseData {
	instance := new(ResponseData)
	return instance
}

func NewResponseData(res *http.Response) (instance *ResponseData, err error) {
	instance = new(ResponseData)
	instance.StatusCode = res.StatusCode
	instance.Header = HttpHeaderToMap(res.Header)
	instance.Body, err = ioutil.ReadAll(res.Body)

	return
}

func (instance *ResponseData) String() string {
	m := map[string]interface{}{
		"status": instance.StatusCode,
		"header": instance.Header,
		"body":   string(instance.Body),
	}
	return qb_utils.JSON.Stringify(m)
}

func (instance *ResponseData) BodyAsMap() map[string]interface{} {
	var m map[string]interface{}
	err := qb_utils.JSON.Read(instance.Body, &m)
	if nil == err {
		return m
	}
	return map[string]interface{}{}
}
