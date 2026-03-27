import { useState, useEffect } from 'react';
import { fetchEntries } from '../api';

/**
 * Component to display a list of entries for a given user and date. Also displays macros for each 
 * meal and entry, and allows adding new entries.
 * @param { Object } props - The component props
 * @param { string } props.username - the username to fetch entries for
 * @param { string } props.date - the date to fetch entries for, in YYYY-MM-DD format
 * @param { boolean } [props.editable=false] - whether the entries should be editable
 * @returns { JSX.Element } The rendered component
 */
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

    return (
        <>
        { meals.map(meal => (
            <article>
                <header>
                    <div style={ flexStyle }>
                        <h3 style={{ margin: 0 }}>{ meal }</h3>
                        {editable && <button className="green">+ Entry</button>}
                    </div>
                </header>

                {(() => {
                    const mealTotals = { protein: 0, carbs: 0, fat: 0, calories: 0 };

                    let list = entries.filter(e => e.meal === meal).map(entry => {
                        mealTotals.protein += entry.food.protein;
                        mealTotals.carbs += entry.food.carbs;
                        mealTotals.fat += entry.food.fat;
                        mealTotals.calories += entry.food.calories;

                        return (
                            <div style={ flexStyle }>
                                <div>
                                    { entry.food.name }
                                    <small style={{ opacity: 0.5, display: "flex", gap: "1ch" }} >
                                        <div>{ entry.food.protein }p</div>
                                        <div>{ entry.food.carbs }c</div>
                                        <div>{ entry.food.fat }f</div>
                                        <div>{ entry.food.calories }kcal</div>
                                    </small>
                                </div>
                                <kbd>{ entry.food.serving_count } Serving(s)</kbd>
                            </div>
                        );
                    });

                    return (
                        <>
                        <center role="group" style={{ margin: 0 }}>
                            <div style={{ color: "dodgerblue" }}>{ mealTotals.protein/100 }p</div>
                            <div style={{ color: "green" }}>{ mealTotals.carbs/100 }c</div>
                            <div style={{ color: "orange" }}>{ mealTotals.fat/100 }f</div>
                            <div style={{ color: "red" }}>{ mealTotals.calories/100 }kcal</div>
                        </center>
                        <hr />

                        { list }
                        </>
                    );
                })()}
            </article>
        )) }
        </>
    );
}