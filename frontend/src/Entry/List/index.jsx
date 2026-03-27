import { useState, useEffect } from 'react';
import { fetchEntries } from '../../api';
import "./index.css";

export default function EntryList({ user_name, date }) {
    const [entries, setEntries] = useState([]);
    const meals = ["Breakfast", "Lunch", "Dinner", "Snack"];

    useEffect(() => {
        fetchEntries(user_name, date)
            .then(setEntries)
            .catch(console.error);
    }, [user_name, date]);

    return (
        <div className="entry-list">

            { meals.map(meal => (
                <div className="meal-section" key={meal}>
                    <h2>{ meal }</h2>
                    { entries.filter(e => e.meal === meal).map(e => (
                        <div className="entry">
                            <div>{ e.food.name }</div>
                            <div>{ e.food.calories } kcal</div>
                        </div>
                    )) }
                </div>
            )) }

        </div>
    )
}