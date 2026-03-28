import { useState, useEffect } from 'react';
import { fetchEntries, fetchFoods } from '../api';
import { CreateEntryModal } from '../Components';

/**
 * Component to display a list of entries for a given user and date. Also displays macros for each 
 * meal and entry, and allows adding new entries.
 * @component
 * @param { Object } props - The component props
 * @param { string } props.username - the username to fetch entries for
 * @param { string } props.date - the date to fetch entries for, in YYYY-MM-DD format
 * @param { boolean } [props.editable=false] - whether the entries should be editable
 * @returns { JSX.Element } The rendered component
 */
export default function EntryList({ username, date, editable=false }) {
    const [entries, setEntries] = useState({
        breakfast: [],
        lunch: [],
        dinner: [],
        snack: []
    });

    const updateEntries = meal => newEntries => {
        setEntries(prev => ({ ...prev, [meal]: [ ...prev[meal], newEntries] }));
    };

    useEffect(() => {
        fetchEntries(username, date.value)
            .then(res => {
                if (res) {
                    const groupedEntries = res.reduce((acc, entry) => {
                        const meal = entry.meal_name.toLowerCase();
                        if (!acc[meal]) acc[meal] = [];
                        acc[meal].push(entry);
                        return acc;
                    }, { breakfast: [], lunch: [], dinner: [], snack: [] });

                    setEntries(groupedEntries);
                }
            });
    }, [username, date]);

    return (
        <>
            { Object.keys(entries).map(meal => (
                <EntryTable key={ meal } 
                            name={ meal.charAt(0).toUpperCase() + meal.slice(1) }
                            entries={ entries[meal] }
                            updateEntries={ updateEntries(meal) }
                            editable={ editable }
                            date={ date } />
            )) }
        </>
    );
}

function EntryTable({ name, entries, updateEntries, editable, date }) {
    const [entryModalOpen, setEntryModalOpen] = useState(false);

    const totals = entries.reduce((acc, entry) => {
        Object.keys(acc).forEach(key => {
            acc[key] += entry.food[key] ?? 0;
        });

        return acc;
    }, { protein: 0, carbs: 0, fat: 0, calories: 0 });

    const handleAddEntrySuccess = newEntry => {
        setEntryModalOpen(false);
        updateEntries(newEntry);
    }

    return (
        <article>
            <header style={{ display: "flex",
                             justifyContent: "space-between",
                             alignItems: "center" }}>
                <h3 style={{ margin: 0 }}>{ name }</h3>
                {editable && <button className="green" onClick={() => setEntryModalOpen(true)}>
                    + Entry
                </button>}
            </header>

            <center role="group" style={{ margin: 0 }}>
                <div style={{ color: "dodgerblue" }}>{ totals.protein }p</div>
                <div style={{ color: "green" }}>{ totals.carbs }c</div>
                <div style={{ color: "orange" }}>{ totals.fat }f</div>
                <div style={{ color: "red" }}>{ totals.calories }kcal</div>
            </center>

            <hr />

            <div style={{ display: "flex", flexDirection: "column", gap: "1rem", marginBottom: "1rem" }}>
                { entries.map(entry => <EntryRow key={ "ER-" + entry.id } entry={ entry } />) }
            </div>

            <CreateEntryModal key={ name }
                              isOpen={ entryModalOpen } 
                              onClose={ () => setEntryModalOpen(false) }
                              initialMeal={ name.toLowerCase() }
                              date={ date }
                              onSuccess={ handleAddEntrySuccess } />

        </article>
    );
}

/**
 * Component to display a single entry, with its food name, macros, and serving count. Used in 
 * EntryTable.
 * @param { Object } props - The component props
 * @param { Object } props.entry - The entry to display, with shape { food: { name, protein, carbs, 
 *                                 fat, calories }, food.serving_count }
 * @returns { JSX.Element } The rendered component
 */
function EntryRow({ entry }) {
    return (
        <div style={{ display: "flex",
             justifyContent: "space-between",
             alignItems: "center" }}>
            <div>
                { entry.food.name }
                <small style={{ opacity: 0.5, display: "flex", gap: "1ch" }} >
                    <div>{ entry.food.protein * entry.servings }p</div>
                    <div>{ entry.food.carbs * entry.servings }c</div>
                    <div>{ entry.food.fat * entry.servings }f</div>
                    <div>{ entry.food.calories * entry.servings }kcal</div>
                </small>
            </div>
            <div>
                <kbd>{ entry.food.serving_count * entry.servings } { entry.food.serving_size }</kbd>
            </div>
        </div>
    );
}