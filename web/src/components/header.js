import React from "react"
import "rbx/index.css"
import { Navbar } from "rbx";

export default (props) => (
    <Navbar color="warning">
        <Navbar.Brand>
            <Navbar.Item href="/">
                <h1>{props.headerText}</h1>
            </Navbar.Item>
            <Navbar.Burger />
        </Navbar.Brand>
        <Navbar.Menu>
            <Navbar.Segment align="end">
                <Navbar.Item href="/" boxed="true">Dashboard</Navbar.Item>
                <Navbar.Item href="/about/">About</Navbar.Item>
                <Navbar.Item href="/settings/">Settings</Navbar.Item>
            </Navbar.Segment>
        </Navbar.Menu>
    </Navbar>
)