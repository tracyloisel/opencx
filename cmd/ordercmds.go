package main

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/mit-dci/opencx/logging"

	"github.com/mit-dci/opencx/cxrpc"
	"github.com/mit-dci/opencx/match"

	"github.com/olekukonko/tablewriter"
)

// OrderCommand submits an order (for now) TODO
func (cl *openCxClient) OrderCommand(args []string) error {
	orderArgs := new(cxrpc.SubmitOrderArgs)
	orderReply := new(cxrpc.SubmitOrderReply)

	var newOrder match.LimitOrder
	newOrder.Client = args[0]
	newOrder.Side = args[1]
	if newOrder.Side != "buy" && newOrder.Side != "sell" {
		return fmt.Errorf("Order's side isn't buy or sell, try again")
	}

	// get the trading pair string from the shell input - third parameter
	err := newOrder.TradingPair.FromString(args[2])
	if err != nil {
		return fmt.Errorf("Error getting asset pair from string: \n%s", err)
	}

	newOrder.AmountHave, err = strconv.ParseUint(args[3], 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing amountHave, please enter something valid:\n%s", err)
	}

	price, err := strconv.ParseFloat(args[4], 64)
	if err != nil {
		return fmt.Errorf("Error parsing price: \n%s", err)
	}
	newOrder.SetAmountWant(price)

	orderArgs.Order = &newOrder
	err = cl.Call("OpencxRPC.SubmitOrder", orderArgs, orderReply)
	if err != nil {
		return fmt.Errorf("Error calling 'SubmitOrder' service method:\n%s", err)
	}

	logging.Infof("Order submitted successfully\n")
	return nil
}

func (cl *openCxClient) GetPrice(args []string) error {
	var err error

	getPriceArgs := new(cxrpc.GetPriceArgs)
	getPriceReply := new(cxrpc.GetPriceReply)

	// can't be a nil pointer to call methods on it
	getPriceArgs.TradingPair = new(match.Pair)

	// get the trading pair string from the shell input - first parameter
	if err = getPriceArgs.TradingPair.FromString(args[0]); err != nil {
		return err
	}

	if err = cl.Call("OpencxRPC.GetPrice", getPriceArgs, getPriceReply); err != nil {
		return err
	}

	logging.Infof("Price: %f\n", getPriceReply.Price)
	return nil
}

// ViewOrderbook return s the orderbook TODO
func (cl *openCxClient) ViewOrderbook(args []string) error {
	var err error

	viewOrderBookArgs := new(cxrpc.ViewOrderBookArgs)
	viewOrderBookReply := new(cxrpc.ViewOrderBookReply)

	// can't be a nil pointer to call methods on it
	viewOrderBookArgs.TradingPair = new(match.Pair)

	// get the trading pair string from the shell input - first parameter
	err = viewOrderBookArgs.TradingPair.FromString(args[0])
	if err != nil {
		return err
	}

	err = cl.Call("OpencxRPC.ViewOrderBook", viewOrderBookArgs, viewOrderBookReply)
	if err != nil {
		return fmt.Errorf("Error calling 'ViewOrderBook' service method:\n%s", err)
	}

	if len(args) == 1 {
		var data [][]string
		buf := new(bytes.Buffer)
		table := tablewriter.NewWriter(buf)
		table.SetHeader([]string{"orderID", "price", "volume", "side"})

		for _, buyOrder := range viewOrderBookReply.BuyOrderBook {
			buyPrice, err := buyOrder.Price()
			if err != nil {
				return err
			}

			// convert stuff to strings
			strPrice := fmt.Sprintf("%f", buyPrice)
			strVolume := fmt.Sprintf("%d", buyOrder.AmountHave)
			// append to the table
			data = append(data, []string{buyOrder.OrderID, strPrice, strVolume, buyOrder.Side})
		}

		for _, sellOrder := range viewOrderBookReply.SellOrderBook {
			sellPrice, err := sellOrder.Price()
			if err != nil {
				return err
			}

			// convert stuff to strings
			strPrice := fmt.Sprintf("%f", sellPrice)
			strVolume := fmt.Sprintf("%d", sellOrder.AmountHave)
			// append to the table
			data = append(data, []string{sellOrder.OrderID, strPrice, strVolume, sellOrder.Side})
		}

		table.AppendBulk(data)
		table.Render()

		// actually print out table stored in buffer
		logging.Infof("\n%s\n", buf.String())
		return nil
	}

	if len(args) == 2 && args[1] == "sell" {
		var data [][]string
		buf := new(bytes.Buffer)
		table := tablewriter.NewWriter(buf)
		table.SetHeader([]string{"orderID", "price", "volume", "side"})

		for _, sellOrder := range viewOrderBookReply.SellOrderBook {
			sellPrice, err := sellOrder.Price()
			if err != nil {
				return err
			}

			// convert stuff to strings
			strPrice := fmt.Sprintf("%f", sellPrice)
			strVolume := fmt.Sprintf("%d", sellOrder.AmountHave)
			// append to the table
			data = append(data, []string{sellOrder.OrderID, strPrice, strVolume, sellOrder.Side})
		}

		table.AppendBulk(data)
		table.Render()

		// actually print out table stored in buffer
		logging.Infof("\n%s\n", buf.String())
		return nil
	} else if len(args) == 2 && args[1] == "buy" {
		var data [][]string
		buf := new(bytes.Buffer)
		table := tablewriter.NewWriter(buf)
		table.SetHeader([]string{"orderID", "price", "volume", "side"})

		for _, buyOrder := range viewOrderBookReply.BuyOrderBook {
			buyPrice, err := buyOrder.Price()
			if err != nil {
				return err
			}

			// convert stuff to strings
			strPrice := fmt.Sprintf("%f", buyPrice)
			strVolume := fmt.Sprintf("%d", buyOrder.AmountHave)
			// append to the table
			data = append(data, []string{buyOrder.OrderID, strPrice, strVolume, buyOrder.Side})
		}

		table.AppendBulk(data)
		table.Render()

		// actually print out table stored in buffer
		logging.Infof("\n%s\n", buf.String())
		return nil
	}

	logging.Warnf("Something went wrong! But I'm not going to quit because this is just the client. I'm lost! I don't know how I got here! Help!")

	return nil
}