import { useEffect, useMemo, useRef, useState } from "react";
import { MdAdd, MdClose, MdRefresh } from "react-icons/md";

import { handleSearchFoods } from "../../api";

const MEAL_OPTIONS = [
    { value: "breakfast", label: "Breakfast" },
    { value: "lunch", label: "Lunch" },
    { value: "dinner", label: "Dinner" },
    { value: "snacks", label: "Snack" },
];

/**
 * Clean entry modal with debounced food search and a compact, scrollable result list.
 * Submit behavior is injected via the onSubmit prop to keep this component reusable.
 */
export default function EntryModal2({
    isOpen,
    onClose,
    onSubmit,
    onCreateFood,
    submitLabel = "Save Entry",
    title = "Add Entry",
    initialMeal = "breakfast",
    initialServings = 1,
    initialFood = null,
}) {
    const [meal, setMeal] = useState(initialMeal);
    const [servings, setServings] = useState(String(initialServings));
    const [query, setQuery] = useState(initialFood?.name ?? "");
    const [selectedFood, setSelectedFood] = useState(initialFood);
    const [results, setResults] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    const requestVersion = useRef(0);

    useEffect(() => {
        if (!isOpen) return;

        setMeal(initialMeal ?? "breakfast");
        setServings(String(initialServings ?? 1));
        setSelectedFood(initialFood ?? null);
        setQuery(initialFood?.name ?? "");
        setResults(initialFood ? [initialFood] : []);
        setError("");
    }, [isOpen, initialMeal, initialServings, initialFood]);

    useEffect(() => {
        if (!isOpen) return;

        const trimmed = query.trim();
        if (trimmed.length < 2) {
            setResults(selectedFood ? [selectedFood] : []);
            setLoading(false);
            return;
        }

        setLoading(true);
        setError("");

        const currentRequest = ++requestVersion.current;
        const timer = setTimeout(() => {
            handleSearchFoods(trimmed)
                .then(data => {
                    if (currentRequest !== requestVersion.current) return;

                    const searchResults = Array.isArray(data) ? data : [];
                    const hasSelected = selectedFood
                        ? searchResults.some(food => food.id === selectedFood.id)
                        : false;

                    setResults(
                        selectedFood && !hasSelected
                            ? [selectedFood, ...searchResults]
                            : searchResults
                    );
                })
                .catch(searchError => {
                    if (currentRequest !== requestVersion.current) return;
                    setError(searchError.message || "Food search failed.");
                })
                .finally(() => {
                    if (currentRequest !== requestVersion.current) return;
                    setLoading(false);
                });
        }, 350);

        return () => clearTimeout(timer);
    }, [query, selectedFood, isOpen]);

    const selectedFoodId = selectedFood?.id ?? null;

    const canSubmit = useMemo(() => {
        const numericServings = Number(servings);
        return !!selectedFood && Number.isFinite(numericServings) && numericServings > 0;
    }, [selectedFood, servings]);

    const refreshSearch = () => {
        if (query.trim().length < 2) return;

        requestVersion.current += 1;
        setLoading(true);
        setError("");

        handleSearchFoods(query.trim())
            .then(data => {
                const searchResults = Array.isArray(data) ? data : [];
                const hasSelected = selectedFood
                    ? searchResults.some(food => food.id === selectedFood.id)
                    : false;

                setResults(
                    selectedFood && !hasSelected
                        ? [selectedFood, ...searchResults]
                        : searchResults
                );
            })
            .catch(searchError => {
                setError(searchError.message || "Food search failed.");
            })
            .finally(() => setLoading(false));
    };

    const handleFormSubmit = event => {
        event.preventDefault();
        event.stopPropagation();

        const numericServings = Number(servings);
        if (!selectedFood) {
            setError("Please choose a food.");
            return;
        }

        if (!Number.isFinite(numericServings) || numericServings <= 0) {
            setError("Servings must be greater than 0.");
            return;
        }

        setError("");
        onSubmit?.({
            meal,
            servings: numericServings,
            food: selectedFood,
            foodId: selectedFood.id,
            query: query.trim(),
        });
    };

    return (
        <FormModal
            isOpen={isOpen}
            title={title}
            onClose={onClose}
            onSubmit={handleFormSubmit}
            submitLabel={submitLabel}
        >
            <label>
                Meal
                <select value={meal} onChange={event => setMeal(event.target.value)}>
                    {MEAL_OPTIONS.map(option => (
                        <option key={option.value} value={option.value}>
                            {option.label}
                        </option>
                    ))}
                </select>
            </label>

            <label>
                Food Search
                <div role="search" style={{ display: "flex", gap: "0.5rem" }}>
                    <button
                        type="button"
                        className="green"
                        style={{ lineHeight: 0 }}
                        title="Create a new food"
                        onClick={() => onCreateFood?.({ query: query.trim() })}
                    >
                        <MdAdd />
                    </button>

                    <input
                        type="search"
                        value={query}
                        placeholder="Type at least 2 characters..."
                        onChange={event => setQuery(event.target.value)}
                    />

                    <button
                        type="button"
                        style={{ lineHeight: 0 }}
                        title="Refresh search"
                        onClick={refreshSearch}
                    >
                        <MdRefresh />
                    </button>
                </div>
            </label>

            <div className="food-results" style={{ marginBottom: "1rem" }}>
                {loading ? <small aria-busy="true">Searching foods...</small> : null}

                {!loading && query.trim().length >= 2 && results.length === 0 ? (
                    <small>No foods found.</small>
                ) : null}

                {results.map(food => (
                    <details key={food.id} open={food.id === selectedFoodId}>
                        <summary>
                            <label style={{ margin: 0, display: "inline-flex", gap: "0.5rem", alignItems: "center" }}>
                                <input
                                    type="radio"
                                    name="food_id"
                                    checked={food.id === selectedFoodId}
                                    onChange={() => setSelectedFood(food)}
                                />
                                <strong>{food.name}</strong>
                            </label>
                        </summary>

                        <div style={{ marginTop: "0.5rem", opacity: 0.85 }}>
                            <small>
                                {food.protein}p • {food.carbs}c • {food.fat}f • {food.calories} kcal
                            </small>
                            <br />
                            <small>
                                Serving: {food.serving_count} {food.serving_size}
                            </small>
                        </div>
                    </details>
                ))}
            </div>

            <label>
                Servings
                <input
                    type="number"
                    min="0.01"
                    step="0.01"
                    value={servings}
                    onChange={event => setServings(event.target.value)}
                    required
                />
            </label>

            {error ? (
                <small style={{ color: "var(--pico-del-color)" }}>
                    {error}
                </small>
            ) : null}

            {!canSubmit ? (
                <small style={{ opacity: 0.8 }}>
                    Select a food and set servings above 0.
                </small>
            ) : null}
        </FormModal>
    );
}

export function FormModal({ isOpen, title, onClose, onSubmit, submitLabel = "Submit", children }) {
    useEffect(() => {
        if (!isOpen) return;

        const handleKeyDown = event => {
            if (event.key === "Escape") {
                event.preventDefault();
                onClose?.();
            }
        };

        window.addEventListener("keydown", handleKeyDown);
        return () => window.removeEventListener("keydown", handleKeyDown);
    }, [isOpen, onClose]);

    if (!isOpen) return null;

    return (
        <>
            <style>{`
                .entry-form-modal article {
                    width: min(44rem, calc(100vw - 2rem));
                    max-height: calc(100vh - 2rem);
                    display: flex;
                    flex-direction: column;
                }

                .entry-form-modal form {
                    display: flex;
                    flex-direction: column;
                    min-height: 0;
                    gap: 0.5rem;
                }

                .entry-form-modal .form-content {
                    min-height: 0;
                    overflow-y: auto;
                    padding-right: 0.25rem;
                }

                .entry-form-modal .food-results {
                    max-height: 14rem;
                    overflow-y: auto;
                    border: 1px solid hsl(from var(--pico-muted-border-color) h s l / 55%);
                    border-radius: var(--pico-border-radius);
                    padding: 0.5rem;
                }

                .entry-form-modal details {
                    margin: 0;
                    padding: 0.35rem 0;
                }
            `}</style>

            <dialog open className="entry-form-modal" onClick={event => event.target === event.currentTarget && onClose?.()}>
                <article onClick={event => event.stopPropagation()}>
                    <header style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                        <strong>{title}</strong>
                        <button type="button" title="Close" style={{ lineHeight: 0, padding: "0.35rem" }} onClick={() => onClose?.()}>
                            <MdClose />
                        </button>
                    </header>

                    <form onSubmit={onSubmit}>
                        <div className="form-content">{children}</div>
                        <input type="submit" value={submitLabel} />
                    </form>
                </article>
            </dialog>
        </>
    );
}
