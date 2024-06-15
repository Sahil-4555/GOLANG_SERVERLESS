import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import {thunk} from 'redux-thunk';
import reducers from './redux'
import './index.css';
import { legacy_createStore, applyMiddleware, compose} from 'redux'
import { Provider } from 'react-redux';

const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;

const store = legacy_createStore(
    reducers,
    composeEnhancers(applyMiddleware(thunk))
)

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    <Provider store={store}>
        <App />
    </Provider>
);
