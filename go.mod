module github.com/bhbosman/goMarketData

go 1.18

require (
	github.com/bhbosman/goCommonMarketData v0.0.0-20230328221032-1140368f4836
	github.com/bhbosman/goCommsDefinitions v0.0.0-20230320101758-971a50fdbf8c
	github.com/bhbosman/goCommsNetDialer v0.0.0-20230328221032-2077fce8d8d6
	github.com/bhbosman/goCommsNetListener v0.0.0-20230328221032-95cc426246fe
	github.com/bhbosman/goCommsStacks v0.0.0-20230328221032-cd6c6063e9ef
	github.com/bhbosman/goFxApp v0.0.0-20230328222821-d6fbe28ebdf1
	github.com/bhbosman/goFxAppManager v0.0.0-20230328221825-d57a6eef05c0
	github.com/bhbosman/gocommon v0.0.0-20230328220050-dafaab862dd2
	github.com/bhbosman/gocomms v0.0.0-20230328221033-4231cc0340d0
	github.com/bhbosman/goprotoextra v0.0.2
	github.com/cskr/pubsub v1.0.2
	github.com/reactivex/rxgo/v2 v2.5.0
	go.uber.org/fx v1.19.2
	go.uber.org/zap v1.24.0
	golang.org/x/net v0.0.0-20220624214902-1bab6f366d9e
	golang.org/x/oauth2 v0.0.0-20220622183110-fd043fe589d2
	google.golang.org/api v0.90.0
	google.golang.org/protobuf v1.28.0
)

require github.com/bhbosman/goConn v0.0.0-20230328145648-e6bab0984275

require (
	cloud.google.com/go/compute v1.7.0 // indirect
	github.com/bhbosman/goConnectionManager v0.0.0-20230328221032-2e48eb2721d0 // indirect
	github.com/bhbosman/goMessages v0.0.0-20230328221033-d00082868706 // indirect
	github.com/bhbosman/goUi v0.0.0-20230328221033-53c557da423d // indirect
	github.com/bhbosman/goerrors v0.0.0-20220623084908-4d7bbcd178cf // indirect
	github.com/bhbosman/gomessageblock v0.0.0-20230308173223-e8144f25444c // indirect
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/gdamore/tcell/v2 v2.5.1 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.1.0 // indirect
	github.com/googleapis/gax-go/v2 v2.4.0 // indirect
	github.com/icza/gox v0.0.0-20220321141217-e2d488ab2fbc // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/tview v0.0.0-20220709181631-73bf2902b59a // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	github.com/teivah/onecontext v0.0.0-20200513185103-40f981bfd775 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/dig v1.16.1 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f // indirect
	golang.org/x/sys v0.0.0-20220624220833-87e55d714810 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220624142145-8cd45d7dbd1f // indirect
	google.golang.org/grpc v1.47.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/gdamore/tcell/v2 => github.com/bhbosman/tcell/v2 v2.5.2-0.20220624055704-f9a9454fab5b

replace github.com/golang/mock => github.com/bhbosman/gomock v1.6.1-0.20230302060806-d02c40b7514e

replace github.com/cskr/pubsub => github.com/bhbosman/pubsub v1.0.3-0.20220802200819-029949e8a8af

replace github.com/rivo/tview => github.com/bhbosman/tview v0.0.0-20230310100135-f8b257a85d36

//replace github.com/bhbosman/goCommsStacks => ../goCommsStacks

//replace github.com/bhbosman/goCommsNetDialer => ../goCommsNetDialer

//replace github.com/bhbosman/goCommsDefinitions => ../goCommsDefinitions

replace github.com/bhbosman/goFxApp => ../goFxApp

//replace github.com/bhbosman/goUi => ../goUi

//replace github.com/bhbosman/goerrors => ../goerrors

//replace github.com/bhbosman/goFxAppManager => ../goFxAppManager

//replace github.com/bhbosman/goConnectionManager => ../goConnectionManager

//replace github.com/bhbosman/goprotoextra => ../goprotoextra

replace github.com/bhbosman/goCommonMarketData => ../goCommonMarketData

replace github.com/bhbosman/goCommsNetListener => ../goCommsNetListener

//replace github.com/reactivex/rxgo/v2 => ../goRx

//replace github.com/bhbosman/goConn => ../goConn
