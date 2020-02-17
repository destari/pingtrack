import React, { Component } from "react"
import "rbx/index.css"
import {Field, Control, Label, Input, Button, Icon, Column, Table } from "rbx"
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faUser } from '@fortawesome/free-solid-svg-icons'
import Layout from "../components/layout"
import Axios from "axios";

class Settings extends Component {
    constructor(props) {
        super(props);
        this.state = {
            hosts: [],
            hostname: "",
        }
        this.handleChange = this.handleChange.bind(this)
        this.addNew = this.addNew.bind(this)
        this.remHost = this.remHost.bind(this)
    }

    componentDidMount() {
        this.getHosts();
    }

    getHosts() {
        let config = {
            headers: {'Access-Control-Allow-Origin': '*',
                'Accept': 'application/json'}
        };
        Axios
            .get(`/api/hosts/`, config)
            .then(response => {
                this.setState({ hosts: response.data });
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
            })
    }

    deleteHost(host) {
        let config = {
            timeout: 2000,
            headers: {'Access-Control-Allow-Origin': '*',
                'Accept': 'application/json'},
                "Content-Type": "application/json"
        };
        Axios
            .delete(`/api/hosts/`+host, config)
            .then(response => {
                this.setState({ hosts: response.data });
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
            })
    }

    addHost(host) {
        let config = {
            timeout: 2000,
            headers: {'Access-Control-Allow-Origin': '*',
                'Accept': 'application/json'},
            "Content-Type": "application/json"
        };
        Axios
            .post(`/api/hosts/`, {hostname: host}, config)
            .then(response => {
                this.setState({ hosts: response.data });
                this.setState({hostname: ""})
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
            })
    }

    handleChange(evt) {
        const newhost = evt.target.value;
        this.setState({ hostname: newhost });
    }

     addNew(e) {
        e.preventDefault();
        this.addHost(this.state.hostname)
    }

    remHost(hostname, e) {
        console.log(hostname)
        //e.preventDefault();
        this.deleteHost(hostname)
    }

    render() {

        const removeHost = hostname => e => {
            console.log(e)
            console.log(hostname)
            e.preventDefault();
            this.deleteHost(hostname)
        }



        const HostsTable = (my) => {
            if (my.state.hosts) {
                const tableItems = my.state.hosts.map(function (hostname) {
                    return (
                        <Table.Row key={hostname}>
                            <Table.Cell>{hostname}</Table.Cell>
                            <Table.Cell><Button size="small" color="danger" key="{hostname}" onClick={() => {my.remHost(hostname)}}>REMOVE</Button></Table.Cell>
                        </Table.Row>
                    )
                });
                return tableItems;
            } else {
                return null;
            }
        };

        return (
            <Layout>
                <Column size="6">
                    <form>
                        <Field horizontal>
                            <Field.Label size="normal">
                                <Label>Ping Interval</Label>
                            </Field.Label>
                            <Field.Body>
                                <Field>
                                    <Control expanded iconLeft>
                                        <Input type="text" placeholder="Ping Interval in seconds" />
                                        <Icon size="small" align="left">
                                            <FontAwesomeIcon icon={faUser} />
                                        </Icon>
                                    </Control>
                                </Field>
                            </Field.Body>
                        </Field>

                        <Table fullwidth narrow>
                            <Table.Head>
                                <Table.Row>
                                    <Table.Heading>Hostname</Table.Heading>
                                    <Table.Heading>Action</Table.Heading>
                                </Table.Row>
                            </Table.Head>

                            <Table.Body>

                                {HostsTable(this)}

                                <Table.Row key="new">
                                    <Table.Cell><Input type="text" name="newhost" placeholder="Hostname or IP" size="small" value={this.state.hostname} onChange={this.handleChange}/></Table.Cell>
                                    <Table.Cell><Button size="small" color="success" onClick={this.addNew}>ADD</Button></Table.Cell>
                                </Table.Row>
                            </Table.Body>
                        </Table>

                    </form>
                </Column>
            </Layout>
        )
    }

}

export default Settings;