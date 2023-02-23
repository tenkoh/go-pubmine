importScripts("wasm_exec.js");
const go = new Go();
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then(
  (result) => {
    go.run(result.instance);
  },
);
onmessage = (e) => {
  const ret = mine(e.data);
  postMessage(ret);
};
