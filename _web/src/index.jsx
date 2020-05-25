import React from 'react'
import ReactDOM from 'react-dom'
import { Provider } from 'react-redux'
import store from './store' // redux store setup
import AppContainer from './AppContainer.jsx' // app entry point

import configData from '../dist/config.json'

ReactDOM.render((
  <Provider store={store}>
    <AppContainer config={config()} />
  </Provider>
), document.getElementById('app'))

function config() {
  return { ...configData, apiHost: location.hostname }
}