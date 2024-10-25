module voolibow_gw

go 1.22.5

require (
	github.com/gofiber/contrib/websocket v1.3.0
	github.com/gofiber/fiber/v2 v2.52.0
	google.golang.org/grpc v1.61.0
	google.golang.org/protobuf v1.31.0
	safir/libs/appconfigs v0.0.0-00010101000000-000000000000
	safir/libs/appstates v0.0.0-00010101000000-000000000000
)

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/fasthttp/websocket v1.5.7 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/klauspost/compress v1.17.3 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/savsgio/gotils v0.0.0-20230208104028-c358bd845dee // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.51.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231106174013-bbf56f31fb17 // indirect
)

replace safir/libs/appconfigs => ../../libs/appconfigs

replace safir/libs/appstates => ../../libs/appstates
