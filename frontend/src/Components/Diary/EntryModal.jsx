import { useState } from 'react';
import { MdClose } from "react-icons/md";

export default function EntryModal({ open, meal, entry, close }) {
    window.onkeydown = e => {
        e.stopImmediatePropagation();
        if (e.key === "Escape")
            close();
    }

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
        <dialog open={ open } onClick={ close } className="modal">
            <article
                onClick={ e => e.stopPropagation() }
                style={{ display: "flex", flexDirection: "column" }}>
                <header style={{ display: "flex", justifyContent: "space-between" }}>
                    <div>{ entry ? "Edit" : "Add" } Entry</div>
                    <div onClick={ close }
                         title="Close"
                         style={{ lineHeight: 0, padding: "0.25rem" }}
                         className="close-button">
                        <MdClose />
                    </div>
                </header>
                <form style={{ display: "flex", flexDirection: "column", flex: 1, minHeight: 0 }}>
                    <section style={{ overflowY: "auto", flex: 1, minHeight: 0 }}>
                        <label>
                            Meal
                            <select value={ entry?.meal_name ?? meal ?? "breakfast" } required>
                                <option value="breakfast">Breakfast</option>
                                <option value="lunch">Lunch</option>
                                <option value="dinner">Dinner</option>
                                <option value="snacks">Snack</option>
                            </select>
                        </label>
                        <label>
                            Food
                            <input type="text" required />
                        </label>
                        <label>
                            Servings
                            <input type="number" min="0" step="0.1" value="1" required />
                        </label>
                    </section>
                    <input type="submit" value="Submit" style={{ margin: 0 }} />
                </form>
            </article>
        </dialog>
    </>;
}