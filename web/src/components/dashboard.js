import React from "react"
import "rbx/index.css"
import { Message, Column, Section, Notification } from "rbx";

export default (props) => (
    <Section>
        <Notification color="link" size={6}>
            {' '}
            <p>Click a host in the menu to the left to see charts.</p>
        </Notification>
        <Column.Group gapSize={8} centered vcentered={true}>
            <Column size="3" >
                <Message color="dark">
                    <Message.Header>
                        <p>Total Pings</p>
                    </Message.Header>
                    <Message.Body textAlign="centered">
                        {' '}
                        <strong>{props.config.PingCount}</strong>
                    </Message.Body>
                </Message>
            </Column>

            <Column size="3">
                <Message color="dark">
                    <Message.Header>
                        <p># of Hosts</p>
                    </Message.Header>
                    <Message.Body textAlign="centered">
                        {' '}
                        <strong>{props.config.Hosts ? props.config.Hosts.length : 0}</strong>
                    </Message.Body>
                </Message>
            </Column>

            <Column size="3">
                <Message color="dark">
                    <Message.Header>
                        <p>Ping Interval</p>
                    </Message.Header>
                    <Message.Body textAlign="centered">
                        {' '}
                        <strong>{props.config.EchoTimes ? props.config.EchoTimes : 0} seconds</strong>
                    </Message.Body>
                </Message>
            </Column>
        </Column.Group>
    </Section>


)