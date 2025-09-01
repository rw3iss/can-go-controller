import type { Component } from 'solid-js';

import logo from './logo.svg';
import styles from './App.module.scss';
import CANDashboard from './pages/CANDashboard/index';

const App: Component = () => {
    return (
        <div class={styles.App}>
            <CANDashboard />
        </div>
    );
};

export default App;
