import { useState, useContext } from "react";

import { NotificationContext } from "../Context";
import Modal, { ModalHeader } from "./Modal";
import { handleCreateFood } from "../api";

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
    const notifications = useContext(NotificationContext);
    const [loading, setLoading] = useState(false);

    const handleSubmit = e => {
        e.preventDefault();
        e.stopPropagation();

        setLoading(true);

        const form = new FormData(e.currentTarget);
        const data = {
            name: form.get("name"),
            brand: form.get("brand") || null,
            calories: parseFloat(form.get("calories")),
            carbs: parseFloat(form.get("carbs")),
            protein: parseFloat(form.get("protein")),
            fat: parseFloat(form.get("fat")),
            serving_count: parseFloat(form.get("serving_count")),
            serving_size: form.get("serving_size")
        };

        handleCreateFood(data)
            .then(onSuccess)
            .catch(err => notifications.add({
                heading: "Failed to create food",
                details: err.message,
                type: "error"
            }))
            .finally(() => setLoading(false));
    }

    const TextInput = ({ name, placeholder, required=true }) => 
        <label>
            { name.charAt(0).toUpperCase() + name.slice(1) }
            <input type="text" name={ name } placeholder={ placeholder } required={ required } />
        </label>;

    const NumberInput = ({ name, value=0, required=true }) =>
        <label>
            { name.charAt(0).toUpperCase() + name.slice(1) }
            <input type="number" name={ name } min="0" step="0.01" inputMode="decimal" 
                   defaultValue={ value } required={ required } />
        </label>;

    return (
        <Modal isOpen={ isOpen } onClose={ onClose }>
            <ModalHeader title={`Create Food Modal`} onClose={ onClose } />

            <form onSubmit={ handleSubmit }>
                <fieldset>
                    <TextInput name="name" placeholder="Enter food name..." />
                    <TextInput name="brand" placeholder="Enter food brand..." required={ false } />
                    <NumberInput name="calories" />

                    <div style={{ display: "flex", gap: "1rem" }}>
                        <NumberInput name="carbs" />
                        <NumberInput name="protein" />
                        <NumberInput name="fat" />
                    </div>

                    <div style={{ display: "flex", gap: "1rem" }}>
                        <NumberInput name="serving_count" value={ 1 } />
                        <TextInput name="serving_size" placeholder="cups, g, ml, etc..." />
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