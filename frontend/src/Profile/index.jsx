import { useState, useEffect } from 'react';
import axios from 'axios';
import "./index.css";
import { fetchEntries, fetchFoods } from '../api';
import { DateSelector, EntryList, SessionTimer } from '../Components';
// import DateSelector from '../Components/DateSelector';
// import EntryList from '../Components/EntryList';
import DateString from '../Classes/DateString';

export default function Profile({ session }) {
    const [date, setDate] = useState(DateString.today());

    return (
        <>  
            <div>
                <h2>Profile Page</h2>
                <p>Welcome to your profile, {session?.user_name}!</p>
            </div>

            <div style={{ display: "grid", placeItems: "center" }}>
                <DateSelector date={ date } setDate={ setDate } />
            </div>

            <EntryList username={ session?.user_name } date={ date } editable={ true } />
        </>
    )
}

function AddEntryForm({ onSuccess, entryDate }) {
    const initialState = {
        date: entryDate,
        meal_name: 'breakfast',
        food_id: undefined,
        servings: 1
    };

    const [formData, setFormData] = useState(initialState);
    const [foods, setFoods] = useState([]);
    const [loadingFoods, setLoadingFoods] = useState(true);
    const [foodModalOpen, setFoodModalOpen] = useState(false);

    useEffect(() => {
        fetchFoods()
            .then(setFoods)
            .catch(console.error)
            .finally(() => setLoadingFoods(false));
    }, []);

    const handleSubmit = e => {
        e.preventDefault();
        e.stopPropagation();

        axios.post("/api/entry", formData)
            .then(res => {
                if(res.status !== 200)
                    alert(`Failed to add entry: ${res.data.message}`);

                setFormData(initialState);
                onSuccess?.(res.data);
            })
            .catch(e => {
                alert(`An error occurred: ${e.message}${'\n' + e.response.data.error ?? ''}`);
                console.error(e);
            });
    }

    return (
        <form onSubmit={handleSubmit}>
            <div>
                <label htmlFor="meal_name">Meal:</label>
                <select id="meal_name" name="meal_name" value={ formData.meal_name } 
                    onChange={ e => setFormData(p => ({ ...p, meal_name: e.target.value })) } 
                    required>
                    <option value="breakfast">Breakfast</option>
                    <option value="lunch">Lunch</option>
                    <option value="dinner">Dinner</option>
                    <option value="snack">Snack</option>
                </select>
            </div>
            <div>

                

                { foodModalOpen ?
                    <AddFoodForm formData={formData} setFormData={setFormData} />
                : 
                    <>
                        <button className="green" onClick={ () => setFoodModalOpen(true) }>+Food</button>

                        { loadingFoods ? (
                            <p>Loading foods...</p>
                        ) : foods.length === 0 ? (
                            <p>No foods available.</p>
                        ) : (
                            <>
                                <label htmlFor="food">Food:</label>
                                <select id="food" name="food" value={ formData.food_id } 
                                    onChange={ e => setFormData(p => ({ ...p, food_id: +e.target.value })) }
                                    required>
                                        <option value="">Select a food</option>
                                    { foods.map(f => (
                                        <option value={ f.id }>{ f.name }</option>
                                    )) }
                                </select>
                            </>
                        ) }
                    </>
                }
            </div>
            <div>
                <label htmlFor="servings">Servings:</label>
                <input type="number" id="servings" name="servings" value={ formData.servings } 
                    onChange={ e => setFormData(p => ({ ...p, servings: +e.target.value }))} min="1"
                     required />
            </div>
            <button type="submit">Add Entry</button>
        </form>
    )
}

function AddFoodForm({ formData, setFormData }) {
    const emptyFood = {
        name: '',
        brand: '',
        calories: 0,
        carbs: 0,
        protein: 0,
        fat: 0,
        serving_size: ''
    };

    const food = formData.food ?? emptyFood;

    useEffect(() => {
        setFormData(p => ({
            ...p, 
            food: emptyFood
        }));
    }, []);

    const onChange = e => {
        const { name, value } = e.target;
        setFormData(p => ({
            ...p,
            food: {
                ...p.food,
                [name]: e.target.type === "number" ? +value : value
            }
        }));
    }

    return (
        <>
            <div>
                <label htmlFor="name">Name:</label>
                <input type="text" id="name" name="name" value={ food.name } 
                    onChange={ onChange } required />
            </div>
            <div>
                <label htmlFor="brand">Brand:</label>
                <input type="text" id="brand" name="brand" value={ food.brand } 
                    onChange={ onChange } />
            </div>
            <div>
                <label htmlFor="calories">Calories:</label>
                <input type="number" id="calories" name="calories" value={ food.calories } 
                    onChange={ onChange } min="0" required />
            </div>
            <div>
                <label htmlFor="carbs">Carbs (g):</label>
                <input type="number" id="carbs" name="carbs" value={ food.carbs } 
                    onChange={ onChange } min="0" required />
            </div>
            <div>
                <label htmlFor="protein">Protein (g):</label>
                <input type="number" id="protein" name="protein" value={ food.protein } 
                    onChange={ onChange } min="0" required />
            </div>
            <div>
                <label htmlFor="fat">Fat (g):</label>
                <input type="number" id="fat" name="fat" value={ food.fat } 
                    onChange={ onChange } min="0" required />
            </div>
            <div>
                <label htmlFor="serving_size">Serving Size:</label>
                <input type="text" id="serving_size" name="serving_size" value={ food.serving_size } 
                    onChange={ onChange } required />
            </div>
        </>
    )
}