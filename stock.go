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
	fmt.Println("START AllCrawlingTarget")
	tbl := NewTable(tableName)
	var stocks []Stock
	tbl.Scan().All(&stocks)

	fmt.Println(len(stocks))

	fmt.Println("END AllCrawlingTarget")

	return stocks
}

// UpdatePrice update to the latest price
func UpdatePrice(stocks []Stock) error {
	fmt.Println("START UpdatePrice")
	wg := new(sync.WaitGroup)
	fmt.Println("After create wg")
	errCodeCh := make(chan string)
	fmt.Println("After create ch")
	for _, stock := range stocks {
		fmt.Println("START Loop")
		wg.Add(1)
		go func(target Stock) {
			defer wg.Done()
			fmt.Println("START goroutine")
			crw := Crawler{Code: target.Code}
			
			price, err := crw.ScrapePrice()
		
			if err != nil {
				fmt.Println(err.Error())
				errCodeCh <- target.Code
				return
			}

			tbl := NewTable(tableName)
			err = tbl.Update("code", target.Code).Set("price", price).Run()

			if err != nil {
				fmt.Println(err.Error())
				errCodeCh <- target.Code
				return
			}
			fmt.Println("end goroutine")
		}(stock)
		fmt.Println("END Loop")
	}
	errCodes := []string{}
	for code := range errCodeCh {
		errCodes = append(errCodes, code)
	}
	wg.Wait()

	if len(errCodes) != 0 {
		return errors.New("Can not update stock price")
	}
	fmt.Println("END UpdatePrice")

	return nil
}