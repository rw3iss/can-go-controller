import type { Component, Getter, Setter } from 'solid-js';
import { createSignal } from 'solid-js';
import appStyles from '../../App.module.scss';
import styles from './CANDashboard.module.scss';

const test_global = () => console.log(`GLOBAL`);

const subscribe = async (events: Getter<any>, setEvents: Setter<any>) => {
	console.log(`Subscribing...`)

	const eventSource = new EventSource('http://localhost:8080/events', {
		withCredentials: false,
	});

	eventSource.onmessage = function(event) {
		console.log(`msg:`, event.data);
		setEvents([event.data, ...events()])
	};

	eventSource.addEventListener('can-msg', (event) => {
		console.log(`can-msg:`, event);
		setEvents([event.data, ...events()])
	});

	eventSource.onerror = function(err) {
		console.error('EventSource failed:', err);
		//alert("Could not connect to server.");
		eventSource.close(); // Close connection on error
	};
}

const CANDashboard: Component = () => {
	const [events, setEvents] = createSignal([]);

	subscribe(events, setEvents);

	return (
		<div class={`${styles.dashboard} ${appStyles.page}`} >
			<header class={appStyles.header}>
				CAN Dashboard
			</header>

			Events:
			{events().map(e => (
				<div class="event">{e}</div>
			))}
		</div>
	);
};

export default CANDashboard;
