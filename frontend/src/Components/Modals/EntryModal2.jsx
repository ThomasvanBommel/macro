import { useContext, useEffect, useState } from 'react';

import Modal, { ModalHeader } from "../Modal";
import { handleFetchFoodList } from "../../api"

export default function EntryModal2({ entry=null, isOpen }) {
    const [open, setOpen] = useState(isOpen);
    const [foodList, setFoodList] = useState(null);
    const [searchString, setSearchString] = useState("");
    const [food, setFood] = useState(entry?.food ?? null);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        setLoading(true);

        handleFetchFoodList()
            .then(setFoodList)
            .finally(() => setLoading(false));
            // handle errors
    }, []);

    return (
        <Modal isOpen={ open } onClose={ () => setOpen(false) }>
            <ModalHeader title={'Entry Modal'} onClose={ () => setOpen(false) } />
            <form onSubmit={ e => {
                e.preventDefault();
                e.stopPropagation();
            }}>
                <label>
                    Meal
                    <select name="meal_name" 
                            defaultValue={ entry?.meal_name ?? "breakfast" }>
                        <option value="breakfast">Breakfast</option>
                        <option value="lunch">Lunch</option>
                        <option value="dinner">Dinner</option>
                        <option value="snack">Snack</option>
                    </select>
                </label>

                <label style={{ width: "100%" }}>
                    Food
                    { food ? null :
                    <input type="text" 
                            value={searchString} 
                            placeholder="Search for food..."
                            onChange={e => setSearchString(e.target.value)} />
                    }
                    { loading ? <span aria-busy="true">Loading foods...</span> : 
                    <div style={{ display: "flex", 
                                    flexDirection: "column", 
                                    gap: "0.5rem",
                                    paddingLeft: "1rem" }}>
                    { foodList?.filter(f => f.name.toLowerCase().includes(searchString.toLowerCase())).map(f => {
                        return food && food.id !== f.id ? null : (
                        <label key={ f.id } style= {{ display: "flex", alignItems: "center", width: "100%" }}>
                            <input type="radio" name="food_id" onClick={ () => setFood(f) } checked={ food?.id === f.id } required />
                            <div style={{ flex: 1 }}>
                                { f.name } 
                                <small style={{ opacity: 0.5, display: "flex", gap: "0.5rem" }}>
                                    <div>{ f.protein }p</div>
                                    <div>{ f.carbs }c</div>
                                    <div>{ f.fat }f</div>
                                    <div>{ f.calories }kcal</div>
                                    <div>({ f.serving_count } { f.serving_size })</div>
                                </small>
                            </div>
                            { !food ? null : 
                            <button style={{ aspectRatio: "1",
                                                lineHeight: 0,
                                                margin: "0.2rem",
                                                padding: "0.75rem", }}
                                    className="red"
                                    onClick={ () => setFood(null) }>-</button> }
                        </label>
                        )
                    }) }
                    </div>}
                </label>

                <label>
                    Servings { food ? `( ${ food.serving_size } )` : null }
                    <input type="number" name="servings" min={ 0 } step={ 0.01 } defaultValue={ entry?.servings ?? 1 } />
                </label>

                <input type="submit" value="Submit" />
            </form>
        </Modal>
    );
}