package ozonseller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	// DefaultTestHost, DefaultTestClientID и DefaultTestApiKey это параметры
	// для осуществления запросов в тестовую среду API
	// подробнее: https://api-seller.ozon.ru/apiref/ru/#t-title_sandbox
	DefaultTestHost     = "https://cb-api.ozonru.me"
	DefaultTestClientID = 836
	DefaultTestApiKey   = "0296d4f2-70a1-4c09-b507-904fd05567b9"
)

const (
	// DefaultProductionHost это адрес на который будут сделаны запросы.
	// Данный адрес вёдет работу с production окружением, будте бдительны
	// подробнее: https://api-seller.ozon.ru/apiref/ru/#t-title_env_production
	DefaultProductionHost = "https://api-seller.ozon.ru"
)

type ProductInfo struct {
	Barcode           string            `json:"barcode"`
	BuyboxPrice       string            `json:"buybox_price"`
	CategoryID        int               `json:"category_id"`
	CreatedAt         time.Time         `json:"created_at"`
	Errors            []Errors          `json:"errors"`
	ID                int               `json:"id"`
	Images            []string          `json:"images"`
	MarketingPrice    string            `json:"marketing_price"`
	MinOzonPrice      string            `json:"min_ozon_price"`
	Name              string            `json:"name"`
	OfferID           string            `json:"offer_id"`
	OldPrice          string            `json:"old_price"`
	PremiumPrice      string            `json:"premium_price"`
	Price             string            `json:"price"`
	RecommendedPrice  string            `json:"recommended_price"`
	Sources           []Sources         `json:"sources"`
	State             string            `json:"state"`
	Stocks            Stocks            `json:"stocks"`
	Vat               string            `json:"vat"`
	VisibilityDetails VisibilityDetails `json:"visibility_details"`
	Visible           bool              `json:"visible"`
}

type Errors struct {
	Field       string `json:"field"`
	AttributeID int    `json:"attribute_id"`
	Code        string `json:"code"`
	Level       string `json:"level"`
}

type Sources struct {
	IsEnabled bool   `json:"is_enabled"`
	SKU       int    `json:"sku"`
	Source    string `json:"source"`
}

type Stocks struct {
	Coming   int `json:"coming"`
	Present  int `json:"present"`
	Reserved int `json:"reserved"`
}

type VisibilityDetails struct {
	ActiveProduct bool `json:"active_product"`
	HasPrice      bool `json:"has_price"`
	HasStock      bool `json:"has_stock"`
}

type ProductInfoFilter struct {
	OfferID   string `json:"offer_id"`
	ProductID int    `json:"product_id"`
	SKU       int    `json:"sku"`
}

type Pagination struct {
	Page     uint `json:"page"`
	PageSize uint `json:"page_size"`
}

type ClientV2 struct {
	Host     string
	ClientID int
	ApiKey   string
	Client   *http.Client
}

func (c ClientV2) ProductInfo(filter ProductInfoFilter) (ProductInfo, error) {
	req, err := c.newPost("/v2/product/info", filter)
	if err != nil {
		return ProductInfo{}, err
	}

	var respStruct struct {
		Result ProductInfo `json:"result"`
	}

	return respStruct.Result, c.do(req, &respStruct)
}

type ProductInfoStocks struct {
	Present  int    `json:"present"`
	Reserved int    `json:"reserved"`
	Type     string `json:"type"`
}

type ProductInfoStock struct {
	OfferID   string              `json:"offer_id"`
	ProductID int                 `json:"product_id"`
	Stocks    []ProductInfoStocks `json:"stocks"`
}

func (c ClientV2) ProductInfoStocks() ([]ProductInfoStock, error) {
	const (
		pageSize  = 100
		startPage = 1
	)

	totalElements := -1

	var list []ProductInfoStock

	for page := uint(1); len(list) != totalElements; page++ {
		req, err := c.newPost("/v2/product/info/stocks", Pagination{
			Page:     page,
			PageSize: pageSize,
		})
		if err != nil {
			return list, err
		}

		type Result struct {
			Items []ProductInfoStock `json:"items"`
			Total int                `json:"total"`
		}
		var respStruct struct {
			Result Result `json:"result"`
		}

		err = c.do(req, &respStruct)
		if err != nil {
			return list, err
		}

		totalElements = respStruct.Result.Total
		list = append(list, respStruct.Result.Items...)
	}

	return list, nil
}

func (c ClientV2) newPost(path string, v interface{}) (*http.Request, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.Host+path, bytes.NewReader(b))
	if err != nil {
		return req, err
	}
	req.Header.Set("Client-Id", strconv.Itoa(c.ClientID))
	req.Header.Set("Api-Key", c.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c ClientV2) do(req *http.Request, v interface{}) error {
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}
