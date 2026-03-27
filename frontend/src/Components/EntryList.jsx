import { useState, useEffect } from 'react';
import { fetchEntries } from '../api';


export default function EntryList({ username, date, editable=false }) {
    const [entries, setEntries] = useState([]);

    const meals = ["Breakfast", "Lunch", "Dinner", "Snack"];

    useEffect(() => {
        fetchEntries(username, date).then(setEntries);
    }, [username, date]);

    const flexStyle = {
        display: "flex",
        justifyContent: "space-between",
        alignItems: "center"
    };

    const greenButton = {
        '--pico-primary-background': '#43a047',
        '--pico-primary-border': '#43a047',
        '--pico-primary-hover-background': '#2e7d32',
        '--pico-primary-hover-border': '#2e7d32',
        '--pico-primary-focus': 'rgba(67, 160, 71, 0.25)',
        '--pico-primary-inverse': '#fff'
    };

    const redButton = {
        '--pico-primary-background': '#e53935',
        '--pico-primary-border': '#e53935',
        '--pico-primary-hover-background': '#c62828',
        '--pico-primary-hover-border': '#c62828',
        '--pico-primary-focus': 'rgba(229, 57, 53, 0.25)',
        '--pico-primary-inverse': '#fff'
    };

    return (
        <>
        { meals.map(meal => {
            return (
                <article>
                    <header>
                        <div style={ flexStyle }>
                            <h3 style={{ margin: 0 }}>{ meal }</h3>
                            <button style={ greenButton }>+ Entry</button>
                        </div>
                    </header>

                    <center role="group" style={{ margin: 0 }}>
                        <div style={{ color: "dodgerblue" }}>0.0p</div>
                        <div style={{ color: "green" }}>0.0c</div>
                        <div style={{ color: "orange" }}>0.0f</div>
                        <div style={{ color: "red" }}>0.0kcal</div>
                    </center>

                    <hr />

                    <div style={ flexStyle }>
                        <div>
                            Coffee
                            <small style={{ opacity: 0.5, display: "flex", gap: "1ch" }} >
                                <div>0.0p</div>
                                <div>0.0c</div>
                                <div>0.0f</div>
                                <div>0.0kcal</div>
                            </small>
                        </div>
                        <kbd>1 Serving</kbd>
                    </div>
                </article>
            );
        }) }
        </>
    );
}