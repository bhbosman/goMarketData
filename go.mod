module github.com/bhbosman/goMarketData

go 1.23.0

toolchain go1.24.0

require (
	github.com/bhbosman/goCommonMarketData v0.0.0-20230730212407-5a7831da9c11
	github.com/bhbosman/goCommsDefinitions v0.0.0-20250307125227-bfc368fdddf9
	github.com/bhbosman/goCommsNetDialer v0.0.0-20250307233555-6c2dfa80f01b
	github.com/bhbosman/goCommsNetListener v0.0.0-20250307153216-6206fd2748ea
	github.com/bhbosman/goCommsStacks v0.0.0-20231011182118-47d6d38b38e4
	github.com/bhbosman/goConn v0.0.0-20250307232512-5c59bbebe521
	github.com/bhbosman/goFxApp v0.0.0-20250307153150-937959817ddb
	github.com/bhbosman/goFxAppManager v0.0.0-20250307145515-bda0fa4d9959
	github.com/bhbosman/gocommon v0.0.0-20250307145810-dee67af89371
	github.com/bhbosman/gocomms v0.0.0-20230730212408-04ba72ddb372
	github.com/bhbosman/goprotoextra v0.0.2
	github.com/cskr/pubsub v1.0.2
	github.com/reactivex/rxgo/v2 v2.5.0
	go.uber.org/fx v1.23.0
	go.uber.org/zap v1.27.0
	golang.org/x/net v0.37.0
	golang.org/x/oauth2 v0.28.0
	google.golang.org/api v0.224.0
	google.golang.org/protobuf v1.36.5
)

require (
	cloud.google.com/go/compute v1.34.0 // indirect
	github.com/bhbosman/goConnectionManager v0.0.0-20250307224538-a79ceb218fd0 // indirect
	github.com/bhbosman/goMessages v0.0.0-20250307224348-83ddb4c19467 // indirect
	github.com/bhbosman/goUi v0.0.0-20250307150712-d06325af4877 // indirect
	github.com/bhbosman/goerrors v0.0.0-20250307194237-312d070c8e38 // indirect
	github.com/bhbosman/gomessageblock v0.0.0-20250307141417-ab783e8e2eba // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/gdamore/encoding v1.0.1 // indirect
	github.com/gdamore/tcell/v2 v2.8.1 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.5 // indirect
	github.com/googleapis/gax-go/v2 v2.14.1 // indirect
	github.com/icza/gox v0.2.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/tview v0.0.0-20241227133733-17b7edb88c57 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/teivah/onecontext v1.3.0 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/dig v1.18.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20250122153221-138b5a5a4fd4 // indirect
	google.golang.org/grpc v1.70.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/gdamore/tcell/v2 => github.com/bhbosman/tcell/v2 v2.5.2-0.20220624055704-f9a9454fab5b

replace github.com/golang/mock => github.com/bhbosman/gomock v1.6.1-0.20230302060806-d02c40b7514e

replace github.com/cskr/pubsub => github.com/bhbosman/pubsub v1.0.3-0.20220802200819-029949e8a8af

replace github.com/rivo/tview => github.com/bhbosman/tview v0.0.0-20230310100135-f8b257a85d36
