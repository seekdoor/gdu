const DONE = "done";
const SCANNING = "scanning";

function formatSize(size) {
    round = (x, y) => Math.round(x / Math.pow(2, y) * 10) / 10

    if (size > 1e12) {
        return `${round(size, 40)} TiB`;
    } else if (size > 1e9) {
        return `${round(size, 30)} GiB`;
    } else if (size > 1e9) {
        return `${round(size, 20)} MiB`;
    } else if (size > 1e9) {
        return `${round(size, 10)} KiB`;
    } else {
        return `${size} B`;
    }
}

function connectWs() {
    let ws = new WebSocket('ws://localhost:8888/ws');
    let status;

    ws.addEventListener('open', function (event) {
        console.log("open");
        const statusElem = document.querySelector("#status");
        statusElem.innerHTML = "";
    });
    ws.addEventListener('error', function (event) {
        console.log("error");
        ws.close();
    });
    ws.addEventListener('close', function (event) {
        const statusElem = document.querySelector("#status");
        statusElem.innerHTML = '"Offline"';
        statusElem.style = "color: red";
        setTimeout(() => connectWs(), 1000);
    });

    ws.addEventListener('message', (event) => {
        const statusElem = document.querySelector("#status");
        const data = JSON.parse(event.data);

        if (data.MsgType === "progress") {
            const progressElem = document.querySelector("#progress");

            if (data.Done === false && status !== SCANNING) {
                status = SCANNING;
                statusElem.innerHTML = '"Scanning..."';
                statusElem.style = "color: green";
                progressElem.style = "display: flex";
            }
            if (data.Done === true && status !== DONE) {
                status = DONE;
                statusElem.innerHTML = '"Done"';
                statusElem.style = "color: green";
                progressElem.style = "display: none";

                ws.send(JSON.stringify({MsgType: "command"}));
            }

            progressElem.innerHTML = `<span class="center">
                Total items: <span class="red">${data.ItemCount}</span><br />
                Size: <span class="red">${formatSize(data.TotalSize)}</span></span>
            `;
        }
        console.log('Message from server ', data);
    });
}

connectWs();

