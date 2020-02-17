import React, { Component } from "react"
import "rbx/index.css"
import {Column, Menu} from "rbx";
import Layout from "../components/layout"
import DataView from "../components/dataview";
import Axios from "axios"

class IndexPage extends Component {

    constructor(props) {
        super(props);
        this.state = {
            hosts: [],
            selected: ""
        };
        this.dataViewElement = React.createRef()
    }


    /*
        declare a member variable to hold the interval ID
        that we can reference later.
    */
    intervalID;
    intervalIDHosts;

    componentDidMount() {
        /*
          need to make the initial call to getData() to populate
         data right away
        */
        this.getHosts();
        //this.getData();

        /*
          Now we need to make it run at a specified interval,
          bind the getData() call to `this`, and keep a reference
          to the interval so we can clear it later.
        */
        //this.intervalID = setInterval(this.getData.bind(this), 5000);
        this.intervalIDHosts = setInterval(this.getHosts.bind(this), 5000);
    }

    componentWillUnmount() {
        /*
          stop getData() from continuing to run even
          after unmounting this component
        */
        //clearInterval(this.intervalID);

        // Uncomment if we want to stop fetching hosts:
        clearInterval(this.intervalIDHosts);
    }

    getHosts = () => {
        let config = {
            headers: {'Access-Control-Allow-Origin': '*',
                'Accept': 'application/json'}
        };
        Axios
            .get(`/api/hosts/`, config)
            .then(response => {
                /*
                console.log(response.data);
                console.log(response.status);
                console.log(response.statusText);
                console.log(response.headers);
                console.log(response.config);
                */
                console.log(response.data);
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
                //this.setState({ loading: false, error })
            })
    }


    render() {
        const setSelected = (hostname) => {
            this.setState({ selected: hostname })
            if (this.dataViewElement.current) {
                this.dataViewElement.current.updateHostname(hostname)
            }

        }

        const selected = () => {
            return this.state.selected;
        }

        const makeMenu = () => {
            if (this.state.hosts) {
                const menuItems = this.state.hosts.map(function (hostname ) {
                    return (
                        <Menu.List.Item key={hostname} active={ hostname === selected()} onClick={() => {setSelected(hostname)}}>{hostname}</Menu.List.Item>
                    )
                });
                return menuItems;
            } else {
                return null;
            }
        };

        return (
            <Layout>
                <Column.Group>
                    <Column size="2">
                        <Menu>
                            <Menu.Label>Hosts</Menu.Label>
                            <Menu.List>
                                { makeMenu() }

                            </Menu.List>
                        </Menu>
                    </Column>
                    <Column>
                        {this.state.selected ? (
                            <DataView host={this.state.selected} ref={this.dataViewElement}> </DataView>
                        ) : (
                            <p>Select a host from the left menu to see details.</p>
                        )}
                    </Column>
                </Column.Group>

            </Layout>
        )
    }

}

export default IndexPage

