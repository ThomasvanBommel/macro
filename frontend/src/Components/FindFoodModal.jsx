import { useState, useEffect } from "react";
import { fetchFoods } from "../api";
import Modal, { ModalHeader } from "./Modal";
import CreateFoodModal from "./CreateFoodModal";

export default function FindFoodModal({ isOpen, onClose, onSelect }) {
    const [foods, setFoods] = useState([]);
    const [loadingFoods, setLoadingFoods] = useState(true);
    const [searchTerm, setSearchTerm] = useState("");

    const [createFoodModalOpen, setCreateFoodModalOpen] = useState(false);

    useEffect(() => {
        fetchFoods()
            .then(setFoods)
            .catch(console.error)
            .finally(() => setLoadingFoods(false));
    }, []);

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
                <ul style={{ display: "flex", flexDirection: "column", gap: "1rem" }}>
                    { foods
                        .filter(food => food.name.toLowerCase().includes(searchTerm.toLowerCase()))
                        .map(food => (
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
                        )) }
                </ul>
            ) }
        </Modal>
    );
}