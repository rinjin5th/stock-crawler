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
	PurchasePrice int `dynamo:"purchase_price"`
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
	fmt.Println("START UpdatePriceXXX")
	defer fmt.Println("END UpdatePrice")
	wg := new(sync.WaitGroup)
	errCodes := []string{}
	for _, stock := range stocks {
		wg.Add(1)
		go func(target Stock) {
			defer wg.Done()
			crw := Crawler{Code: target.Code}
			
			price, isNoPrice, err := crw.ScrapePrice()
		
			if err != nil {
				fmt.Println(err.Error())
				errCodes = append(errCodes, target.Code)
				return
			}

			if isNoPrice || target.Price == price {
				fmt.Sprintf("No update %s", target.Code)
				return
			}
			alert(target, price)
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

func alert(stock Stock, scrapedPrice int) (){
	fmt.Println("START alert")
	defer fmt.Println("END alert")
	
	if stock.Price == 0 || stock.Price == scrapedPrice {
		return
	}
	diff := scrapedPrice - stock.PurchasePrice

	fmt.Printf("diff -> %d", diff)

	if diff <= LowerLimit {
		slack := NewSlack(fmt.Sprintf("<!here>\n%s is bad price. %d yen", stock.Code, scrapedPrice))
		slack.Send(SlackWebHookURL)
	} else if diff >= UpperLimit {
		slack := NewSlack(fmt.Sprintf("<!here>\n%s is good price. %d yen", stock.Code, scrapedPrice))
		slack.Send(SlackWebHookURL)
	}
}
