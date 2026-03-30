import { useState, useEffect } from "react";
import Modal, { ModalHeader } from "./Modal";
import FindFoodModal from "./FindFoodModal";
import { createEntry } from "../api";

export default function CreateEntryModal({ isOpen, onClose, initialMeal, date, onSuccess }) {
    const [meal, setMeal] = useState(initialMeal ?? "breakfast");
    const [servings, setServings] = useState(1);
    const [food, setFood] = useState(null);

    const [loading, setLoading] = useState(false);
    const [foodModalOpen, setFoodModalOpen] = useState(false);

    const onFoodChoice = food => {
        setFoodModalOpen(false);
        setFood(food);
    }

    const handleSubmit = e => {
        e.preventDefault();
        e.stopPropagation();

        setLoading(true);

        try {
            createEntry({
                food_id: food.id,
                meal_name: meal,
                date: date.value,
                servings: servings
            })
            .then(res => {
                if (res) onSuccess?.(res);
            })
            .finally(() => setLoading(false));
        } catch (err) {
            alert(`Failed to create entry: ${err.response?.data?.message || err.message}`);
            setLoading(false);
        }
    }

    if (foodModalOpen) {
        return <FindFoodModal isOpen={ foodModalOpen } 
                              onClose={ () => setFoodModalOpen(false) } 
                              onSelect={ onFoodChoice } />
    }

    return (
        <Modal isOpen={ isOpen } onClose={ onClose }>
            <ModalHeader title={`Create Entry Modal`} onClose={ onClose } />

            <form onSubmit={ handleSubmit }>
                <fieldset>
                    <MealSelector meal={ meal } setMeal={ setMeal } />

                    <label>
                        Food
                        <div style={{ display: "flex", width: "100%", gap: "1rem", marginBottom: "1rem" }}>
                            <div style= {{ display: "flex", alignItems: "center" }}>
                                <input type="radio" 
                                    name="food_id" 
                                    value={ food?.id }
                                    onClick={ () => setFoodModalOpen(true) }
                                    checked={ !!food }
                                    required />
                            </div>
                            <div>
                                { food ? food.name : "Select a food..." } 
                                { food ? (
                                <small style={{ opacity: 0.5, display: "flex", gap: "0.5rem" }}>
                                    <div>{ food.protein }p</div>
                                    <div>{ food.carbs }c</div>
                                    <div>{ food.fat }f</div>
                                    <div>{ food.calories }kcal</div>
                                    <div>({ food.serving_count } { food.serving_size })</div>
                                </small>
                                ) : null }
                            </div>
                        </div>
                    </label>

                    <label>
                        Servings
                        <input type="number" 
                               name="servings" 
                               min={ 0 } 
                               step={ 0.5 } 
                               value={ servings } 
                               onChange={ (e) => setServings(e.target.value) } />
                    </label>

                    { loading ? (
                        <div role="button" aria-busy="true" disabled></div>
                    ) : (
                        <input type="submit" value="Add Entry" />
                    ) }
                </fieldset>
            </form>
        </Modal>
    );
}

function MealSelector({ meal, setMeal }) {
    return (
        <label>
            Meal
            <select name="meal_name" 
                    onChange={ (e) => setMeal(e.target.value) }
                    defaultValue={ meal }>
                <option value="breakfast">Breakfast</option>
                <option value="lunch">Lunch</option>
                <option value="dinner">Dinner</option>
                <option value="snack">Snack</option>
            </select>
        </label>
    );
}