GOOS=js GOARCH=wasm tinygo build -o build/wasm/CANServer.wasm main.go

cp build/wasm/CANServer.wasm ../can-ui/public
