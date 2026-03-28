import { useState, useEffect } from 'react';
import { fetchEntries, fetchFoods } from '../api';
import { CreateEntryModal } from '../Components';

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

    useEffect(() => {
        fetchEntries(username, date.value).then(setEntries);
    }, [username, date]);

    const breakfast = entries.filter(e => e.meal === "breakfast");
    const lunch = entries.filter(e => e.meal === "lunch");
    const dinner = entries.filter(e => e.meal === "dinner");
    const snack = entries.filter(e => e.meal === "snack");

    return (
        <>
            <EntryTable name="Breakfast" entries={ breakfast } editable={ editable } />
            <EntryTable name="Lunch" entries={ lunch } editable={ editable } />
            <EntryTable name="Dinner" entries={ dinner } editable={ editable } />
            <EntryTable name="Snack" entries={ snack } editable={ editable } />
        </>
    );
}

function EntryTable({ name, entries, editable }) {
    const [entryModalOpen, setEntryModalOpen] = useState(false);

    const totals = entries.reduce((acc, entry) => {
        Object.keys(acc).forEach(key => {
            acc[key] += entry.food[key] ?? 0;
        });

        return acc;
    }, { protein: 0, carbs: 0, fat: 0, calories: 0 });

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

            { entries.map(entry => <EntryRow entry={ entry } />) }

            <center role="group" style={{ margin: 0 }}>
                <div style={{ color: "dodgerblue" }}>{ totals.protein/100 }p</div>
                <div style={{ color: "green" }}>{ totals.carbs/100 }c</div>
                <div style={{ color: "orange" }}>{ totals.fat/100 }f</div>
                <div style={{ color: "red" }}>{ totals.calories/100 }kcal</div>
            </center>

            <CreateEntryModal isOpen={ entryModalOpen } 
                            onClose={ () => setEntryModalOpen(false) }
                            initialMeal={ name.toLowerCase() } />

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
                    <div>{ entry.food.protein }p</div>
                    <div>{ entry.food.carbs }c</div>
                    <div>{ entry.food.fat }f</div>
                    <div>{ entry.food.calories }kcal</div>
                </small>
            </div>
            <kbd>{ entry.food.serving_count } Serving(s)</kbd>
        </div>
    );
}

function AddEntryForm({ isOpen, onSubmit, onCancel }) {
    const [foods, setFoods] = useState([]);
    const [formData, setFormData] = useState({
        food: null,
        servings: 1
    });

    useEffect(() => {
        fetchFoods().then(setFoods);
    }, []);

    const handleChange = e => {
        const { name, value } = e.target;
        setFormData(p => ({ ...p, [name]: value }));

        console.log(formData);
    }

    if (!isOpen) return null;

    return (
        <form onSubmit={ onSubmit }>
            <hr />
            <fieldset>
                <label>
                    Food
                    <select name="food.id"
                            onChange={ handleChange }>
                        <option selected disabled value="">
                            Pick a food you'd like to add...
                        </option>
                        { foods.map(food => (
                            <option value={ food.id }>{ food.name }</option>
                        )) }
                    </select>
                </label>

                <label>
                    Servings { formData.food?.serving_size }
                    <input name="servings" 
                           type="number" 
                           placeholder="Enter number of servings"
                           value={ formData.servings }
                           onChange={ handleChange }
                           min="0"
                           step="0.5"
                           required />
                </label>
            </fieldset>
        </form>
    );
}