import { useState, useEffect, useContext } from "react";
import { handleFetchFoodList } from "../api";
import Modal, { ModalHeader } from "./Modal";
import CreateFoodModal from "./CreateFoodModal";
import { NotificationContext } from '../Context';

// 
export default function FindFoodModal({ isOpen, onClose, onSelect }) {
    const [foods, setFoods] = useState([]);
    const [loadingFoods, setLoadingFoods] = useState(true);
    const [searchTerm, setSearchTerm] = useState("");

    const [createFoodModalOpen, setCreateFoodModalOpen] = useState(false);

    const notifications = useContext(NotificationContext);

    useEffect(() => {
        handleFetchFoodList()
            .then(setFoods)
            .catch(error => notifications.add({ 
                heading: "Failed to fetch food list",
                details: error.message,
                type: "error"
            }))
            .finally(() => setLoadingFoods(false));
    }, []);

    const handleSubmit = e => {
        e.preventDefault();
        e.stopPropagation();

        const formData = new FormData(e.target);
        const foodId = parseInt(formData.get("food"));

        const selectedFood = foods.find(food => food.id === foodId);
        onSelect?.(selectedFood);
    }

    if (createFoodModalOpen) {
        return <CreateFoodModal isOpen={ createFoodModalOpen } 
                                onClose={ () => setCreateFoodModalOpen(false) }
                                onSuccess={ (food) => {
                                    setFoods(prev => [...prev, food]);
                                    setCreateFoodModalOpen(false);
                                }} />
    }

    return (
        <Modal isOpen={ isOpen } onClose={ onClose }>
            <ModalHeader title={`Find Food Modal`} onClose={ onClose } />

            <div role="group">
                <input type="text" 
                       placeholder="Search for food..." 
                       value={ searchTerm } 
                       onChange={ (e) => setSearchTerm(e.target.value) } />
                <input type="button" 
                       value="New" 
                       className="green" 
                       onClick={ () => setCreateFoodModalOpen(true) } />
            </div>

            { loadingFoods ? (
                <article aria-busy="true"></article>
            ) : (
                <form style={{ display: "flex", flexDirection: "column", gap: "1rem" }}
                      onSubmit={ handleSubmit }>
                    { foods
                        .filter(food => food.name.toLowerCase().includes(searchTerm.toLowerCase()))
                        .map(food => (
                            <label key={"FFM-" + food.id} style={{ display: "flex", width: "100%", gap: "1rem" }}>
                                <div style= {{ display: "flex", alignItems: "center" }}>
                                    <input type="radio" name="food" value={food.id} required />
                                </div>
                                <div>
                                    { food.name } 
                                    <small style={{ opacity: 0.5, display: "flex", gap: "0.5rem" }}>
                                        <div>{ food.protein }p</div>
                                        <div>{ food.carbs }c</div>
                                        <div>{ food.fat }f</div>
                                        <div>{ food.calories }kcal</div>
                                        <div>({ food.serving_count } { food.serving_size })</div>
                                    </small>
                                </div>
                            </label>
                        )) }

                    <input type="submit" value="Done" />
                </form>
            ) }
        </Modal>
    );
}