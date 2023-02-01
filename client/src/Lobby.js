import React, {Component} from 'react';
import {Button, Header, Icon, Modal} from 'semantic-ui-react';
import axios from "axios";
import {CopyToClipboard} from 'react-copy-to-clipboard';
import './App.css';
import {NavLink} from 'react-router-dom';

const API_URL = process.env.REACT_APP_API_URL ?? "http://localhost:80/v1/"
const WEBSOCKET = process.env.REACT_APP_WEBSOCKET ?? "ws://localhost:80/ws"

class Lobby extends Component {
    constructor(props) {
        super(props);
        const slug = this.props.match.params.id;
        this.state = {
            slug: slug,
            // gameState information
            // 0 means waiting for player to join
            // 1 means player 1's turn
            // 2 means player 2's turn
            // 3 means final result displayed + update score + next round button appears
            // 4 is waiting page for both players to hit next round
            gameState: 0,
            playerNumber: 0,
            myScore: 0,
            opponentScore: 0,
            modal: true,
            myAction: "",
            opponentAction: "",
            winner: 0,
            rules: false,
            finish: 0,
            waiting: false,
        }
        this.buttonPanel = this.buttonPanel.bind(this)
        this.score = this.score.bind(this)
        this.join = this.join.bind(this)
        this.handleMessage = this.handleMessage.bind(this)
        this.handleJoin = this.handleJoin.bind(this)
        this.move = this.move.bind(this)
        this.handleState = this.handleState.bind(this)
    }

    componentDidMount() {
        this.conn = new WebSocket(WEBSOCKET);
        this.conn.onmessage = (message) => {
            const messageArray = message.data.split("\n");
            for (const messageData of messageArray) {
                console.log(messageData)
                this.handleMessage(JSON.parse(messageData));
            }
        };
    }

    handleMessage(message) {
        // filter messages only relevant to current room
        if (message.slug === this.state.slug) {
            if (message.type === "join") this.handleJoin(message.player_number);
            else if (message.type === "state") this.handleState(message.player_number, message.action);
            else if (message.type === "round" && message.valid) this.handleCards(message);
            else if (message.type === "finish") this.handleFinish(message);
        }
    }

    handleJoin(playerNumber) {
        console.log("player number: " + playerNumber)
    }


    join() {
        axios.put(API_URL + "lobby/" + this.state.slug + "/join",
            {},
            {
                headers: {"Access-Control-Allow-Origin": "*"}
            }).then((res) => {
            let joinPacket = {
                type: 'join',
                slug: this.state.slug
            };

            // Connected to remote
            this.conn.send(JSON.stringify(joinPacket));

            this.setState({
                modal: false
            })
        })
    }

    move(pit_index) {
        let turnPacket = {
            type: 'turn',
            player_number: this.state.playerNumber,
            slug: this.state.slug,
            pit_index: pit_index
        };

        // Connected to remote
        this.conn.send(JSON.stringify(turnPacket));
    }

    buttonPanel() {
        if (this.state.waiting) {
            return (
                <div className="message">
                    Waiting for opponent...
                </div>
            )
        }

        // render table with 2 rows and 6 buttons each
        return (
            <div>
                <div className="message"> It's your turn!</div>

                <div className="">Opponent side</div>
                <div className="button-panel">
                    <div className="column">
                        <Button size="large" color="yellow" onClick={() => {
                        }}>Opponent's pit</Button></div>
                    <div className="column">
                        <Button size="large" active={false}>Pit #1</Button></div>

                    <div className="column">
                        <Button size="large" active={false}>Pit #2</Button></div>
                    <div className="column">
                        <Button size="large" active={false}>Pit #3</Button></div>
                    <div className="column">
                        <Button size="large" active={false}>Pit #4</Button></div>
                    <div className="column">
                        <Button size="large" active={false}>Pit #5</Button></div>
                    <div className="column">
                        <Button size="large" active={false}>Pit #6</Button></div>
                </div>
                <div className="">Your side</div>
                <div className="button-panel">
                    <div className="column">
                        <Button size="large" onClick={() => {
                            this.move(0)
                        }}>Pit #1</Button></div>

                    <div className="column">
                        <Button size="large" onClick={() => {
                            this.move(1)
                        }}>Pit #2</Button></div>
                    <div className="column">
                        <Button size="large" onClick={() => {
                            this.move(2)
                        }}>Pit #3</Button></div>
                    <div className="column">
                        <Button size="large" onClick={() => {
                            this.move(3)
                        }}>Pit #4</Button></div>
                    <div className="column">
                        <Button size="large" onClick={() => {
                            this.move(4)
                        }}>Pit #5</Button></div>
                    <div className="column">
                        <Button size="large" onClick={() => {
                            this.move(5)
                        }}>Pit #6</Button></div>

                    <div className="column">
                        <Button size="large" color="yellow" onClick={() => {
                        }}>Your pit</Button></div>
                </div>
            </div>
        )
    }

    score() {
        return (
            <div className="scoreboard">
                <div className="score">
                    Your score: {this.state.myScore}
                </div>
                <div className="score">
                    Opponent's score: {this.state.opponentScore}
                </div>
            </div>
        )

    }

    render() {
        return (
            <div className="App">
                <header className="App-header">
                    <div className="menu" id="navbar">
                        <div className="nav-menu">
                            <div id="home-nav">
                                <NavLink to="/" style={{color: "white"}}> <Icon name='home'/> </NavLink>
                            </div>

                        </div>
                    </div>
                    <div className="game-code-panel">
                        <div className="game-code">
                            Game code is: {this.state.slug}
                            <CopyToClipboard text={this.state.slug}
                            >
                                <span id="clipboard-icon"><Icon name='clipboard'/></span>
                            </CopyToClipboard>
                        </div>
                    </div>

                    {this.buttonPanel()}
                    {this.score()}

                </header>
                <Modal
                    size='small'
                    closeOnDimmerClick={false}
                    open={this.state.modal}
                >
                    <Header icon='gamepad' content='Join Game'/>
                    <Modal.Content>
                        <div className="modal-content">
                            <p>
                                Welcome to the Mancala game!
                            </p>

                            <ul>
                                <li> Your game code is {this.state.slug} </li>
                                <li> Only share the game code with one other person.</li>
                                <li>Never refresh or press back; all progress will be lost.</li>
                                <li> There may be lag at times so please be patient.</li>
                                <li> If the site breaks down or is taking too long, both players should exit and create
                                    a new game.
                                </li>
                            </ul>
                        </div>

                    </Modal.Content>
                    <Modal.Actions>

                        <Button color='purple' onClick={this.join}>
                            Join
                        </Button>
                    </Modal.Actions>
                </Modal>
            </div>
        )
    }

    handleState(player_number, state) {
        console.log("handleState", player_number, state);
    }
}

export default Lobby;