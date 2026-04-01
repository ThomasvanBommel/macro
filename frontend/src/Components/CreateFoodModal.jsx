import { useState } from "react";
import Modal, { ModalHeader } from "./Modal";

/**
 * Callback for when the user closes the modal without creating a food.
 * @callback onClose
 * @returns { void }
 */

/**
 * Callback for successful food creation. Receives the new food data as an argument.
 * @callback onSuccess
 * @param { object } food - The newly created food data.
 * @returns { void }
 */

/**
 * Modal component for creating a new food. Contains a form with fields for all necessary food 
 * information. On successful creation, calls the onSuccess callback with the new food data.
 * @component
 * @param { Object } props - The component props
 * @param { boolean } props.isOpen - Whether the modal is open or not
 * @param { onClose } props.onClose - Callback for when the user closes the modal without creating a 
 *                                    food
 * @param { onSuccess } props.onSuccess - Callback for successful food creation. Receives the new 
 *                                        food data as an argument.
 * @returns { JSX.Element } The rendered component
 */
export default function CreateFoodModal({ isOpen, onClose, onSuccess }) {
    const [loading, setLoading] = useState(false);
    const [data, setData] = useState({
        name: "",
        brand: "",
        calories: 0,
        carbs: 0,
        protein: 0,
        fat: 0,
        serving_size: "",
        serving_count: 1
    });

    const handleChange = e => {
        const { name, value } = e.target;
        setData(prev => ({ ...prev, [name]: value }));
    };

    const handleSubmit = e => {
        e.preventDefault();
        e.stopPropagation();

        setLoading(true);

        // createFood(data)
        //     .then(res => {
        //         if (res) onSuccess(res);
        //         setLoading(false);
        //     });
    }

    return (
        <Modal isOpen={ isOpen } onClose={ onClose }>
            <ModalHeader title={`Create Food Modal`} onClose={ onClose } />

            <form onSubmit={ handleSubmit }>
                <fieldset>
                    <label>
                        Name
                        <input type="text" 
                               name="name" 
                               placeholder="Enter food name..." 
                               value={ data.name } 
                               onChange={ handleChange } 
                               required />
                    </label>

                    <label>
                        Brand
                        <input type="text" 
                               name="brand" 
                               placeholder="Enter food brand..." 
                               value={ data.brand } 
                               onChange={ handleChange } />
                    </label>

                    <label>
                        Calories (kcal)
                        <input type="number" 
                            name="calories" 
                            min="0" 
                            step="0.01"
                            inputMode="decimal"
                            value={ data.calories } 
                            onChange={ handleChange } 
                            required />
                    </label>

                    <div style={{ display: "flex", gap: "1rem" }}>
                        <label>
                            Carbs (g)
                            <input type="number" 
                                name="carbs" 
                                min="0" 
                                step="0.01"
                                inputmode="decimal"
                                value={ data.carbs } 
                                onChange={ handleChange } 
                                required />
                        </label>

                        <label>
                            Protein (g)
                            <input type="number" 
                                name="protein" 
                                min="0" 
                                step="0.01"
                                inputmode="decimal"
                                value={ data.protein } 
                                onChange={ handleChange } 
                                required />
                        </label>

                        <label>
                            Fat (g)
                            <input type="number" 
                                name="fat" 
                                min="0"
                                step="0.01"
                                inputmode="decimal"
                                value={ data.fat } 
                                onChange={ handleChange } 
                                required />
                        </label>
                    </div>

                    <div style={{ display: "flex", gap: "1rem" }}>
                        <label>
                            Serving Count
                            <input type="number" 
                                name="serving_count" 
                                placeholder="" 
                                min="0"
                                step="0.01"
                                inputmode="decimal"
                                value={ data.serving_count } 
                                onChange={ handleChange }
                                required />
                        </label>

                        <label>
                            Serving Size
                            <input type="text" 
                                name="serving_size" 
                                placeholder="cups, g, ml, etc..." 
                                value={ data.serving_size } 
                                onChange={ handleChange }
                                required />
                        </label>
                    </div>

                    { loading ? (
                        <div role="button" disabled aria-busy="true"></div>
                    ) : (
                        <input type="submit" value="Create Food" />
                    ) }
                </fieldset>
            </form>
        </Modal>
    );
}