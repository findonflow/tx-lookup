package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bjartek/overflow"
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
	o := overflow.Overflow(overflow.WithNetwork(network))

	t, err := o.GetOverflowTransactionById(context.Background(), flow.HexToID(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot find transaction with id %s on network %s error:%v", id, network, err), http.StatusNotFound)
		return
	}

	b, err := o.GetBlockById(context.Background(), t.BlockId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot find block for transaction with id %s on network %s error:%v", id, network, err), http.StatusNotFound)
		return
	}

	events := []Event{}
	for _, ev := range t.Events {
		events = append(events, Event{
			Name:   ev.Name,
			Fields: ev.Fields,
		})
	}

	arguments := []KV{}
	for _, a := range t.Arguments {
		arguments = append(arguments, KV{
			Key:   a.Key,
			Value: a.Value,
		})
	}

	imports := []Import{}
	for _, a := range t.Imports {
		imports = append(imports, Import{
			Name:    a.Name,
			Address: a.Address,
		})
	}

	hash := sha256.Sum256(t.Script)

	status := t.Status
	txErr := ""
	if t.Error != nil {
		txErr = t.Error.Error()
		if txErr != "" {
			status = "ERROR"

		}
	}

	tx := Transaction{
		Id:              t.Id,
		Block:           Block{Height: b.Height, ID: b.ID.String(), Time: b.Timestamp},
		Events:          events,
		Error:           txErr,
		Fee:             t.Fee,
		Status:          status,
		Arguments:       arguments,
		Authorizers:     t.Authorizers,
		Stakeholders:    t.Stakeholders,
		Payer:           t.Payer,
		ProposalKey:     t.ProposalKey,
		GasLimit:        t.GasLimit,
		GasUsed:         t.GasUsed,
		ExecutionEffort: t.ExecutionEffort,
		Body:            string(t.Script),
		BodyHash:        hex.EncodeToString(hash[:]),
		Import:          imports,
	}

	res, err := json.Marshal(tx)
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

type Transaction struct {
	Id              string
	Block           Block
	Events          []Event
	Error           string
	Fee             float64
	Status          string
	Arguments       []KV
	Authorizers     []string
	Stakeholders    map[string][]string
	Payer           string
	ProposalKey     flow.ProposalKey
	GasLimit        uint64
	GasUsed         uint64
	ExecutionEffort float64
	Body            string
	BodyHash        string
	Import          []Import
}

type KV struct {
	Key   string
	Value interface{}
}

type Import struct {
	Name    string
	Address string
}

type Event struct {
	Fields map[string]interface{}
	Name   string
}

type Block struct {
	Height uint64
	ID     string
	Time   time.Time
}
