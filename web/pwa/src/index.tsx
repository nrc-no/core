import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import {store} from './app/store';
import {Provider} from 'react-redux';
import * as serviceWorker from './serviceWorker';
import {AuthProvider} from 'oidc-react';
import * as log from "loglevel"
import {SQLContextProvider} from "./app/db";

log.setDefaultLevel(log.levels.TRACE)

const oidcConfig = {
    onSignIn: () => {
        // Redirect?
    },
    authority: 'https://dev-53701279.okta.com',
    clientId: '0oa2c6odfoi6s2JXH5d7',
    redirectUri: 'http://localhost:3001'

};

ReactDOM.render(
    <React.StrictMode>
        <AuthProvider
            scope={"openid profile email"}
            autoSignIn={false}
            {...oidcConfig} >
            <Provider store={store}>
                <SQLContextProvider>
                    <App/>
                </SQLContextProvider>
            </Provider>
        </AuthProvider>
    </React.StrictMode>,
    document.getElementById('root')
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
