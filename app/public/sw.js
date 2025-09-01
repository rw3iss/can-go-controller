// importScripts('https://cdn.jsdelivr.net/gh/golang/go@go1.18.4/misc/wasm/wasm_exec.js')
// importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@v1.1.0/sw.js')

addEventListener('install', event => {
	event.waitUntil(skipWaiting())
	console.log(`SW install done.`)
})

addEventListener('activate', event => {
	event.waitUntil(clients.claim())
	console.log(`SW activate done.`)
})

addEventListener('message', (msg) => {
	console.log(`sw msg:`, msg)
})

loadApp = () => {
	const go = new Go();
	WebAssembly.instantiateStreaming(fetch("/CANServer.wasm"), go.importObject)
		.then((result) => {
			go.run(result.instance);
			//resolve(true);
		})
		.catch((e) => reject(e));
}