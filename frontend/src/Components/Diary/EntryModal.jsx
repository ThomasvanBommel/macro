import { useState, useEffect } from 'react';
import { MdAdd, MdClose, MdRefresh } from "react-icons/md";

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

    useEffect(() => {
        if (!isOpen) return;

        setMeal(initialMeal ?? "breakfast");
        setFood(initialFood ?? null);
        setServings(initialServings ?? 1);
        setSearchString(initialFood?.name ?? "");
        setSearchResults(initialFood ? [initialFood] : []);
        setError("");
    }, [isOpen, initialMeal, initialFood, initialServings]);
    


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

                & > div[role="search"] {

                    & > button { line-height: 0; }
                }

                & > input { margin-bottom: 0; }
                & label { width: 100%; }

                & .search-results {
                    // background-color:red;
                    // padding: 0.5rem;
                    border: 2px solid hsl(from var(--pico-color) h s l / 20%);

                    & label {
                        display: flex;
                        gap: 0.5rem;
                        padding: 0.25rem;
                        margin: 0;
                    }

                    & label:hover {
                        cursor: pointer;
                        background-color: hsl(from var(--pico-color) h s l / 5%);
                    }

                    & label:active {
                        background-color: hsl(from var(--pico-color) h s l / 10%);
                    }

                    & input[type="radio"] {
                        margin:0;
                    }
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
            <label>
                Food
                <div className="search-results">
                    <label>
                        <div>
                            <input type="radio" />
                        </div>
                        <div>
                            Apple
                            <div>100p - 50c - 200f - 666kcal</div>
                        </div>
                    </label>
                </div>
            </label>
            <label>
                Servings
                <input type="number" min="0" step="0.01" value={ servings } onChange={ e => setServings(e.target.value) } required />
            </label>
            <input type="submit" value={ submitLabel } />
        </form>
    </Modal>;
}