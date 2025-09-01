# Setup:
go get

Also install can-utils cli if sending test messages, ie:
sudo dnf install can-utils

# Start virtual can device:
sudo ip link add dev vcan0 type vcan
sudo ip link set vcan0 up

# Run program:
Start the web server/api on port 8080, and the can listener for vcan0 (or another):
go run main.go vcan0

Send a test frame in another terminal and see it:
cansend vcan0 123#deadbeef

Can also see the events on can-ui if running: http://localhost:3000


# OTHER STUFF (later): ====================================

## build wasm:
GOOS=js GOARCH=wasm go build -o build/wasm/CANServer.wasm main.go

## build with tinygo:
GOOS=js GOARCH=wasm tinygo build -o build/wasm/CANServer.wasm main.go

## Copy wasm to frontend:
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" .

Check out
https://github.com/nlepage/go-wasm-http-server/tree/master

## Canable Firmware:
https://github.com/normaldotcom/canable2-fw
