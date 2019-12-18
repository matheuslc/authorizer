package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	ac "github.com/matheuslc/authorizer/internal/account"
	"github.com/matheuslc/authorizer/internal/eventstore/memorystore"
	tr "github.com/matheuslc/authorizer/internal/transaction"
)

func main() {
	accountStore := memorystore.NewStorage("account")
	transactionStore := memorystore.NewStorage("transaction")

	dec := json.NewDecoder(os.Stdin)
	json := make(map[string]map[string]interface{})

	for {
		err := dec.Decode(&json)

		if err == io.EOF {
			return
		}

		if err != nil {
			log.Fatal(err)
		}

		if json["transaction"] != nil {
			event := tr.CommandHandler(json["transaction"], &accountStore, &transactionStore)

			fmt.Println(tr.JSONPresenter(event))
		} else if json["account"] != nil {
			event := ac.CommandHandler(json["account"], &accountStore)

			fmt.Println(ac.JSONPresenter(event))
		}
	}
}
