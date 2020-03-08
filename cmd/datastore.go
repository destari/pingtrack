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
	err = store.CreateIndex("hostname", "*", buntdb.IndexJSON("Host"))
	if err != nil {
		fmt.Println(err)
	}

	err = store.CreateIndex("time", "*", buntdb.IndexJSON("EpochTime"))
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func CloseStore() {
	store.Close()
}

func StoreResult(result Results) error {
	var err error

	key := fmt.Sprintf("%s:%d", result.Host, result.EpochTime)
	pattern := fmt.Sprintf("%s:*", result.Host)

	_ = store.CreateIndex(result.Host, pattern, buntdb.IndexJSON("EpochTime"))

	jsonResults, _ := json.Marshal(result)

	err = store.Update(func(tx *buntdb.Tx) error {
		tx.Set(key, string(jsonResults), &buntdb.SetOptions{Expires:true, TTL:time.Second*time.Duration(ttl*24*60*60)})
		return nil
	})

	return err
}

func StoreRetrieve(host string, from int64, duration int64) []json.RawMessage {
	results := []json.RawMessage{}

	start := fmt.Sprintf("{\"EpochTime\":%d}", from)
	end := fmt.Sprintf("{\"EpochTime\":%d}", from+duration)

	_ = store.View(func(tx *buntdb.Tx) error {
		_ = tx.AscendRange(host, start, end, func(key, value string) bool {
			//fmt.Printf("%s: %s\n", key, value)
			results = append(results, json.RawMessage(value))
			return true
		})
		return nil
	})
	fmt.Println(len(results))
	return results
}