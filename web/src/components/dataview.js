import React from "react"
import "rbx/index.css"
import {Table, Container, Notification} from "rbx";
import {ColumnChart, LineChart} from 'react-chartkick'
import 'chart.js'

function DataGraph(props) {

    if (props.data.Results) {
        let data = props.data.Results;
        let hostname = props.hostname;

        const dataBlob = data[hostname].map(function (d) {
            return (
                <Table.Row key={d.Time}>
                    <Table.Cell>{d.Time}</Table.Cell>
                    <Table.Cell>{d.MinRtt/1000000} ms</Table.Cell>
                    <Table.Cell>{d.AvgRtt/1000000} ms</Table.Cell>
                    <Table.Cell>{d.MaxRtt/1000000} ms</Table.Cell>
                    <Table.Cell>{d.PacketLoss} %</Table.Cell>
                </Table.Row>
            )
        });

        return dataBlob;
    } else {
        return (
            <tr><td>Loading...</td><td>Loading...</td></tr>
        );
    }
}

function dataTwist(inputData) {
    let outputData = []
    let packetLoss = []

    if (inputData) {

        let series = {}
        let seriesMax = {}
        let seriesMin = {}
        let seriesLoss = {}
        inputData.forEach( function (datum) {
            series[datum.Time] = datum.AvgRtt/1000000
            seriesMax[datum.Time] = datum.MaxRtt/1000000
            seriesMin[datum.Time] = datum.MinRtt/1000000
            seriesLoss[datum.Time] = datum.PacketLoss
            }
        )

        outputData.push({"name": "Avg RTT", "data": series})
        outputData.push({"name": "Min RTT", "data": seriesMin})
        outputData.push({"name": "Max RTT", "data": seriesMax})

        packetLoss.push({"name": "Packet Loss", "data": seriesLoss})

    }

    return [outputData, packetLoss]
}

class DataView extends React.Component {


    render() {

        let data;
        let packetLoss;

        if (this.props.data && this.props.data.Results) {
            [data, packetLoss] = dataTwist(this.props.data.Results[this.props.host]);
        } else {
            data = [];
            packetLoss = [];
        }

        return (
            <div>
                <Container fluid>
                    <Notification>
                        Viewing: <strong>{this.props.host}</strong>
                    </Notification>
                </Container>

                <div>
                    <LineChart
                        data={data}
                        height="500px"
                        legend={true}
                        precision={4}
                        suffix="ms"
                        xtitle="Time"
                        ytitle="Latency (ms)"
                        dataset={{pointRadius: 1}}
                    />

                    <ColumnChart
                        data={packetLoss}
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
                            <DataGraph data={this.props.data} hostname={this.props.host}/>
                        </Table.Body>
                    </Table>
                </div>
            </div>

        )
    }


}

export default DataView

