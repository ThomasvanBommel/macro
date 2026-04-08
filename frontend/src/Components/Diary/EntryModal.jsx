import { useState, useEffect } from 'react';
import { MdClose, MdRefresh } from "react-icons/md";


import { handleSearchFoods } from "../../api";

export default function EntryModal({ open, defaultMeal, entry, close }) {
    const [meal, setMeal] = useState(entry?.meal_name ?? defaultMeal ?? "breakfast");
    const [food, setFood] = useState(entry?.food ?? null);
    
    const [query, setQuery] = useState("");
    const [lastQuery, setLastQuery] = useState("");
    const [lastQueryTime, setLastQueryTime] = useState(0);

    window.onkeydown = e => {
        e.stopImmediatePropagation();
        if (e.key === "Escape")
            close();
    }

    useEffect(() => {
        if (query.length > 2 && Date.now() - lastQueryTime > 500) {
            handleSearchFoods(query).then(setFood);
            setLastQueryTime(Date.now());
        }
    }, [query]);


    return (
        <FormModal title={ entry ? "Edit Entry" : "Add Entry" } isOpen={ open } onClose={ close }>
            <label>
                Meal
                <select value={ meal } onChange={ e => setMeal(e.target.value) } required>
                    <option value="breakfast">Breakfast</option>
                    <option value="lunch">Lunch</option>
                    <option value="dinner">Dinner</option>
                    <option value="snacks">Snack</option>
                </select>
            </label>
            <label>
                Food
                <div role="search">
                    <input type="search" 
                        value={ query } 
                        onChange={ e => setQuery(e.target.value) }
                        placeholder="Search for food..." />
                    <button><MdRefresh /></button>
                </div>
            </label>
            <label>
                Servings
                <input type="number" min="0" step="0.1" value="1" required />
            </label>
        </FormModal>
    );
}

function FormModal({ title, isOpen, onClose, onSubmit, children }) {
    window.onkeydown = e => {
        e.stopImmediatePropagation();
        if (e.key === "Escape")
            onClose?.();
    };

    return <>
        <style>{`
            .modal .close-button:hover {
                cursor: pointer;
                background-color: hsl(from var(--pico-color) h s l / 5%);
                border-radius: 0.25rem;
            }

            .modal .close-button:active {
                background-color: hsl(from var(--pico-color) h s l / 10%);
            }
        `}</style>
        <dialog open={ isOpen } onClick={ onClose } className="modal">
            <article
                onClick={ e => e.stopPropagation() }
                style={{ display: "flex", flexDirection: "column" }}>
                <header style={{ display: "flex", justifyContent: "space-between" }}>
                    <div>{ title }</div>
                    <div onClick={ onClose }
                         title="Close"
                         style={{ lineHeight: 0, padding: "0.25rem" }}
                         className="close-button">
                        <MdClose />
                    </div>
                </header>
                <form style={{ display: "flex", flexDirection: "column", flex: 1, minHeight: 0 }}
                      onSubmit={ onSubmit }>
                    <section style={{ overflowY: "auto", flex: 1, minHeight: 0 }}>
                        { children }
                    </section>
                    <input type="submit" value="Submit" style={{ margin: 0 }} />
                </form>
            </article>
        </dialog>
    </>;
}