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
func AllCrawlingTarget() ([]Stock, error)  {
	fmt.Println("START AllCrawlingTarget")
	defer fmt.Println("END AllCrawlingTarget")
	tbl := NewTable(tableName)
	var stocks []Stock
	err := tbl.Scan().All(&stocks)

	if err != nil {
		return nil, err
	}

	if len(stocks) == 0 {
		return nil, errors.New("Not found targets")
	}

	return stocks, nil
}

// UpdatePrice update to the latest price
func UpdatePrice(stocks []Stock) error {
	fmt.Println("START UpdatePrice")
	defer fmt.Println("END UpdatePrice")
	wg := new(sync.WaitGroup)
	errCodes := []string{}
	for _, stock := range stocks {
		wg.Add(1)
		go func(target Stock) {
			defer wg.Done()
			crw := Crawler{Code: target.Code}
			
			price, err := crw.ScrapePrice()
		
			if err != nil {
				fmt.Println(err.Error())
				errCodes = append(errCodes, target.Code)
				return
			}

			tbl := NewTable(tableName)
			err = tbl.Update("code", target.Code).Set("price", price).Run()

			if err != nil {
				fmt.Println(err.Error())
				errCodes = append(errCodes, target.Code)
				return
			}
		}(stock)
	}
	wg.Wait()

	if len(errCodes) != 0 {
		return errors.New("Can not update stock price")
	}

	return nil
}