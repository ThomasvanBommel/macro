import { useState, useEffect } from "react";

import Modal from "./Modal";
import { handleCreateFood } from "../../api";

export default function FoodModal({
    isOpen,
    onClose,
    onSuccess
}) {
    const [loading, setLoading] = useState(false);
    const [foodName, setFoodName] = useState("");
    const [brand, setBrand] = useState("");
    const [calories, setCalories] = useState(0);
    const [carbs, setCarbs] = useState(0);
    const [protein, setProtein] = useState(0);
    const [fat, setFat] = useState(0);
    const [servingCount, setServingCount] = useState(1);
    const [servingSize, setServingSize] = useState("");

    // reset
    useEffect(() => {
        if (!isOpen) return;

        setFoodName("");
        setBrand("");
        setCalories(0);
        setCarbs(0);
        setProtein(0);
        setFat(0);
        setServingCount(1);
        setServingSize("");
        setLoading(false);
    }, [isOpen])

    const handleSubmit = e => {
        e.preventDefault();

        const foodData = {
            name: foodName,
            brand: brand || null,
            calories: calories,
            carbs: carbs,
            protein: protein,
            fat: fat,
            serving_count: servingCount,
            serving_size: servingSize
        };

        setLoading(true);
        handleCreateFood(foodData)
            .then(data => {
                if (data.error)
                    throw new Error(data.error);

                onSuccess(data);
            })
            .catch(console.error)
            .finally(() => setLoading(false));
    }

    if (!isOpen) return null;

    return <Modal title="Add Food" isOpen={ isOpen } onClose={ onClose }>
        <style>{`
            .food-form {
                padding: var(--pico-block-spacing-vertical) var(--pico-block-spacing-horizontal);

                & fieldset {
                    display: flex;
                    gap: 1rem;
                    margin: 0;
                }

                & input[type="submit"] { margin: 0; }
            }
        `}</style>

        <form className="food-form" onSubmit={ handleSubmit }>
            <label>
                Food Name
                <input 
                    type="text" 
                    name="name" 
                    placeholder="Enter food name..." 
                    value={foodName}
                    onChange={(e) => setFoodName(e.target.value)}
                    required />
            </label>
            <label>
                Brand
                <input 
                    type="text" 
                    name="brand" 
                    placeholder="Enter food brand..." 
                    value={brand}
                    onChange={(e) => setBrand(e.target.value)} />
            </label>
            <label>
                Calories (kcal)
                <input 
                    type="number" 
                    name="calories" 
                    min="0" 
                    step="0.01" 
                    inputMode="decimal" 
                    value={calories}
                    onChange={(e) => setCalories(parseFloat(e.target.value))}
                    required />
            </label>
            <fieldset>
                <label>
                    Protein (g)
                    <input 
                        type="number" 
                        name="protein" 
                        min="0" 
                        step="0.01" 
                        inputMode="decimal" 
                        value={protein}
                        onChange={(e) => setProtein(parseFloat(e.target.value))}
                        required />
                </label>
                <label>
                    Carbs (g)
                    <input 
                        type="number" 
                        name="carbs" 
                        min="0" 
                        step="0.01" 
                        inputMode="decimal" 
                        value={carbs}
                        onChange={(e) => setCarbs(parseFloat(e.target.value))}
                        required />
                </label>
                <label>
                    Fat (g)
                    <input 
                        type="number" 
                        name="fat" 
                        min="0" 
                        step="0.01" 
                        inputMode="decimal" 
                        value={fat}
                        onChange={(e) => setFat(parseFloat(e.target.value))}
                        required />
                </label>
            </fieldset>
            <fieldset>
                <label>
                    Serving Count
                    <input 
                        type="number" 
                        name="serving_count" 
                        min="0.01" 
                        step="0.01" 
                        value={servingCount}
                        onChange={(e) => setServingCount(parseFloat(e.target.value))}
                        inputMode="decimal" 
                        required />
                </label>
                <label>
                    Serving Size
                    <input 
                        type="text" 
                        name="serving_size" 
                        placeholder="cup, slice, g, ml, etc..." 
                        value={servingSize}
                        onChange={(e) => setServingSize(e.target.value)}
                        required />
                </label>
            </fieldset>
            { loading ? (
                <div role="button" aria-busy="true" disabled>Loading...</div>
            ) : (
                <input type="submit" value="Save" />
            ) }
        </form>
    </Modal>;
}