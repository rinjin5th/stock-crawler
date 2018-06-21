package main

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

const (
	sbiURL = "https://www.sbisec.co.jp/ETGate"
)

const (
	companyName = iota
	price
)

// Crawler gets stock info from sbi
type Crawler struct {
	Code string
	// CompanyName string
}

// ScrapePrice gets stock price from sbi
func (crw Crawler) ScrapePrice() (int, error) {
	if len(crw.Code) == 0 {
		return -1, errors.New("Must set stock code")
	}

	resp, err := http.PostForm(sbiURL, newParams(crw.Code))
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return -1, errors.New("Can not crawl to sbi")
	}

	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	priceVal := -1
	doc.Find("span.fxx01").Each(func(i int, s *goquery.Selection) {
		switch i {
		case companyName:
			// nop
		case price:
			priceVal, err = strconv.Atoi(strings.Replace(s.Text(), ",", "", -1))
		}
	})

	if err != nil {
		return -1, errors.New("Illegal price value")
	}

	if priceVal < 0 {
		return -1, errors.New("Can not scrape price")
	}
	
	return priceVal, nil
}

func newParams(code string) url.Values {
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
	values.Set("i_stock_sec", code)
	values.Set("json_status", "1")
	values.Set("json_content", "")
	values.Set("i_output_type", "0")
	values.Set("qr_keyword", "1")
	values.Set("qr_suggest", "1")
	values.Set("qr_sort", "1")

	return values
}
