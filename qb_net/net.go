package qb_net

import "github.com/rskvp/qb-core/qb_net/qb_net_http"

type NetHelper struct {
}

var Net *NetHelper

func init() {
	Net = new(NetHelper)
}

func (*NetHelper) NewHttpClient() *qb_net_http.HttpClient {
	return qb_net_http.NewHttpClient()
}
