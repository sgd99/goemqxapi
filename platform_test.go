package goemqxapi_test

import (
	"testing"

	"github.com/sgd99/goemqxapi"
)

func TestClientsReq(t *testing.T) {
	q := goemqxapi.ClientsRequest{
		ConnState: "connected",
	}
	t.Log(q.QueryString())
}
