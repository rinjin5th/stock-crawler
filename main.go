package main

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/PuerkitoBio/goquery"
)

const (
	sbiURL = "https://www.sbisec.co.jp/ETGate"
)

const (
	companyName = iota
	price
)

// Handler is executed by AWS Lambda in the main function. Once the request
// is processed, it returns an Amazon API Gateway response object to AWS Lambda
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	jar, err := cookiejar.New(nil)
    if err != nil {
        return events.APIGatewayProxyResponse{}, err
    }

    client := &http.Client{
        Jar: jar,
	}
	
	values := url.Values{}
	values.Set("_ControlID", "WPLETsiR001Control")
	values.Set("_PageID", "WPLETsiR001Iser10")
	values.Set("_DataStoreID", "DSWPLETsiR001Control")
	values.Set("_ActionID", "clickToSearchStockPriceJP")
	values.Set("i_dom_flg", "1")
	values.Set("ref_from", "1")
	values.Set("ref_to", "20")
	values.Set("wstm4130_sort_id", "++")
	values.Set("wstm4130_sort_kbn", "+")
	values.Set("i_exchange_code", "JPN")
	values.Set("i_stock_sec", "2685")
	values.Set("json_status", "1")
	values.Set("json_content", "")
	values.Set("i_output_type", "0")
	values.Set("qr_keyword", "1")
	values.Set("qr_suggest", "1")
	values.Set("qr_sort", "1")

	resp, err := client.PostForm(sbiURL, values)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
	    return events.APIGatewayProxyResponse{}, err
	}

	var stockPrice string
	doc.Find("span.fxx01").Each(func(i int, s *goquery.Selection) {
		switch i {
			case companyName:
				// nop
			case price: 
				stockPrice = strings.Replace(s.Text(), ",", "", -1)
		}
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       stockPrice,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}, nil

}

func main() {
	lambda.Start(Handler)
}
