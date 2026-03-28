import { useState, useEffect } from "react";
import { fetchFoods } from "../api";
import Modal, { ModalHeader } from "./Modal";
import FindFoodModal from "./FindFoodModal";

export default function CreateEntryModal({ isOpen, onClose, initialMeal="breakfast" }) {
    const [meal, setMeal] = useState(initialMeal);
    const [servings, setServings] = useState(1);

    const [foodModalOpen, setFoodModalOpen] = useState(false);

    if (foodModalOpen) {
        return <FindFoodModal isOpen={ foodModalOpen } onClose={ () => setFoodModalOpen(false) } />
    }

    return (
        <Modal isOpen={ isOpen } onClose={ onClose }>
            <ModalHeader title={`Add Entry Form`} onClose={ onClose } />

            <form>
                <fieldset>
                    <MealSelector meal={ meal } setMeal={ setMeal } />

                    <label>
                        Food
                        <input type="button" 
                               value="🔍 Search..." 
                               onClick={ () => setFoodModalOpen(true) } />
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

                    <input type="submit" value="Add Entry" />
                </fieldset>
            </form>
        </Modal>
    );
}

function MealSelector({ meal, setMeal }) {
    return (
        <label>
            Meal
            <select name="meal_name" onChange={ (e) => setMeal(e.target.value) }>
                <option value="breakfast" selected={meal === "breakfast"}>
                    Breakfast
                </option>
                <option value="lunch" selected={meal === "lunch"}>
                    Lunch
                </option>
                <option value="dinner" selected={meal === "dinner"}>
                    Dinner
                </option>
                <option value="snack" selected={meal === "snack"}>
                    Snack
                </option>
            </select>
        </label>
    );
}