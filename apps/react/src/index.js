import React from 'react'
import ReactDOM from 'react-dom'
import { BrowserRouter } from "react-router-dom"
import App from './app'

const config = require('./config');

console.log(config);

let aaa = <BrowserRouter><App /></BrowserRouter>;

ReactDOM.render(
  aaa,
  document.getElementById('root')
);
