import { useState, useEffect, useRef } from 'react';
import { MdAdd, MdRefresh, MdStar } from "react-icons/md";

import {
    handleSearchFoods, 
    handleCreateEntry, 
    handleEditEntry, 
    handleDeleteEntry } from "../../api";
import { useSession } from '../../Context';
import Modal from './Modal';

const MEAL_OPTIONS = [
    { value: "breakfast", label: "Breakfast" },
    { value: "lunch",     label: "Lunch" },
    { value: "dinner",    label: "Dinner" },
    { value: "snacks",    label: "Snack" },
];

export default function EntryModal({
    date,
    isOpen, 
    onClose,
    onSuccess,
    onNewFood,
    initialMeal = "breakfast",
    initialFood = null,
    initialServings = 1,
    entryID = null // used to identify entry for editing, null for creating new entry
}) {
    const [meal, setMeal] = useState(initialMeal);
    const [food, setFood] = useState(initialFood);
    const [servings, setServings] = useState(initialServings);
    const [searchString, setSearchString] = useState(initialFood?.name ?? "");
    const [searchResults, setSearchResults] = useState([]);
    const [loadingFood, setLoadingFood] = useState(false);
    const [loading, setLoading] = useState(false);

    const { username, notifications } = useSession();

    // form reset
    useEffect(() => {
        if (!isOpen) return;

        setMeal(initialMeal ?? "breakfast");
        setFood(initialFood ?? null);
        setServings(initialServings ?? 1);
        setSearchString(initialFood?.name ?? "");
        setSearchResults(initialFood ? [initialFood] : []);
    }, [date, isOpen, initialMeal, initialFood, initialServings]);

    const clear = e => {
        e?.preventDefault(); // eww, fix this

        setFood(null);
        setServings(1);
    };

    // search
    const version = useRef(0);
    const search = e => {
        e?.preventDefault();

        setLoadingFood(true);

        const currentVersion = ++version.current;
        const t = setTimeout(() => {
            handleSearchFoods(searchString.trim())
                .then(data => {
                    if (currentVersion !== version.current) return;
                    setSearchResults(data ?? []);
                })
                .catch(err => {
                    if (currentVersion !== version.current) return;

                    console.error(err);
                    notifications.add({
                        heading: "Error searching foods",
                        details: err.message || "Please try again.",
                        type: "error",
                        ttl: 5000
                    });
                })
                .finally(() => currentVersion === version.current && setLoadingFood(false));
        }, 350);

        return t;
    };

    // search effect
    useEffect(() => {
        if (!isOpen) return;

        const trimmed = searchString.trim();
        if (trimmed.length < 2) {
            setSearchResults(food ? [food] : []);
            setLoadingFood(false);
            return;
        }

        const t = search();
        return () => clearTimeout(t);
    }, [searchString]);

    const handleSubmit = e => {
        e.preventDefault();

        if (!food || !meal || servings <= 0) return;
        // TODO: inform user of missing fields

        setLoading(true);
        const action = entryID
            ? handleEditEntry(entryID, food.id, meal, date, servings)
            : handleCreateEntry(food.id, meal, date, servings);

        action
            .then(onSuccess)
            .catch(err => {
                console.error(err);
                notifications.add({
                    heading: `Failed to ${entryID ? "edit" : "create"} entry`,
                    details: err.message || "Please try again.",
                    type: "error",
                    ttl: 5000
                });
            })
            .finally(() => setLoading(false));
    }

    const handleDelete = e => {
        e.preventDefault();

        if (!entryID) return;

        setLoading(true);
        handleDeleteEntry(entryID)
            .then(onSuccess)
            .catch(err => {
                console.error(err);
                notifications.add({
                    heading: "Failed to delete entry",
                    details: err.message || "Please try again.",
                    type: "error",
                    ttl: 5000
                });
            })
            .finally(() => setLoading(false));
    }

    return <Modal title={ entryID ? `Edit Entry (${ entryID })` : "Add Entry" } isOpen={ isOpen } 
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

        <form className="entry-form" onSubmit={ handleSubmit }>
            <label>
                Meal
                <select value={ meal } onChange={ e => setMeal(e.target.value) } required>
                    { MEAL_OPTIONS.map(opt => 
                        <option key={ opt.value } value={ opt.value }>{ opt.label }</option>) }
                </select>
            </label>
            <div role="search">
                <button className="green" title="Add a new food item" onClick={ onNewFood }>
                    <MdAdd />
                </button>
                <input
                    type="search" 
                    value={ searchString }
                    onChange={ e => setSearchString(e.target.value) }
                    placeholder="Search for food..." 
                />
                <button title="Refresh search results" onClick={ e => clear(e) || search() }>
                    <MdRefresh />
                </button>
            </div>
            <fieldset className="food-search-results">
                <legend>Food</legend>
                <div className="search-results">
                    { loadingFood ? <div className="loading" aria-busy="true">Loading...</div> : (<>
                        { searchResults.length === 0 ? <div>No results found.</div> : 
                            searchResults.map(f => {
                                if (food && food.id !== f.id) return;
                                return (
                                    <label key={ f.id }>
                                        <div>
                                            <input 
                                                type="radio" 
                                                name="food_id" 
                                                value={ f.id } 
                                                onClick={ () => setFood(f) } 
                                                defaultChecked={ food?.id === f.id }
                                                required />
                                        </div>
                                        <div>
                                            <div>
                                                { f.username === username && 
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
                <input type="number" min="0.01" step="0.01" value={ servings } 
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
                { entryID && 
                    <input type="submit" className="red" value="Delete" onClick={ handleDelete } />}
                    <input type="submit" value="Save"
                        disabled={ !food || servings === "0" || loading } />
            </fieldset>
        </form>
    </Modal>;
}