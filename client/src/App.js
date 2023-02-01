import React, { useState } from 'react';
import { Button } from 'semantic-ui-react';
import { BrowserRouter, Redirect, Route, Switch } from "react-router-dom";
import axios from "axios";
import Lobby from "./Lobby";
import './App.css';
import { Icon } from 'semantic-ui-react'
import { NavLink } from 'react-router-dom';

const API_URL = process.env.REACT_APP_API_URL ?? "http://localhost:80/v1"

export const App = () => {

  const [redirect, setRedirect] = useState(null)

  const createGame = () => {
    axios.post(API_URL + "/lobby",
      { "num_clients": 0 },
      {
        headers: { "Access-Control-Allow-Origin": "*" }
      }).then(res => {
        setRedirect("/lobby/" + res.data.slug)
      })
  }

  const joinGame = () => {
    let slug = prompt("Enter lobby code:")
    if (slug !== null) {
      setRedirect("/lobby/" + slug)
    } else {
      setRedirect("/")
    }
  }

  return (
    <BrowserRouter>
      {redirect && <Redirect to={redirect} />}
      <main className="content-container">
        <Switch>
          <Route path="/lobby/:id" component={Lobby} />
          <Route path="/fail" render={() => { }} />
          <Route path="/" render={() => {
            return (
              <div className="App">
                <header className="App-header">
                  <div className="menu" id="navbar">
                  <div className="nav-menu">
                        <div id="home-nav">
                            <NavLink to="/" style = {{color: "white"}}> <Icon name='home'/> </NavLink>
                        </div>

                    </div>
                  </div>
                  <div>
                    <h1>Mancala</h1>
                  </div>

                  <div className="directions">
                    The mancala games are a family of two-player turn-based strategy board games played with small stones, beans, or seeds and rows of holes or pits in the earth, a board or other playing surface. The objective is usually to capture all or some set of the opponent's pieces.

                    Versions of the game date back past the 3rd century and evidence suggests the game existed in Ancient Egypt. It is among the oldest known games to still be widely played today.
                  </div>
                  <div className="home-buttons">
                    <div className = "flex-button">
                      <Button size="large" onClick={createGame}>Create Game</Button>
                    </div>
                    <div className = "flex-button">
                      <Button size="large" onClick={joinGame}>Join Game</Button>
                    </div>
                  </div>
                </header>
              </div>
            );
          }} />
        </Switch>
      </main>
    </BrowserRouter>
  );
}

export default App;
