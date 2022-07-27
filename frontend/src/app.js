import React from 'react';
import { Route, Switch } from 'react-router-dom';

import { NavBar, Footer, Loading } from "./components";
import { Home, Profile } from './views';
import ProtectedRoute from './authenitcation/protected-route'

import "./app.css";
import { useAuth0 } from "@auth0/auth0-react";

const App = () => {
  const { isLoading } = useAuth0();

  if (isLoading) {
    return <Loading />
  }


  return (
    <div id="app" className="d-flex flex-column h-100">
      <NavBar />
      <div className="container flex-grow-1">
        <div className="mt-5">
          <Switch>
            <Route path="/" exact component={Home} />
            <ProtectedRoute path="/profile" component={Profile} />
          </Switch>
        </div>
      </div>
      <Footer />
    </div>
  );
};

export default App;
