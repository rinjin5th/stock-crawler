package main

import (
	"sync"
	"errors"
	"fmt"
)

const (
	tableName = "stock"
)

// Stock is stock object for DynamoDB
type Stock struct {
	Code  string `dynamo:"code"`
	CompanyName string `dynamo:"company_name"`
	Price int `dynamo:"price"`
}

// AllCrawlingTarget gets crawling target from DynamoDB
func AllCrawlingTarget() ([]Stock) {
	fmt.Print("START AllCrawlingTarget")
	tbl := NewTable(tableName)
	var stocks []Stock
	tbl.Scan().All(&stocks)
	fmt.Print("END AllCrawlingTarget")

	return stocks
}

// UpdatePrice update to the latest price
func UpdatePrice(stocks []Stock) error {
	fmt.Print("START UpdatePrice")
	wg := new(sync.WaitGroup)
	errCodeCh := make(chan string)
	for _, stock := range stocks {
		wg.Add(1)
		go func(target Stock) {
			defer wg.Done()
			crw := Crawler{Code: target.Code}
			
			price, err := crw.ScrapePrice()
		
			if err != nil {
				errCodeCh <- target.Code
				return
			}

			tbl := NewTable(tableName)
			err = tbl.Update("code", target.Code).Set("price", price).Run()

			if err != nil {
				errCodeCh <- target.Code
				return
			}
		}(stock)
	}
	errCodes := []string{}
	for code := range errCodeCh {
		errCodes = append(errCodes, code)
	}
	wg.Wait()

	if len(errCodes) != 0 {
		return errors.New("Can not update stock price")
	}
	fmt.Print("END UpdatePrice")

	return nil
}