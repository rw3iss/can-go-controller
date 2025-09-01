

const startApp = async (cb: () => void) => {
	console.log(`Starting app...`)
	await loadWasm();
	cb();
}

const loadWasm = async () => {
	return new Promise((resolve, reject) => {
		try {
			if (!navigator?.serviceWorker) throw "service worker not supported.";
			console.log(`Loading WASM service worker...`)
			navigator.serviceWorker.register('sw.js').then(registration => {
				const sw = registration.active ?? registration.installing ?? registration.waiting;
				console.log(`WASM Loaded.`)
				return resolve(true);
				//				if (sw) setInterval(() => sw.postMessage(null), 15000)
			})

			// const go = new Go();
			// WebAssembly.instantiateStreaming(fetch("/CANServer.wasm"), go.importObject)
			// 	.then((result) => {
			// 		go.run(result.instance);
			// 		resolve(true);
			// 	})
			// 	.catch((e) => reject(e));
		} catch (e) {
			reject(e);
		}
	});
}