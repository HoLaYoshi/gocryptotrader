package liqui

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/thrasher-/gocryptotrader/config"
	"github.com/thrasher-/gocryptotrader/currency/pair"
	"github.com/thrasher-/gocryptotrader/currency/symbol"
	exchange "github.com/thrasher-/gocryptotrader/exchanges"
)

var l Liqui

const (
	apiKey         = ""
	apiSecret      = ""
	canPlaceOrders = false
)

func TestSetDefaults(t *testing.T) {
	l.SetDefaults()
}

func TestSetup(t *testing.T) {
	cfg := config.GetConfig()
	cfg.LoadConfig("../../testdata/configtest.json")
	liquiConfig, err := cfg.GetExchangeConfig("Liqui")
	if err != nil {
		t.Error("Test Failed - liqui Setup() init error")
	}

	liquiConfig.AuthenticatedAPISupport = true
	liquiConfig.APIKey = apiKey
	liquiConfig.APISecret = apiSecret

	l.Setup(liquiConfig)
}

func TestGetAvailablePairs(t *testing.T) {
	t.Parallel()
	v := l.GetAvailablePairs(false)
	if len(v) != 0 {
		t.Error("Test Failed - liqui GetFee() error")
	}
}

func TestGetInfo(t *testing.T) {
	t.Parallel()
	_, err := l.GetInfo()
	if err != nil {
		t.Error("Test Failed - liqui GetInfo() error", err)
	}
}

func TestGetTicker(t *testing.T) {
	t.Parallel()
	_, err := l.GetTicker("eth_btc")
	if err != nil {
		t.Error("Test Failed - liqui GetTicker() error", err)
	}
}

func TestGetDepth(t *testing.T) {
	t.Parallel()
	_, err := l.GetDepth("eth_btc")
	if err != nil {
		t.Error("Test Failed - liqui GetDepth() error", err)
	}
}

func TestGetTrades(t *testing.T) {
	t.Parallel()
	_, err := l.GetTrades("eth_btc")
	if err != nil {
		t.Error("Test Failed - liqui GetTrades() error", err)
	}
}

func TestAuthRequests(t *testing.T) {
	if l.APIKey != "" && l.APISecret != "" {
		_, err := l.GetAccountInfo()
		if err == nil {
			t.Error("Test Failed - liqui GetAccountInfo() error", err)
		}

		_, err = l.Trade("", "", 0, 1)
		if err == nil {
			t.Error("Test Failed - liqui Trade() error", err)
		}

		_, err = l.GetActiveOrders("eth_btc")
		if err == nil {
			t.Error("Test Failed - liqui GetActiveOrders() error", err)
		}

		_, err = l.GetOrderInfo(1337)
		if err == nil {
			t.Error("Test Failed - liqui GetOrderInfo() error", err)
		}

		_, err = l.CancelExistingOrder(1337)
		if err == nil {
			t.Error("Test Failed - liqui CancelExistingOrder() error", err)
		}

		_, err = l.GetTradeHistory(url.Values{}, "")
		if err == nil {
			t.Error("Test Failed - liqui GetTradeHistory() error", err)
		}

		_, err = l.WithdrawCoins("btc", 1337, "someaddr")
		if err == nil {
			t.Error("Test Failed - liqui WithdrawCoins() error", err)
		}
	}
}

func TestUpdateTicker(t *testing.T) {
	p := pair.NewCurrencyPairDelimiter("ETH_BTC", "_")
	_, err := l.UpdateTicker(p, "SPOT")
	if err != nil {
		t.Error("Test Failed - liqui UpdateTicker() error", err)
	}
}

func TestUpdateOrderbook(t *testing.T) {
	p := pair.NewCurrencyPairDelimiter("ETH_BTC", "_")
	_, err := l.UpdateOrderbook(p, "SPOT")
	if err != nil {
		t.Error("Test Failed - liqui UpdateOrderbook() error", err)
	}
}

func setFeeBuilder() exchange.FeeBuilder {
	return exchange.FeeBuilder{
		Amount:              1,
		Delimiter:           "-",
		FeeType:             exchange.CryptocurrencyTradeFee,
		FirstCurrency:       symbol.LTC,
		SecondCurrency:      symbol.BTC,
		IsMaker:             false,
		PurchasePrice:       1,
		CurrencyItem:        symbol.USD,
		BankTransactionType: exchange.WireTransfer,
	}
}
func TestGetFee(t *testing.T) {
	l.SetDefaults()
	var feeBuilder = setFeeBuilder()
	// CryptocurrencyTradeFee Basic
	if resp, err := l.GetFee(feeBuilder); resp != float64(0.0025) || err != nil {
		t.Error(err)
		t.Errorf("Test Failed - GetFee() error. Expected: %f, Recieved: %f", float64(0.0025), resp)
	}

	// CryptocurrencyTradeFee High quantity
	feeBuilder = setFeeBuilder()
	feeBuilder.Amount = 1000
	feeBuilder.PurchasePrice = 1000
	if resp, err := l.GetFee(feeBuilder); resp != float64(2500) || err != nil {
		t.Errorf("Test Failed - GetFee() error. Expected: %f, Recieved: %f", float64(2000), resp)
		t.Error(err)
	}

	// CryptocurrencyTradeFee IsMaker
	feeBuilder = setFeeBuilder()
	feeBuilder.IsMaker = true
	if resp, err := l.GetFee(feeBuilder); resp != float64(0.001) || err != nil {
		t.Errorf("Test Failed - GetFee() error. Expected: %f, Recieved: %f", float64(0.001), resp)
		t.Error(err)
	}

	// CryptocurrencyTradeFee Negative purchase price
	feeBuilder = setFeeBuilder()
	feeBuilder.PurchasePrice = -1000
	if resp, err := l.GetFee(feeBuilder); resp != float64(0) || err != nil {
		t.Errorf("Test Failed - GetFee() error. Expected: %f, Recieved: %f", float64(0), resp)
		t.Error(err)
	}
	// CryptocurrencyWithdrawalFee Basic
	feeBuilder = setFeeBuilder()
	feeBuilder.FeeType = exchange.CryptocurrencyWithdrawalFee
	if resp, err := l.GetFee(feeBuilder); resp != float64(0.01) || err != nil {
		t.Errorf("Test Failed - GetFee() error. Expected: %f, Recieved: %f", float64(0.01), resp)
		t.Error(err)
	}

	// CryptocurrencyWithdrawalFee Invalid currency
	feeBuilder = setFeeBuilder()
	feeBuilder.FirstCurrency = "hello"
	feeBuilder.FeeType = exchange.CryptocurrencyWithdrawalFee
	if resp, err := l.GetFee(feeBuilder); resp != float64(0) || err != nil {
		t.Errorf("Test Failed - GetFee() error. Expected: %f, Recieved: %f", float64(0), resp)
		t.Error(err)
	}

	// CyptocurrencyDepositFee Basic
	feeBuilder = setFeeBuilder()
	feeBuilder.FeeType = exchange.CyptocurrencyDepositFee
	if resp, err := l.GetFee(feeBuilder); resp != float64(0) || err != nil {
		t.Errorf("Test Failed - GetFee() error. Expected: %f, Recieved: %f", float64(0), resp)
		t.Error(err)
	}

	// InternationalBankDepositFee Basic
	feeBuilder = setFeeBuilder()
	feeBuilder.FeeType = exchange.InternationalBankDepositFee
	if resp, err := l.GetFee(feeBuilder); resp != float64(0) || err != nil {
		t.Errorf("Test Failed - GetFee() error. Expected: %f, Recieved: %f", float64(0), resp)
		t.Error(err)
	}

	// InternationalBankWithdrawalFee Basic
	feeBuilder = setFeeBuilder()
	feeBuilder.FeeType = exchange.InternationalBankWithdrawalFee
	feeBuilder.CurrencyItem = symbol.USD
	if resp, err := l.GetFee(feeBuilder); resp != float64(0) || err != nil {
		t.Errorf("Test Failed - GetFee() error. Expected: %f, Recieved: %f", float64(0), resp)
		t.Error(err)
	}
}

func TestFormatWithdrawPermissions(t *testing.T) {
	// Arrange
	l.SetDefaults()
	expectedResult := exchange.NoAPIWithdrawalMethodsText
	// Act
	withdrawPermissions := l.FormatWithdrawPermissions()
	// Assert
	if withdrawPermissions != expectedResult {
		t.Errorf("Expected: %s, Recieved: %s", expectedResult, withdrawPermissions)
	}
}

// This will really really use the API to place an order
// If you're going to test this, make sure you're willing to place real orders on the exchange
func TestSubmitOrder(t *testing.T) {
	l.SetDefaults()
	TestSetup(t)
	l.Verbose = true

	if l.APIKey == "" || l.APISecret == "" ||
		l.APIKey == "Key" || l.APISecret == "Secret" ||
		!canPlaceOrders {
		t.Skip(fmt.Sprintf("ApiKey: %s. Can place orders: %v", l.APIKey, canPlaceOrders))
	}
	var p = pair.CurrencyPair{
		Delimiter:      "",
		FirstCurrency:  symbol.BTC,
		SecondCurrency: symbol.EUR,
	}
	response, err := l.SubmitOrder(p, exchange.Buy, exchange.Market, 1, 10, "hi")
	if err != nil || !response.IsOrderPlaced {
		t.Errorf("Order failed to be placed: %v", err)
	}
}
