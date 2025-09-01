(Note: This may work on other platforms besides linux, with real can devices, but the virtal device testing only works on linux)

## Setup:
To setup the can server listener and web api, first install Go, then:
```
cd server
go get
```

Install can-utils and start a virtual can device, if testing on linux:
```
sudo dnf install can-utils
sudo ip link add dev vcan0 type vcan
sudo ip link set vcan0 up
```

For the web app, install node, then:
```
cd app
npm install
npm run
```

## Run program:
To start the web api (port 8080) and can listener for the vcan0 device, run:
`go run main.go vcan0`

Send a test frame in another terminal and see it (if can-utils are installed):
`cansend vcan0 123#deadbeef`

See the events in the app dashboard if running: http://localhost:3000


## Other stuff (later):

build wasm:
`GOOS=js GOARCH=wasm go build -o build/wasm/CANServer.wasm main.go`

build with tinygo:
`GOOS=js GOARCH=wasm tinygo build -o build/wasm/CANServer.wasm main.go`

Copy wasm to frontend:
`cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" .`

## Check out

https://github.com/nlepage/go-wasm-http-server/tree/master

Canable Firmware: https://github.com/normaldotcom/canable2-fw
