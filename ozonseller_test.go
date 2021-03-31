package ozonseller

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Client(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch req.URL.String() {
		case "/v2/product/info":
			const data = `{
"result": {
  "barcode": "",
  "buybox_price": "",
  "category_id": 17034461,
  "created_at": "2019-11-26T10:40:44.940Z",
  "errors": [
    {
      "field": "string",
      "attribute_id": 0,
      "code": "string",
      "level": "string"
    }
  ],
  "id": 7154396,
  "images": [
    "https://cdn1.ozone.ru/multimedia/1028110514.jpg"
  ],
  "marketing_price": "",
  "min_ozon_price": "3599.0000",
  "name": "Туалетная вода VALENTINO UOMO ACQUA spray 75 ml",
  "offer_id": "item_6060091",
  "old_price": "",
  "premium_price": "",
  "price": "3599.0000",
  "recommended_price": " ",
  "sources": [
    {
      "is_enabled": true,
      "sku": 150583609,
      "source": "fbo"
    }
  ],
  "state": "processed",
  "stocks": {
    "coming": 0,
    "present": 120,
    "reserved": 0
  },
  "vat": "0.2",
  "visibility_details": {
    "active_product": true,
    "has_price": true,
    "has_stock": true
  },
  "visible": true
}
}`
			rw.Write([]byte(data))
		case "/v2/product/info/stocks":
			const data = `{"result":{"items": [{"offer_id":"aa", "stocks": []}], "total": 1}}`
			rw.Write([]byte(data))
		}
	}))
	defer server.Close()
	c := ClientV2{
		Host:     server.URL,
		Client:   server.Client(),
		ClientID: DefaultTestClientID,
	}

	t.Run("product info unmarshaling", func(t *testing.T) {
		info, err := c.ProductInfo(ProductInfoFilter{
			OfferID:   "0",
			ProductID: 0,
			SKU:       0,
		})
		if err != nil {
			t.Fatal(err)
		}

		t.Log(info)
	})

	t.Run("product stock", func(t *testing.T) {
		stocks, err := c.ProductInfoStocks()
		if err != nil {
			t.Fatal(err)
		}

		if stocks[0].OfferID != "aa" {
			t.Fatal("not matched offer_id, want: 'aa', but got", stocks[0].OfferID)
		}
	})
}
