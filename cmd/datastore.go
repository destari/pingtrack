package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/buntdb"
	"time"
)

var store *buntdb.DB

func OpenStore(dbfilepath string) error {
	var err error
	store, err = buntdb.Open(dbfilepath)
	if err != nil {
		return err
	}
	store.CreateIndex("hostname", "*", buntdb.IndexJSON("Host"))
	store.CreateIndex("time", "*", buntdb.IndexJSON("Time"))

	return nil
}

func CloseStore() {
	store.Close()
}

func StoreResult(result Results) error {

	key := fmt.Sprintf("%s:%d", result.Host, result.EpochTime)
	pattern := fmt.Sprintf("%s:*", result.Host)

	store.CreateIndex(result.Host, pattern, buntdb.IndexJSON("Time"))

	jsonResults, _ := json.Marshal(result)

	err := store.Update(func(tx *buntdb.Tx) error {
		tx.Set(key, string(jsonResults), &buntdb.SetOptions{Expires:true, TTL:time.Second*time.Duration(ttl*24*60*60)})
		return nil
	})

	return err
}

func StoreRetrieve(host string, from int64) []json.RawMessage {
	results := []json.RawMessage{}

	store.View(func(tx *buntdb.Tx) error {
		pivot := fmt.Sprintf("{\"Time\":%d}", from)
		tx.AscendGreaterOrEqual(host, pivot, func(key, value string) bool {
			//fmt.Printf("%s: %s\n", key, value)
			results = append(results, json.RawMessage(value))
			return true
		})
		return nil
	})

	return results
}