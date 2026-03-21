import { useState, useEffect } from 'react';
import axios from 'axios';
import "./index.css";
import { ModalForm } from '../Modal';
import { fetchEntries, fetchFoods } from '../api';

export default function Profile({ session }) {
    const [timeLeft, setTimeLeft] = useState("-");
    const expiry = new Date(session?.expires);

    // Update time left every second
    useEffect(() => {
        const interval = setInterval(() => {
            const diff = expiry - new Date();
            setTimeLeft(diff > 0 ? 
                `${Math.floor(diff/60000)}m ${Math.floor((diff%60000)/1000)}s` : "Expired");
        }, 1000);
        return () => clearInterval(interval);
    }, []);

    return (
        <>  
            <div>
                <h2>Profile Page</h2>
                <p>Welcome to your profile, {session?.user_name}!</p>
                <p>Session expires in: {timeLeft} ({new Date(session?.expires).toLocaleString()})</p>
            </div>

            <EntryList session={ session } />
        </>
    )
}

function EntryList({ session }) {
    // const [date, setDate] = useState(new Date().toLocaleString().split(',')[0]);
    const toInputDate = d => new Date(d).toISOString().slice(0, 10);
    const [date, setDate] = useState(toInputDate(new Date()));
    const [entries, setEntries] = useState([]);
    const [loading, setLoading] = useState(true);
    const [entryUpdate, setEntryUpdate] = useState(null);
    const [EntryModalOpen, setEntryModalOpen] = useState(false);

    const reloadEntries = () => {
        setLoading(true);
        setEntryUpdate(Date.now());
    }

    // Helper to add days to current date and trigger entry reload
    const addDays = n => {
        // const d = new Date(date + " 00:00");
        // d.setDate(d.getDate() + n);
        // setDate(d.toLocaleString().split(',')[0]);
        const d = new Date(date + "T00:00:00");
        d.setDate(d.getDate() + n);
        setDate(toInputDate(d));
        reloadEntries();
    }

    // Fetch entries for the selected date
    useEffect(() => {
        fetchEntries(session.user_name, date)
            .then(setEntries)
            .catch(console.error)
            .finally(() => setLoading(false));
    }, [date, entryUpdate, session.user_name]);
    // }, [entryUpdate]);
    
    return (
        <div>
            <h3>Your Food Entries</h3>
            <div id="date-nav">
                <div className="spacer"></div>
                <input type="submit" value="Previous Day" onClick={ () => { addDays(-1) } } />
                <div>
                    {/* <input type="date" value={ date } onChange={ e => setDate(e.target.value) } /> */}
                    <input type="date"
                        value={date}
                        onChange={e => {
                            setLoading(true);
                            setDate(e.target.value);
                        }} />
                </div>
                <input type="submit" value="Next Day" onClick={ () => { addDays(1) } } />
                <div className="spacer"></div>

            </div>

            <div className="green">
                <input type="submit" 
                    value="Insert New Entry" 
                    onClick={ () => setEntryModalOpen(true) } />
            </div>

            <ModalForm isOpen={ EntryModalOpen } onClose={ () => setEntryModalOpen(false) }>
                <AddEntryForm onSuccess={ entry => { 
                    setEntryModalOpen(false); 
                    setEntries(prev => [...prev, entry]);
                    // reloadEntries(); 
                } } entryDate={ date } />
            </ModalForm>

            { loading ? <p>Loading entries...</p> : 
                <> 
                <div className="totals">
                    <div className="calories">Calories: { entries.reduce((sum, e) => sum + e.food.calories * e.servings, 0) }</div>
                    <div className="carbs">Carbs: { entries.reduce((sum, e) => sum + e.food.carbs * e.servings, 0) }g</div>
                    <div className="protein">Protein: { entries.reduce((sum, e) => sum + e.food.protein * e.servings, 0) }g</div>
                    <div className="fat">Fat: { entries.reduce((sum, e) => sum + e.food.fat * e.servings, 0) }g</div>
                </div>
                {
                    ["breakfast", "lunch", "dinner", "snack"].map(meal => {
                        const mealEntries = entries.filter(e => e.meal_name === meal);
                        if(mealEntries.length === 0) return null;

                        return (
                            <div className={ meal + " meal" }>
                                <h4>{ meal.charAt(0).toUpperCase() + meal.slice(1) }</h4>
                                <table>
                                    <thead>
                                        <tr>
                                            <th>Food</th>
                                            <th>Calories</th>
                                            <th>Carbs (g)</th>
                                            <th>Protein (g)</th>
                                            <th>Fat (g)</th>
                                            <th>Servings</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        { mealEntries.map(e => (
                                            <tr key={ e.id }>
                                                <td>{ e.food.name }</td>
                                                <td>{ e.food.calories }</td>
                                                <td>{ e.food.carbs }</td>
                                                <td>{ e.food.protein }</td>
                                                <td>{ e.food.fat }</td>
                                                <td>{ e.servings }</td>
                                            </tr>
                                        )) }
                                    </tbody>
                                    <tfoot>
                                        <tr>
                                            <th>Total</th>
                                            <th>{ mealEntries.reduce((sum, e) => sum + e.food.calories * e.servings, 0) }</th>
                                            <th>{ mealEntries.reduce((sum, e) => sum + e.food.carbs * e.servings, 0) }</th>
                                            <th>{ mealEntries.reduce((sum, e) => sum + e.food.protein * e.servings, 0) }</th>
                                            <th>{ mealEntries.reduce((sum, e) => sum + e.food.fat * e.servings, 0) }</th>
                                            <th>{ mealEntries.reduce((sum, e) => sum + e.servings, 0) }</th>
                                        </tr>
                                    </tfoot>
                                </table>
                            </div>
                        )  
                    })
                } </>
            }
        </div>
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