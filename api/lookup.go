package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/bjartek/overflow"
	"github.com/onflow/flow-go-sdk"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Specify id and/or network", http.StatusInternalServerError)
		return
	}
	network := r.URL.Query().Get("network")
	if network == "" {
		network = "mainnet"
	}
	o := Overflow(WithNetwork(network))

	t, err := o.GetTransactionById(context.Background(), flow.HexToID(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot find transaction with id %s on network %s error:%v", id, network, err), http.StatusNotFound)
		return
	}

	res, err := json.Marshal(t)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot marshal result as json for transaction with id %s on network %s error:%v", id, network, err), http.StatusNotFound)
		return
	}

	_, err = w.Write(res)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing resultfor transaction with id %s on network %s error:%v", id, network, err), http.StatusNotFound)
		return
	}

}
