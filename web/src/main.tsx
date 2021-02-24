import * as React from 'react'
import * as ReactDOM from 'react-dom'
import { Pie } from 'react-chartjs-2';

const DONE = "done";
const SCANNING = "scanning";

function formatSize(size: number): string {
    let round = (x: number, y: number) => Math.round(x / Math.pow(2, y) * 10) / 10

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

function connectWs(dataReceivedCallback: Function) {
    let ws = new WebSocket('ws://localhost:8888/ws');

    ws.addEventListener('open', function (event) {
        console.log("open");
    });
    ws.addEventListener('error', function (event) {
        console.log("error");
        ws.close();
    });
    ws.addEventListener('close', function (event) {
        dataReceivedCallback({MsgType: "close"});
        setTimeout(() => connectWs(dataReceivedCallback), 1000);
    });

    ws.addEventListener('message', (event) => {
        const data = JSON.parse(event.data);
        dataReceivedCallback(data);
    });
}

interface AppProps {}

interface AppState {
    status: string
    itemCount: number
    totalSize: number
}

interface DataRecord {
    MsgType: string
    Done: boolean
    ItemCount: number
    TotalSize: number
}

class App extends React.Component<AppProps, AppState> {
    constructor(props: AppProps) {
        super(props);
        this.state = {status: null, itemCount: 0, totalSize: 0};
        this.onMessage = this.onMessage.bind(this);
    }

    onMessage(data: DataRecord) {
        if (data.MsgType === "progress") {
            const status = this.state.status;

            if (data.Done === false && status !== SCANNING) {
                this.setState({status: SCANNING});
            }
            if (data.Done === true && status !== DONE) {
                this.setState({status: DONE});
            }

            this.setState({
                itemCount: data.ItemCount,
                totalSize: data.TotalSize,
            });
        }
    }


    componentDidMount() {
        connectWs(this.onMessage);
    }

    renderStatus() {
        if (![SCANNING, DONE].includes(this.state.status)) {
            return "";
        }

        if (this.state.status === SCANNING) {
            return <span className="green">"Scanning..."</span>;
        }
        if (this.state.status === DONE) {
            return <span className="green">"Done"</span>;
        }
    }

    renderProgress() {
        if (this.state.status !== SCANNING) {
            return "";
        }

        return (
            <div id="progress">
                <span className="center">
                    Total items: <span className="red">{this.state.itemCount}</span><br />
                    Size: <span className="red">{formatSize(this.state.totalSize)}</span>
                </span>
            </div>
        );
    }

    renderPie() {
        if (this.state.status !== DONE) {
            return "";
        }

        return (
            <Pie
                width={500}
                data={{
                    datasets: [{
                        data: [1, 2, 3],
                        backgroundColor: ["#f00", "#0f0", "#00f"],
                    }],
                }}
            />
        );

    }

    render() {
        return (
            <div>
                <h1>go DiskUsage({this.renderStatus()})</h1>
                {this.renderProgress()}
                {this.renderPie()}
            </div>
        );
    }
}

const domContainer = document.querySelector('#container');
ReactDOM.render(<App />, domContainer);