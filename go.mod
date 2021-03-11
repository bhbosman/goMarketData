module github.com/bhbosman/goMarketData

go 1.15

require (
	github.com/bhbosman/goMessages v0.0.0-20201004192822-66d168b4b744
	github.com/bhbosman/gocommon v0.0.0-20201004145117-eae02ab42c9a
	github.com/bhbosman/gocomms v0.0.0-20201004142558-41c4e9c3302c
	github.com/bhbosman/gologging v0.0.0-20200921180328-d29fc55c00bc
	github.com/bhbosman/gomessageblock v0.0.0-20200921180725-7cd29a998aa3
	github.com/bhbosman/goprotoextra v0.0.15
	github.com/cskr/pubsub v1.0.2
	github.com/reactivex/rxgo/v2 v2.1.0
	go.uber.org/fx v1.13.1
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	google.golang.org/api v0.32.0
)

replace (
	github.com/bhbosman/goMessages => ../goMessages
	github.com/bhbosman/gocommon => ../gocommon
	github.com/bhbosman/gocomms => ../gocomms
	github.com/bhbosman/gologging => ../gologging
	github.com/bhbosman/gomessageblock => ../gomessageblock
	github.com/bhbosman/goprotoextra => ../goprotoextra
	github.com/reactivex/rxgo/v2 => ../../reactivex/rxgo
)
