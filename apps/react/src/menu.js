import React, { Component } from 'react';
import { Link, Route, Switch } from 'react-router-dom'

// import Dashboard from './components/dashboard'
import Alerts from './components/alerts/alerts'
import News from './components/news/news'
import MySpace from './components/myspace/myspace'

import './main.css';

class Menu extends Component {
    render() {
        return (
            <div className="Menu">
            <ul>
            <li><Link to='/'>Alerts</Link></li>
            <li><Link to='/news'>News</Link></li>
            <li><Link to='/myspace'>My Space</Link></li>
            </ul>
            </div>
        );
    }
}

class Workspace extends Component {
    render() {
        return (
            <div className="Workspace">
            <Switch>
            <Route path='/news' component={News} />
            <Route path='/myspace' component={MySpace} />
            <Route path='/' component={Alerts} />
            </Switch>
            </div>
        );
    }
}


export {Menu, Workspace}


// <div className="rt-cnt">
//     </div>
