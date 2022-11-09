package liquid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"

	"github.com/shopspring/decimal"
)

type Orders struct {
	Models      []Order `json:"models"`
	CurrentPage int     `json:"current_page"`
	TotalPages  int     `json:"total_pages"`
}

type Order struct {
	ID                        int64           `json:"id,omitempty"`
	OrderType                 string          `json:"order_type,omitempty"`
	Quantity                  decimal.Decimal `json:"quantity,omitempty"`
	DiscQuantity              decimal.Decimal `json:"disc_quantity,omitempty"`
	IcebergTotalQuantity      decimal.Decimal `json:"iceberg_total_quantity,omitempty"`
	Side                      string          `json:"side,omitempty"`
	FilledQuantity            decimal.Decimal `json:"filled_quantity,omitempty"`
	Price                     decimal.Decimal `json:"price,omitempty"`
	CreatedAt                 int             `json:"created_at,omitempty"`
	UpdatedAt                 int             `json:"updated_at,omitempty"`
	Status                    string          `json:"status,omitempty"`
	LeverageLevel             int             `json:"leverage_level,omitempty"`
	SourceExchange            interface{}     `json:"source_exchange,omitempty"`
	ProductID                 int             `json:"product_id,omitempty"`
	MarginType                interface{}     `json:"margin_type,omitempty"`
	TakeProfit                interface{}     `json:"take_profit,omitempty"`
	StopLoss                  interface{}     `json:"stop_loss,omitempty"`
	TradingType               string          `json:"trading_type,omitempty"`
	ProductCode               string          `json:"product_code,omitempty"`
	FundingCurrency           string          `json:"funding_currency,omitempty"`
	CryptoAccountID           interface{}     `json:"crypto_account_id,omitempty"`
	CurrencyPairCode          string          `json:"currency_pair_code,omitempty"`
	AveragePrice              decimal.Decimal `json:"average_price,omitempty"`
	Target                    string          `json:"target,omitempty"`
	OrderFee                  decimal.Decimal `json:"order_fee,omitempty"`
	SourceAction              string          `json:"source_action,omitempty"`
	UnwoundTradeID            interface{}     `json:"unwound_trade_id,omitempty"`
	TradeID                   interface{}     `json:"trade_id,omitempty"`
	ClientOrderID             string          `json:"client_order_id,omitempty"`
	Settings                  interface{}     `json:"settings,omitempty"`
	TrailingStopType          interface{}     `json:"trailing_stop_type,omitempty"`
	TrailingStopValue         interface{}     `json:"trailing_stop_value,omitempty"`
	StopTriggeredTime         interface{}     `json:"stop_triggered_time,omitempty"`
	MarginUsed                decimal.Decimal `json:"margin_used,omitempty"`
	MarginInterest            decimal.Decimal `json:"margin_interest,omitempty"`
	UnwoundTradeLeverageLevel interface{}     `json:"unwound_trade_leverage_level,omitempty"`
}

type OrderExecutions []struct {
	ID        int     `json:"id"`
	Quantity  float64 `json:"quantity,string"`
	Price     float64 `json:"price,string"`
	TakerSide string  `json:"taker_side"`
	MySide    string  `json:"my_side"`
	CreatedAt int64   `json:"created_at"`
}

func (c *Client) GetOrder(orderID int) (Order, error) {
	spath := fmt.Sprintf("/orders/%d", orderID)

	var order Order
	res, err := c.sendRequest("GET", spath, nil, nil)
	if err != nil {
		return order, err
	}

	if err := decode(res, &order); err != nil {
		return order, err
	}

	return order, nil
}

func (c *Client) GetOrderByClientID(clientID string) (Order, error) {
	spath := fmt.Sprintf("/orders/client:%s", clientID)

	var order Order
	res, err := c.sendRequest("GET", spath, nil, nil)
	if err != nil {
		return order, err
	}

	if err := decode(res, &order); err != nil {
		return order, err
	}

	return order, nil
}

type OrdersFilter struct {
	ProductID       string `json:"product_id,omitempty"`
	WithDetails     string `json:"with_details,omitempty"`
	Status          string `json:"status,omitempty"`
	FundingCurrency string `json:"funding_currency,omitempty"`

	// 下記Doc非公開フィルター
	// WebAPI: page=1&limit=24&currency_pair_code=BTCJPY&status=live&trading_type=cfd
	Page             string `json:"page,omitempty"`
	Limit            string `json:"limit,omitempty"`
	CurrencyPairCode string `json:"currency_pair_code,omitempty"`
	TradingType      string `json:"trading_type,omitempty"`
}

func (c *Client) GetOrders(filters OrdersFilter) (Order, error) {
	var orders Order

	j, err := json.Marshal(filters)
	if err != nil {
		return orders, err
	}

	res, err := c.sendRequest("GET", "/orders", bytes.NewReader(j), nil)
	if err != nil {
		return orders, err
	}

	if err := decode(res, &orders); err != nil {
		return orders, err
	}

	return orders, nil
}

type RequestOrder struct {
	Order OrderParams `json:"order"`
}

type OrderParams struct {
	OrderType  string `json:"order_type"` // limit, market or market_with_range
	ProductID  int    `json:"product_id"`
	Side       string `json:"side"`
	Quantity   string `json:"quantity"`
	Price      string `json:"price,omitempty"`
	PriceRange string `json:"price_range,omitempty"` // order type optional
	// Margin trade
	LeverageLevel   int    `json:"leverage_level,omitempty"`
	FundingCurrency string `json:"funding_currency,omitempty"`
	OrderDirection  string `json:"order_direction,omitempty"`
	// CFD Optionals
	TradingType string `json:"trading_type,omitempty"` // margin or cfd, only available if leverage_level > 1
	MarginType  string `json:"margin_type,omitempty"`  // cross or isolated, only available if leverage_level > 1, default is cross

	// ClientID
	ClientOrderID string `json:"client_order_id,omitempty"`
}

// orderType, side, quantity, price, priceRange string, productID int
func (c *Client) CreateOrder(o *RequestOrder) (Order, error) {
	var order Order

	body, err := json.Marshal(o)
	if err != nil {
		return order, err
	}

	res, err := c.sendRequest("POST", "/orders/", bytes.NewReader(body), nil)
	if err != nil {
		return order, err
	}

	if err := decode(res, &order); err != nil {
		return order, err
	}

	return order, nil
}

func (c *Client) CancelOrder(orderID int) (Order, error) {
	spath := fmt.Sprintf("/orders/%d/cancel", orderID)

	var order Order
	res, err := c.sendRequest("PUT", spath, nil, nil)
	if err != nil {
		return order, err
	}

	if err := decode(res, &order); err != nil {
		return order, err
	}

	return order, nil
}

type EditOrderParams struct {
	Quantity string `json:"quantity"`
	Price    string `json:"price"`
}

// func (c *Client) EditLiveOrder(orderID int, quantity, price string) (Order, error) {
func (c *Client) EditLiveOrder(id int, e *EditOrderParams) (Order, error) {
	spath := fmt.Sprintf("/orders/%d", id)
	// bodyTemplate :=
	// 	`{
	// 		"order": { // Orderセクションは実はいらない
	// 			"quantity":	"%s",
	// 			"price":	"%s"
	// 		}
	// 	}`
	// body := fmt.Sprintf(bodyTemplate, quantity, price)

	var order Order

	body, err := json.Marshal(e)
	if err != nil {
		return order, err
	}

	res, err := c.sendRequest("PUT", spath, bytes.NewReader(body), nil)
	if err != nil {
		return order, err
	}

	if err := decode(res, &order); err != nil {
		return order, err
	}

	return order, nil
}

// https://api.liquid.com/trades/close_all?funding_currency=JPY&product_id=5
// OrderAllClose
func (p *Client) OrderAllClose() {

}

// ToPriceString is float to string price.00001
func ToPriceString(price float64) string {
	price = math.RoundToEven(price*100000) / 100000

	return fmt.Sprintf("%.5f", price)
}

// ToQtyString is float to string size.001
func ToQtyString(size float64) string {
	size = math.RoundToEven(size*1000) / 1000

	return fmt.Sprintf("%.3f", size)
}
