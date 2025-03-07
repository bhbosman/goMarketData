package main

import (
	"fmt"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerViewer"
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
	"github.com/bhbosman/goFxApp"
	"github.com/bhbosman/goMarketData/internal/incomingMarketData"
	"github.com/bhbosman/goMarketData/internal/listener"
	app2 "github.com/bhbosman/gocommon/Providers"
	"go.uber.org/fx"
	"net/url"
	"os"
)

func main() {
	u, _ := url.Parse("tcp4://127.0.0.1:4001")
	var runTimeManager *app2.RunTimeManager
	app := goFxApp.NewFxMainApplicationServices(
		"MarketData",
		false,
		fullMarketDataManagerViewer.Provide(),
		fullMarketDataManagerService.Provide(true),
		fullMarketDataHelper.Provide(),
		instrumentReference.Provide(),
		//googleSheets.Provide(),
		listener.CompressedListener(1024, false, nil, u),
		app2.RegisterRunTimeManager(),
		incomingMarketData.ProvideLunoMarketDataDialer(1, "tcp4://127.0.0.1:3001"),
		incomingMarketData.ProvideKrakenMarketDataDialer(1, "tcp4://127.0.0.1:3011"),
		fx.Populate(&runTimeManager),
	)

	if app.FxApp.Err() != nil {
		_, _ = fmt.Fprint(os.Stderr, app.FxApp.Err().Error())
		return
	}
	app.RunTerminalApp()
}
