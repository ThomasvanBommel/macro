import { useState, useEffect, useRef } from 'react';
import { MdAdd, MdRefresh } from "react-icons/md";

import { handleSearchFoods } from "../../api";
import Modal from './Modal';

const MEAL_OPTIONS = [
    { value: "breakfast", label: "Breakfast" },
    { value: "lunch", label: "Lunch" },
    { value: "dinner", label: "Dinner" },
    { value: "snacks", label: "Snack" },
];

export default function EntryModal({ 
    isOpen, 
    onClose,
    onSubmit,
    initialMeal = "breakfast",
    initialFood = null,
    initialServings = 1,
    title = "Add Entry",
    submitLabel = "Save"
}) {
    const [meal, setMeal] = useState(initialMeal);
    const [food, setFood] = useState(initialFood);
    const [servings, setServings] = useState(initialServings); // string?
    const [searchString, setSearchString] = useState(initialFood?.name ?? "");
    const [searchResults, setSearchResults] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

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
    useEffect(() => {
        if (!isOpen) return;

        const trimmed = searchString.trim();
        if (trimmed.length < 2) {
            setSearchResults(food ? [food] : []);
            setLoading(false);
            return;
        }

        setLoading(true);
        setError("");

        const currentVersion = ++version.current;
        const t = setTimeout(() => {
            handleSearchFoods(trimmed)
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

        return () => clearTimeout(t);
    }, [searchString]);
    


    // const refreshSearchResults = () => {
    //     setSelectedFood(null);

    //     handleSearchFoods(query)
    //         .then(setFood);
    //         // TODO: handle errors

    //     setLastQuery(query);
    // }

    // useEffect(() => {
    //     let t = setTimeout(() => {
    //         if (query.length >= 2 && query !== lastQuery)
    //             refreshSearchResults();
    //     }, 500);

    //     return () => clearTimeout(t);
    // }, [query]);

    // const handleSubmit = e => {
    //     e.preventDefault();
    // };

    return <Modal title={ title } isOpen={ isOpen } onClose={ onClose }>
        <style>{`
            .entry-form {
                padding: var(--pico-block-spacing-vertical) var(--pico-block-spacing-horizontal);
                display: flex;
                flex-direction: column;
                max-height: 50vh;
                gap: var(--pico-spacing);

                & > * { margin: 0; }

                & > div[role="search"] button { line-height: 0; }
                & > input { margin-bottom: 0; }
                & legend { margin-bottom: calc(var(--pico-spacing) * .25); }
                & label { width: 100%; }

                & .search-results {
                    border: 2px solid hsl(from var(--pico-color) h s l / 20%);
                    // min-height: 0;
                    // display: flex;
                    // flex-direction: column;

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
                <button className="green" title="Add a new food item">
                    <MdAdd />
                </button>
                <input
                    type="search" 
                    value={ searchString }
                    onChange={ e => setSearchString(e.target.value) }
                    placeholder="Search for food..." 
                />
                <button title="Refresh search results">
                    <MdRefresh />
                </button>
            </div>
            <fieldset className="food-search-results">
                <legend>Food</legend>
                <div className="search-results">
                    { loading ? <div className="loading" aria-busy="true">Loading...</div> : (<>
                        { searchResults.length === 0 ? <div>No results found.</div> : searchResults.map(food => (
                            <label key={ food.id }>
                                <div><input type="radio" name="food_id" value={ food.id } required /></div>
                                <div>
                                    { food.name }
                                    {  food.brand && <><br /><small>({ food.brand })</small></> }
                                    <br />
                                    <small>
                                        { food.protein }p - { food.carbs }c - { food.fat }f - { food.calories }kcal
                                    </small>
                                </div>
                            </label>
                        )) }
                    </>) }
                </div>
            </fieldset>
            <label>
                Servings
                <input type="number" min="0" step="0.01" value={ servings } onChange={ e => setServings(e.target.value) } required />
            </label>
            <input type="submit" value={ submitLabel } />
        </form>
    </Modal>;
}