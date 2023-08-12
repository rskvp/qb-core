package qb_net_http

import (
	"io"
	"net/http"
	"strings"

	"github.com/rskvp/qb-core/qb_utils"
)

const MethodGet = "GET"
const MethodPost = "POST"
const MethodPut = "PUT"
const MethodDelete = "DELETE"

func HttpHeaderToMap(header http.Header) map[string]string {
	response := make(map[string]string)
	for k, v := range header {
		response[k] = strings.Join(v, ",")
	}
	return response
}

func ReadPayload(rawPayload interface{}) io.Reader {
	if nil != rawPayload {
		var payload string
		if v, b := rawPayload.(string); b {
			payload = v
		} else if v, b := rawPayload.([]byte); b {
			payload = string(v)
		} else if v, b := rawPayload.([]uint8); b {
			payload = string(v)
		} else if v, b := rawPayload.(map[string]interface{}); b {
			payload = qb_utils.JSON.Stringify(v)
		}
		return strings.NewReader(payload)
	}
	return nil
}
