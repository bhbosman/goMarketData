module github.com/bhbosman/goMarketData

go 1.15

require (
	github.com/bhbosman/goMessages v0.0.0-20200922081308-8c8f88094624
	github.com/bhbosman/gocommon v0.0.0-20200921215456-bfddd9bb050e
	github.com/bhbosman/gocomms v0.0.0-20200922180846-e1a99ac26663
	github.com/bhbosman/gologging v0.0.0-20200921180328-d29fc55c00bc
	github.com/bhbosman/gomessageblock v0.0.0-20200921180725-7cd29a998aa3
	github.com/bhbosman/goprotoextra v0.0.1
	github.com/cskr/pubsub v1.0.2
	github.com/gobwas/httphead v0.0.0-20200921212729-da3d93bc3c58 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/reactivex/rxgo/v2 v2.1.0
	go.uber.org/fx v1.13.1
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	google.golang.org/api v0.32.0
)
replace github.com/reactivex/rxgo/v2 v2.1.0 => github.com/bhbosman/rxgo/v2 v2.1.1-0.20200922152528-6aef42e76e00
replace github.com/bhbosman/gocomms => /Users/brendanbosman/src/github.com/bhbosman/gocomms