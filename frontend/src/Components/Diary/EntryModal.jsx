import { useState, useEffect, useContext, useRef } from 'react';
import { MdAdd, MdRefresh, MdStar } from "react-icons/md";

import { handleSearchFoods, handleCreateEntry } from "../../api";
import { SessionContext } from '../../Context';
import Modal from './Modal';

const MEAL_OPTIONS = [
    { value: "breakfast", label: "Breakfast" },
    { value: "lunch", label: "Lunch" },
    { value: "dinner", label: "Dinner" },
    { value: "snacks", label: "Snack" },
];

export function AddEntryModal({
    isOpen,
    onClose,
    onSubmit,
    initialMeal = "breakfast"
}) {
    const handleSubmit = ({ meal, food, servings }) => {
        if (!food || !meal || !servings) return new Error("Missing required fields.");

        // handleCreateEntry(food.id, meal, new Date(), servings)
    }

    return <EntryModal 
        isOpen={ isOpen } 
        onClose={ onClose } 
        onSubmit={ handleSubmit }
        initialMeal={ initialMeal } />;
}

export function EditEntryModal({
    isOpen,
    onClose,
    onSubmit,
    initialMeal,
    initialFood,
    initialServings
}) {
    return <EntryModal 
        isOpen={ isOpen } 
        onClose={ onClose }
        onSubmit={ onSubmit }
        initialMeal={ initialMeal }
        initialFood={ initialFood }
        initialServings={ initialServings }
        editMode={ true } />;
}

export default function EntryModal({ 
    isOpen, 
    onClose,
    onSubmit,
    initialMeal = "breakfast",
    initialFood = null,
    initialServings = 1,
    editMode = false
}) {
    const [meal, setMeal] = useState(initialMeal);
    const [food, setFood] = useState(initialFood);
    const [servings, setServings] = useState(initialServings);
    const [searchString, setSearchString] = useState(initialFood?.name ?? "");
    const [searchResults, setSearchResults] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    const session = useContext(SessionContext);

    // form reset
    useEffect(() => {
        if (!isOpen) return;

        setMeal(initialMeal ?? "breakfast");
        setFood(initialFood ?? null);
        setServings(initialServings ?? 1);
        setSearchString(initialFood?.name ?? "");
        setSearchResults(initialFood ? [initialFood] : []);
        setError("");
    }, [isOpen, initialMeal, initialFood, initialServings]);

    // search
    const version = useRef(0);
    const search = e => {
        e?.preventDefault();

        setFood(null);
        setServings(1);
        setLoading(true);
        setError("");

        const currentVersion = ++version.current;
        const t = setTimeout(() => {
            handleSearchFoods(searchString.trim())
                .then(data => {
                    if (currentVersion !== version.current) return;
                    setSearchResults(data ?? []);
                })
                .catch(err => {
                    if (currentVersion !== version.current) return;
                    setError(err.error ?? "Failed to fetch search results.");
                })
                .finally(() => currentVersion === version.current && setLoading(false));
        }, 350);

        return t;
    };

    // search effect
    useEffect(() => {
        if (!isOpen) return;

        const trimmed = searchString.trim();
        if (trimmed.length < 2) {
            setSearchResults(food ? [food] : []);
            setLoading(false);
            return;
        }

        const t = search();
        return () => clearTimeout(t);
    }, [searchString]);

    const handleSubmit = e => {
        e.preventDefault();

        // setLoading(true);
        onSubmit?.({
            meal, food, servings
        }).finally(() => setLoading(false));
    }

    return <Modal title={ editMode ? "Edit Entry" : "Add Entry" } isOpen={ isOpen } 
        onClose={ onClose }>
        <style>{`
            .entry-form {
                padding: var(--pico-block-spacing-vertical) var(--pico-block-spacing-horizontal);
                display: flex;
                flex-direction: column;
                max-height: calc(100vh - var(--pico-block-spacing-vertical) * 4 - 1.6rem);
                gap: var(--pico-spacing);

                &  *, & input { margin: 0; }

                & > div[role="search"] button {
                    line-height: 0;
                    padding: var(--pico-form-element-spacing-vertical) 
                        calc(var(--pico-form-element-spacing-horizontal) * 0.75);
                }

                & > input { margin-bottom: 0; }
                & label { width: 100%; }

                & .food-search-results {
                    flex: 1;
                    min-height: 4rem;
                    overflow: auto;
                }

                & .search-results {
                    background-color: hsl(from var(--pico-color) h s l / 2.5%);

                    & > div {
                        padding: calc(var(--pico-spacing) * 0.5);
                        opacity: 0.5;
                    }

                    & label {
                        display: flex;
                        gap: calc(var(--pico-spacing) * 0.5);
                        padding: calc(var(--pico-spacing) * 0.5);
                        margin: 0;
                    }

                    & label:hover {
                        cursor: pointer;
                        background-color: hsl(from var(--pico-color) h s l / 5%);
                    }

                    & label:active {
                        background-color: hsl(from var(--pico-color) h s l / 10%);
                    }

                    & input[type="radio"] { margin:0; }
                    & small { opacity: 0.5; }

                    & .star {
                        color: gold;

                        & > * { transform: translateY(-0.12rem); }
                    }


                }

                & .totals {
                    display: flex;
                    justify-content: space-between;

                    & > div { flex: 1; }
                    & .p { color: var(--protein-color); }
                    & .c { color: var(--carb-color); }
                    & .f { color: var(--fat-color); }
                    & .kcal { color: var(--kcal-color); }
                }

                & .submit-row {
                    display: flex;
                    gap: calc(var(--pico-spacing) * 0.5);
                }
            }
        `}</style>

        <form className="entry-form" onSubmit={ onSubmit }>
            <label>
                Meal
                <select value={ meal } onChange={ e => setMeal(e.target.value) } required>
                    { MEAL_OPTIONS.map(opt => 
                        <option key={ opt.value } value={ opt.value }>{ opt.label }</option>) }
                </select>
            </label>
            <div role="search">
                <button className="green" title="Add a new food item" disabled>
                    <MdAdd />
                </button>
                <input
                    type="search" 
                    value={ searchString }
                    onChange={ e => setSearchString(e.target.value) }
                    placeholder="Search for food..." 
                />
                <button title="Refresh search results" onClick={ search }>
                    <MdRefresh />
                </button>
            </div>
            <fieldset className="food-search-results">
                <legend>Food</legend>
                <div className="search-results">
                    { loading ? <div className="loading" aria-busy="true">Loading...</div> : (<>
                        { searchResults.length === 0 ? <div>No results found.</div> : 
                            searchResults.map(f => {
                                if (food && food.id !== f.id) return;
                                return (
                                    <label key={ f.id }>
                                        <div>
                                            <input type="radio" name="food_id" value={ f.id } 
                                                onClick={ () => setFood(f) } required />
                                        </div>
                                        <div>
                                            <div>
                                                { f.username === session?.username && 
                                                    <span className="star" 
                                                        title="You created this item">
                                                        <MdStar title="You created this item!" />
                                                    </span>
                                                }
                                                <span>{ f.name }</span>
                                            </div>
                                            {  f.brand && <><small>{ f.brand }</small><br/></> }
                                            <small>
                                                { f.protein }p
                                                , { f.carbs }c
                                                , { f.fat }f
                                                , { f.calories }kcal
                                                <strong>
                                                    &nbsp;in { f.serving_count } { f.serving_size }
                                                </strong>
                                            </small>
                                        </div>
                                    </label>
                                );
                            }) 
                        }
                    </>) }
                </div>
            </fieldset>
            <label>
                Servings { food && `( ${food?.serving_size} )` }
                <input type="number" min={ editMode ? 0 : 1 } step="0.01" value={ servings } 
                    onChange={ e => setServings(e.target.value) } required />
            </label>
            <center className="totals">
                <div className="p">
                    { food ? (food.protein / food.serving_count * servings).toFixed(2) : 0 }p</div>
                <div className="c">
                    { food ? (food.carbs / food.serving_count * servings).toFixed(2) : 0 }c</div>
                <div className="f">
                    { food ? (food.fat / food.serving_count * servings).toFixed(2) : 0 }f</div>
                <div className="kcal">
                    { food ? (food.calories / food.serving_count * servings).toFixed(2) : 0 }kcal
                </div>
            </center>
            <fieldset className="submit-row">
                { editMode && 
                    <input type="submit" className="red" value="Delete" onClick={ handleSubmit } />}
                    <input type="submit" value="Save" onClick={ handleSubmit } 
                        disabled={ !food || servings === "0" } />
            </fieldset>
        </form>
    </Modal>;
}