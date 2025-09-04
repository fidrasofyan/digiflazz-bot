package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/fidrasofyan/digiflazz-bot/internal/config"
)

var digiflazzHttpClient = &http.Client{
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:   true,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 5 * time.Second,
	},
	Timeout: 10 * time.Second,
}

type DigiflazzErrorResponse struct {
	Data struct {
		RC      string `json:"rc"`
		Message string `json:"message"`
	} `json:"data"`
}

func (e *DigiflazzErrorResponse) Error() string {
	return fmt.Sprintf("Kode: %s - Pesan: %s", e.Data.RC, e.Data.Message)
}

// Sign
func DigiflazzSign(suffix string) string {
	h := md5.New()
	h.Write([]byte(config.Cfg.DigiflazzUsername))
	h.Write([]byte(config.Cfg.DigiflazzApiKey))
	h.Write([]byte(suffix))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Check balance
type DigiflazzCheckBalanceResponse struct {
	Data struct {
		Deposit int `json:"deposit"`
	} `json:"data"`
}

func DigiflazzCheckBalance(ctx context.Context) (*DigiflazzCheckBalanceResponse, error) {
	data := struct {
		Cmd      string `json:"cmd"`
		Username string `json:"username"`
		Sign     string `json:"sign"`
	}{
		Cmd:      "deposit",
		Username: config.Cfg.DigiflazzUsername,
		Sign:     DigiflazzSign("depo"),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	url := config.Cfg.DigiflazzBaseUrl + "/cek-saldo"

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := digiflazzHttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		if res.StatusCode == 404 {
			return nil, errors.New("404 page not found")
		}
		var digiflazzError DigiflazzErrorResponse
		err = json.NewDecoder(res.Body).Decode(&digiflazzError)
		if err != nil {
			return nil, err
		}
		return nil, &digiflazzError
	}

	var response DigiflazzCheckBalanceResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Get prepaid products
type DigiflazzPrepaidProduct struct {
	ProductName         string `json:"product_name"`
	Category            string `json:"category"`
	Brand               string `json:"brand"`
	Type                string `json:"type"`
	SellerName          string `json:"seller_name"`
	Price               int64  `json:"price"`
	BuyerSKUCode        string `json:"buyer_sku_code"`
	BuyerProductStatus  bool   `json:"buyer_product_status"`
	SellerProductStatus bool   `json:"seller_product_status"`
	UnlimitedStock      bool   `json:"unlimited_stock"`
	Stock               int64  `json:"stock"`
	Multi               bool   `json:"multi"`
	StartCutOff         string `json:"start_cut_off"`
	EndCutOff           string `json:"end_cut_off"`
	Description         string `json:"desc"`
}

type DigiflazzGetPrepaidPriceListResponse struct {
	Data []DigiflazzPrepaidProduct `json:"data"`
}

func DigiflazzGetPrepaidPriceList(ctx context.Context) (*DigiflazzGetPrepaidPriceListResponse, error) {
	data := struct {
		Cmd      string `json:"cmd"`
		Username string `json:"username"`
		Sign     string `json:"sign"`
	}{
		Cmd:      "prepaid",
		Username: config.Cfg.DigiflazzUsername,
		Sign:     DigiflazzSign("pricelist"),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	url := config.Cfg.DigiflazzBaseUrl + "/price-list"

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := digiflazzHttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		if res.StatusCode == 404 {
			return nil, errors.New("404 page not found")
		}
		var digiflazzError DigiflazzErrorResponse
		err = json.NewDecoder(res.Body).Decode(&digiflazzError)
		if err != nil {
			return nil, err
		}
		return nil, &digiflazzError
	}

	var response DigiflazzGetPrepaidPriceListResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Create transaction
type DigiflazzTrxStatus string

var (
	DigiflazzTrxStatusPending DigiflazzTrxStatus = "Pending"
	DigiflazzTrxStatusSuccess DigiflazzTrxStatus = "Sukses"
	DigiflazzTrxStatusFailed  DigiflazzTrxStatus = "Gagal"
)

type DigiflazzTrxData struct {
	RefID        string             `json:"ref_id"`
	CustomerNo   string             `json:"customer_no"`
	BuyerSKUCode string             `json:"buyer_sku_code"`
	Message      string             `json:"message"`
	Status       DigiflazzTrxStatus `json:"status"`
	RC           string             `json:"rc"`
	SN           *string            `json:"sn"`
	Price        int32              `json:"price"`
}

type DigiflazzCreateTrxResponse struct {
	Data DigiflazzTrxData `json:"data"`
}

type DigiflazzCreateTrxParams struct {
	RefID        string
	BuyerSKUCode string
	CustomerNo   string
}

func DigiflazzCreateTrx(ctx context.Context, params *DigiflazzCreateTrxParams) (*DigiflazzCreateTrxResponse, error) {
	data := struct {
		Username     string `json:"username"`
		BuyerSKUCode string `json:"buyer_sku_code"`
		CustomerNo   string `json:"customer_no"`
		RefID        string `json:"ref_id"`
		Sign         string `json:"sign"`
		Testing      bool   `json:"testing"`
	}{
		Username:     config.Cfg.DigiflazzUsername,
		BuyerSKUCode: params.BuyerSKUCode,
		CustomerNo:   params.CustomerNo,
		RefID:        params.RefID,
		Sign:         DigiflazzSign(params.RefID),
		Testing:      config.Cfg.AppEnv != "production",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	url := config.Cfg.DigiflazzBaseUrl + "/transaction"

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := digiflazzHttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		if res.StatusCode == 404 {
			return nil, errors.New("404 page not found")
		}
		var digiflazzError DigiflazzErrorResponse
		err = json.NewDecoder(res.Body).Decode(&digiflazzError)
		if err != nil {
			return nil, err
		}
		return nil, &digiflazzError
	}

	var response DigiflazzCreateTrxResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
