module github.com/bhbosman/goMarketData

go 1.23.0

toolchain go1.24.0

require (
	github.com/bhbosman/goCommonMarketData v0.0.0-20250308054451-6dade8c5c3f9
	github.com/bhbosman/goCommsDefinitions v0.0.0-20250308000247-4306925b3dfd
	github.com/bhbosman/goCommsNetDialer v0.0.0-20250308053613-a763632f1b3d
	github.com/bhbosman/goCommsNetListener v0.0.0-20250308054504-9e394f3cd581
	github.com/bhbosman/goCommsStacks v0.0.0-20250308055233-8992b9d0c219
	github.com/bhbosman/goConn v0.0.0-20250307235008-177f4ffe3521
	github.com/bhbosman/goFxApp v0.0.0-20250308052840-c690d60e479d
	github.com/bhbosman/goFxAppManager v0.0.0-20250307225418-ef314d0a9319
	github.com/bhbosman/gocommon v0.0.0-20250308052839-0ebeb121f996
	github.com/bhbosman/gocomms v0.0.0-20250308053431-cfd3e830d8a5
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
	cloud.google.com/go/auth v0.15.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.7 // indirect
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	github.com/bhbosman/goConnectionManager v0.0.0-20250307224538-a79ceb218fd0 // indirect
	github.com/bhbosman/goMessages v0.0.0-20250308052839-b11fc966f96a // indirect
	github.com/bhbosman/goUi v0.0.0-20250308052840-a0e5fd7e5f88 // indirect
	github.com/bhbosman/goerrors v0.0.0-20250307194237-312d070c8e38 // indirect
	github.com/bhbosman/gomessageblock v0.0.0-20250307141417-ab783e8e2eba // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/gdamore/encoding v1.0.1 // indirect
	github.com/gdamore/tcell/v2 v2.8.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
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
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.59.0 // indirect
	go.opentelemetry.io/otel v1.34.0 // indirect
	go.opentelemetry.io/otel/metric v1.34.0 // indirect
	go.opentelemetry.io/otel/trace v1.34.0 // indirect
	go.uber.org/dig v1.18.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250303144028-a0af3efb3deb // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250303144028-a0af3efb3deb // indirect
	google.golang.org/grpc v1.71.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/gdamore/tcell/v2 => github.com/bhbosman/tcell/v2 v2.5.2-0.20220624055704-f9a9454fab5b

replace github.com/golang/mock => github.com/bhbosman/gomock v1.6.1-0.20230302060806-d02c40b7514e

replace github.com/cskr/pubsub => github.com/bhbosman/pubsub v1.0.3-0.20220802200819-029949e8a8af

replace github.com/rivo/tview => github.com/bhbosman/tview v0.0.0-20230310100135-f8b257a85d36
