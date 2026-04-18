import { useState, useEffect } from 'react';
import { createPortal } from 'react-dom';
import { MdAdd, MdArrowBackIos, MdArrowForwardIos } from "react-icons/md";

import { handleGetDiary } from "../../api";
import EntryModal from './EntryModal';
import FoodModal from './FoodModal';

const valid = d => d instanceof Date && !isNaN(d.getTime());
const today = (new Date()).toLocaleString("sv-SE").slice(0, 10);
const dtos = d => valid(d) ? d.toISOString().slice(0, 10) : today;
const stod = s => new Date(s);

// Component to display a diary for a given date, with sections for each meal and daily totals.
// Props:
// - username: the username of the diary to display
// - defaultDate: the initial date to display, in YYYY-MM-DD format (defaults to today)
export default function Diary({ username, defaultDate }) {
    const [refresh, setRefresh] = useState(0);
    const [data, setData] = useState(null);
    const [date, setDate] = useState(defaultDate ?? today);
    const [loading, setLoading] = useState(true);

    const [AddEntryModalOpen, setAddEntryModalOpen] = useState(false);
    const [AddEntryInitialMeal, setAddEntryInitialMeal] = useState("breakfast");
    const [AddEntryInitialFood, setAddEntryInitialFood] = useState(null);
    const [AddEntryInitialServings, setAddEntryInitialServings] = useState(1);
    const [AddEntryEntryID, setAddEntryEntryID] = useState(null);

    const [AddFoodModalOpen, setAddFoodModalOpen] = useState(false);

    const handleAddEntrySuccess = () => {
        setAddEntryInitialFood(null);
        setAddEntryInitialServings(1);
        setAddEntryEntryID(null);
        setAddEntryModalOpen(false);
        setRefresh(prev => prev + 1);
    };

    const handleAddEntryOnClose = () => {
        setAddEntryInitialFood(null);
        setAddEntryInitialServings(1);
        setAddEntryEntryID(null);
        setAddEntryModalOpen(false);
    }

    const handleOpenAddFoodModal = e => {
        e.preventDefault();
        setAddFoodModalOpen(true);
    };

    const handleAddFoodModalSuccess = food => {
        setAddEntryInitialFood(food);
        setAddFoodModalOpen(false);
    };

    useEffect(() => {
        setLoading(true);

        handleGetDiary(username, date)
            .then(setData)
            .finally(() => setLoading(false));
            // TODO: handle errors
    }, [username, date, refresh]);

    const mealSectionProps = meal => {
        const key = meal.toLowerCase();
        return { 
            meal, diary: data?.[key], 
            addEntry: m => {
                setAddEntryInitialMeal(m);
                setAddEntryModalOpen(true);
            }, 
            editEntry: ({ id, food, servings }) => {
                setAddEntryInitialMeal(meal);
                setAddEntryInitialFood(food);
                setAddEntryInitialServings(servings);
                setAddEntryEntryID(id);
                setAddEntryModalOpen(true);
            } };
    }

    return (
        <>
            <div style={{ display: "flex", justifyContent: "center" }}>
                <DateSelector defaultDate={ date } onChange={ setDate } />
            </div>
            { loading ? <div aria-busy="true" style={{ marginTop: "3rem" }} /> : <>
            <DailyTotals totals={ data?.totals } />
            <hr />
            { !data?.breakfast ? null : <MealSection {...mealSectionProps("Breakfast")} /> }
            { !data?.lunch     ? null : <MealSection {...mealSectionProps("Lunch")} /> }
            { !data?.dinner    ? null : <MealSection {...mealSectionProps("Dinner")} /> }
            { !data?.snacks    ? null : <MealSection {...mealSectionProps("Snacks")} /> }
            </> }
            { createPortal(
                <EntryModal 
                    date={ date }
                    isOpen={ AddEntryModalOpen } 
                    onClose={ handleAddEntryOnClose } 
                    onNewFood={ handleOpenAddFoodModal }
                    initialMeal={ AddEntryInitialMeal } 
                    initialFood={ AddEntryInitialFood }
                    initialServings={ AddEntryInitialServings }
                    entryID={ AddEntryEntryID }
                    onSuccess={ handleAddEntrySuccess } />
                , document.body) }

            { createPortal(
                <FoodModal 
                    isOpen={ AddFoodModalOpen } 
                    onClose={ () => setAddFoodModalOpen(false) }
                    onSuccess={ handleAddFoodModalSuccess } />
                , document.body) }
        </>
    );
}

// Component to select a date, with buttons to increment/decrement the date by one day.
// Props:
// - onChange(d): callback function for when changes occur, d is the new date in YYYY-MM-DD format
// - defaultDate: the initial date to display, in YYYY-MM-DD format (defaults to today)
function DateSelector({ onChange, defaultDate }) {
    const [date, setDate] = useState(stod(defaultDate ?? today));

    const update = d => {
        setDate(valid(d) ? d : stod(today));
        onChange?.(dtos(d));
    };

    const add = n => update(new Date(date.getTime() + n * 24 * 60 * 60 * 1000));
    const handleChange = e => update(stod(e.target.value));

    return (
        <fieldset role="group" style={{ width: "fit-content" }}>
            <button onClick={ () => add(-1) } 
                    style={{ lineHeight: 0 }}
                    title="Previous day">
                <MdArrowBackIos />
            </button>
            <input type="date" value={ dtos(date) } onChange={ handleChange } />
            <button onClick={ () => add(1) } 
                    style={{ lineHeight: 0 }}
                    title="Next day">
                <MdArrowForwardIos />
            </button>
        </fieldset>
    );
}

// Component to display daily totals each macro, along with a percent-based bar visualization.
// Props:
// - totals: an object with the following properties:
//   - protein: total protein in grams
//   - carbs: total carbs in grams
//   - fat: total fat in grams
//   - calories: total calories
function DailyTotals({ totals }) {
    const t = totals ?? { protein: 0, carbs: 0, fat: 0, calories: 0 };
    const pct = key => {
        const n = (t[key] / (t.protein + t.carbs + t.fat) * 100).toFixed(0);
        return n > 0 ? n + "%" : "";
    };

    return <>
        <article className="daily-totals">
            <center className="totals color-children">
                <div>
                    <strong>{ t.protein.toFixed(1) } g</strong>
                    <div>Protein</div>
                </div>
                <div>
                    <strong>{ t.carbs.toFixed(1) } g</strong>
                    <div>Carb</div>
                </div>
                <div>
                    <strong>{ t.fat.toFixed(1) } g</strong>
                    <div>Fat</div>
                </div>
                <div>
                    <strong>{ t.calories.toFixed(1) } g</strong>
                    <div>Kcal</div>
                </div>
            </center>

            { t.protein + t.carbs + t.fat <= 0 ? null : <>
                <center className="bar bg-children size-children">
                    <div title={ `${ pct("protein") } Protein` }></div>
                    <div title={ `${ pct("carbs") } Carbs` }></div>
                    <div title={ `${ pct("fat") } Fat` }></div>
                </center>

                <center className="label color-children size-children">
                    <div title={ `${ pct("protein") } Protein` }>{ pct("protein") }</div>
                    <div title={ `${ pct("carbs") } Carbs` }>{ pct("carbs") }</div>
                    <div title={ `${ pct("fat") } Fat` }>{ pct("fat") }</div>
                </center>
            </>}
        </article>

<style>{`
    .daily-totals {
        & > .totals {
            display: flex;
            margin-bottom: 1rem;

            & > div { flex: 1; }
            & > div > div { line-height: 0.5rem; }
        }

        & > .bar {
            height: 1rem;
            display:flex;
            margin-top: 0.5rem;
            border-radius: 0.5rem;
            overflow: hidden;

            & > div { height: 100%; }
        }

        & > .label {
            height: 1rem;
            display: flex;
            gap: 0.5rem;
            margin-top: 0.25rem;
            line-height: 1rem;

            & > div { height: 100%; }
        }

        .color-children > *:nth-child(1) { color: var(--protein-color); }
        .color-children > *:nth-child(2) { color: var(--carb-color); }
        .color-children > *:nth-child(3) { color: var(--fat-color); }
        .color-children > *:nth-child(4) { color: var(--kcal-color); }

        .bg-children > *:nth-child(1) { background: var(--protein-color); }
        .bg-children > *:nth-child(2) { background: var(--carb-color); }
        .bg-children > *:nth-child(3) { background: var(--fat-color); }
        .bg-children > *:nth-child(4) { background: var(--kcal-color); }

        .size-children > *:nth-child(1) { flex: ${ t.protein * 100 }; }
        .size-children > *:nth-child(2) { flex: ${ t.carbs * 100 }; }
        .size-children > *:nth-child(3) { flex: ${ t.fat * 100 }; }
    }
`}</style>
    </>;
}

// Component to display a meal section with its entries and totals, with a button to add a new entry
// Props:
// - meal: the name of the meal (e.g. "Breakfast")
// - diary: an object with the following properties:
// TODO: create diary object structure
function MealSection({ meal, diary, addEntry, editEntry }) {
    const entry = e => (
        <div key={`${ e.id }@${ e.created }`} 
             style={{ display: "flex", padding: "0.5rem" }}
             className="entry"
             title="Edit entry"
             onClick={ () => editEntry(e) }>
            <div style={{ flex: 1 }}>
                { e.food.name }
                <div style={{ display: "flex", gap: "0.5rem", opacity: 0.5 }}>
                    <small>{ (e.servings * e.food.protein / e.food.serving_count).toFixed(2) }p</small>
                    <small>{ (e.servings * e.food.carbs / e.food.serving_count).toFixed(2) }c</small>
                    <small>{ (e.servings * e.food.fat / e.food.serving_count).toFixed(2) }f</small>
                    <small>{ (e.servings * e.food.calories / e.food.serving_count).toFixed(2) }kcal</small>
                </div>
            </div>
            <div style={{ display: "flex", flexDirection: "column", justifyContent: "center" }}>
                <kbd>{ e.servings } { e.food.serving_size }</kbd>
            </div>
        </div>
    );

    return (
        <article>
            <header style={{ display: "flex", gap: "0.75rem" }}>
                <button className="green" 
                        style={{ padding: "0.25rem", lineHeight: "0" }}
                        title="Add entry"
                        onClick={ () => addEntry(meal.toLowerCase()) }>
                    <MdAdd />
                </button>
                <div>{ meal }</div>
            </header>
            { diary.entries.length <= 0 ? <i>No entries</i> : <>
                <center role="group" style={{ margin: 0 }}>
                    <div style={{ color: "var(--protein-color)" }}>
                        { diary.totals.protein.toFixed(1) } p
                    </div>
                    <div style={{ color: "var(--carb-color)" }}>
                        { diary.totals.carbs.toFixed(1) } c
                    </div>
                    <div style={{ color: "var(--fat-color)" }}>
                        { diary.totals.fat.toFixed(1) } f
                    </div>
                    <div style={{ color: "var(--kcal-color)" }}>
                        { diary.totals.calories.toFixed(1) } kcal
                    </div>
                </center>
                <hr />
            </>}
            <style>
            {`
                .meal-entries .entry:hover {
                    cursor: pointer;
                    background-color: hsl(from var(--pico-color) h s l / 5%);
                    border-radius: 0.25rem;
                }

                .meal-entries .entry:active {
                    background-color: hsl(from var(--pico-color) h s l / 10%);
                }
            `}
            </style>
            <div className="meal-entries" 
                 style={{ display: "flex", flexDirection: "column", gap: "0.25rem" }}>
            {  diary.entries.map(entry) }
            </div>
        </article>
    );
}