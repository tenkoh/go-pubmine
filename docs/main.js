const workers = [];
const threads = navigator.hardwareConcurrency;
const maxWorkers = threads > 1 ? threads - 1 : threads;
for (let i = 0; i < maxWorkers; i++) {
  const w = new Worker("mine.js");
  workers.push(w);
}
const mine = () => {
  const btn = document.getElementById("run");
  btn.disabled = true;
  const prefix = document.getElementById("input").value;
  const jobs = [];
  for (const w of workers) {
    const promise = new Promise((resolve) => {
      w.onmessage = (e) => {
        resolve(e.data);
      };
    });
    jobs.push(promise);
    w.postMessage(prefix);
  }
  Promise.any(jobs).then((result) => {
    console.log(result);
    document.getElementById("public").value = result.public;
    document.getElementById("private").value = result.private;
    for (let [i, w] of workers.entries()) {
      w.terminate();
      w = null;
      const nw = new Worker("mine.js");
      workers[i] = nw;
    }
    btn.disabled = false;
  });
};

const copyToClip = async (id) => {
  const elem = document.getElementById(id);
  try {
    await navigator.clipboard.writeText(elem.value);
  } catch (err) {
    alert("Failed to copy: ", err);
  }
};
