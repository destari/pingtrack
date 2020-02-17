import React from "react"
import "rbx/index.css"
import {Container, Notification, Table} from "rbx";
import {ColumnChart, LineChart} from 'react-chartkick'
import 'chart.js'
import Axios from "axios";

function DataTable(props) {

    if (props.data) {
        let data = props.data;

        const dataTable = data.map(function (d) {
            return (
                <Table.Row key={d.Time}>
                    <Table.Cell>{d.Time}</Table.Cell>
                    <Table.Cell>{d.MinRtt / 1000000} ms</Table.Cell>
                    <Table.Cell>{d.AvgRtt / 1000000} ms</Table.Cell>
                    <Table.Cell>{d.MaxRtt / 1000000} ms</Table.Cell>
                    <Table.Cell>{d.PacketLoss} %</Table.Cell>
                </Table.Row>
            )
        });
        return dataTable;
    } else {
        return (
            <tr><td>Loading...</td><td>Loading...</td></tr>
        );
    }
}

function dataTwist(inputData) {
    let outputData = [];
    let packetLoss = [];

    if (inputData) {

        let series = {};
        let seriesMax = {};
        let seriesMin = {};
        let seriesLoss = {};
        inputData.forEach( function (datum) {
            series[datum.Time] = datum.AvgRtt/1000000;
            seriesMax[datum.Time] = datum.MaxRtt/1000000;
            seriesMin[datum.Time] = datum.MinRtt/1000000;
            seriesLoss[datum.Time] = datum.PacketLoss
            }
        );

        outputData.push({"name": "Avg RTT", "data": series});
        outputData.push({"name": "Min RTT", "data": seriesMin});
        outputData.push({"name": "Max RTT", "data": seriesMax});

        packetLoss.push({"name": "Packet Loss", "data": seriesLoss})

    }

    return [outputData, packetLoss]
}

class DataView extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            hostStats: [],
            data: [],
            packetLoss: [],
            hostname: props.host
        }
    }

    intervalID;
    Results;

    updateHostname = (newHostname) => {
        this.setState({hostStats: [], data: [], packetLoss: []});
        this.setState({hostname: newHostname});
        this.getData();
    };

    componentDidMount() {
        this.getData();
        this.intervalID = setInterval(this.getData.bind(this), 15000);
    }

    componentWillUnmount() {
        clearInterval(this.intervalID);
    }

    getData = () => {
        if (this.state.hostname === "") {
            return
        }
        var config = {
            headers: {'Access-Control-Allow-Origin': '*',
                'Accept': 'application/json'}
        };

        console.log("getData called: /api/data/" + this.state.hostname)
        Axios
            .get(`/api/data/`+this.state.hostname, config)
            .then(response => {
                /*
                console.log(response.data);
                console.log(response.status);
                console.log(response.statusText);
                console.log(response.headers);
                console.log(response.config);
                */

                //console.log(response.data)
                this.setState({ hostStats: response.data.slice(-100, -1) });
                let [tmpData, packetLoss] = dataTwist(this.state.hostStats);
                this.setState({ data: tmpData });
                this.setState({ packetLoss: packetLoss });
            })
            .catch(error => {
                if (error.response) {
                    // The request was made and the server responded with a status code
                    // that falls out of the range of 2xx
                    console.log(error.response.data);
                    console.log(error.response.status);
                    console.log(error.response.headers);
                } else if (error.request) {
                    // The request was made but no response was received
                    // `error.request` is an instance of XMLHttpRequest in the browser and an instance of
                    // http.ClientRequest in node.js
                    console.log(error.request);
                } else {
                    // Something happened in setting up the request that triggered an Error
                    console.log('Error', error.message);
                }
                console.log(error.config);
                //this.setState({ loading: false, error })
            })
    }

    render() {
        /*
                let data;
                let packetLoss;

                if (this.state.data) {
                    data = this.state.data;
                    packetLoss = this.state.packetLoss;
                } else {
                    data = [];
                    packetLoss = [];
                }

                if (this.state.hostStats) {
                    [data, packetLoss] = dataTwist(this.state.hostStats);
                } else {
                    data = [];
                    packetLoss = [];
                }
        */
        return (
            <div>
                <Container fluid>
                    <Notification>
                        Viewing: <strong>{this.state.hostname}</strong>
                    </Notification>
                </Container>

                <div>
                    <LineChart
                        data={this.state.data}
                        height="500px"
                        legend={true}
                        precision={4}
                        suffix="ms"
                        xtitle="Time"
                        ytitle="Latency (ms)"
                        dataset={{pointRadius: 1}}
                    />

                    <ColumnChart
                        data={this.state.packetLoss}
                        height="200px"
                        min={0} max={100}
                        legend={true}
                        precision={4}
                        suffix="%"
                        xtitle="Time"
                        ytitle="% Loss"
                    />

                    <Table fullwidth narrow>
                        <Table.Head>
                            <Table.Row>
                                <Table.Heading>Time</Table.Heading>
                                <Table.Heading>Min RTT</Table.Heading>
                                <Table.Heading>Avg RTT</Table.Heading>
                                <Table.Heading>Max RTT</Table.Heading>
                                <Table.Heading>Packet Loss %</Table.Heading>
                            </Table.Row>
                        </Table.Head>

                        <Table.Body>
                            <DataTable data={this.state.hostStats} host={this.state.hostname}/>
                        </Table.Body>
                    </Table>
                </div>
            </div>

        )
    }


}

export default DataView

