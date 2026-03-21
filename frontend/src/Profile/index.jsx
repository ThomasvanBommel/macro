import { useState, useEffect } from 'react';
import axios from 'axios';
import "./index.css";
import { createPortal } from 'react-dom';

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

function AddEntryForm({ onSuccess }) {
    const initialState = {
        date: new Date().toLocaleString().split(',')[0],
        meal: 'breakfast',
        food_id: undefined,
        servings: 1
    };

    const [formData, setFormData] = useState(initialState);
    const [foods, setFoods] = useState([]);
    const [loadingFoods, setLoadingFoods] = useState(true);
    const [foodModalOpen, setFoodModalOpen] = useState(false);
    const [foodUpdate, setFoodUpdate] = useState(null);

    useEffect(() => {
        axios.get("/api/food")
            .then(res => {
                if(res.data) {
                    console.log(res.data)
                    setFoods(res.data);
                }
                setLoadingFoods(false);
            })
            .catch(console.error);
    }, [foodUpdate]);

    const handleSubmit = e => {
        e.preventDefault();
        e.stopPropagation();

        axios.put("/api/entry", formData )
            .then(res => {
                if(res.status === 201) {
                    setFormData(initialState);
                    onSuccess?.();
                } else {
                    alert(`Failed to add entry: ${res.data.message}`);
                }
            })
            .catch(console.error);
    }

    return (
        <form onSubmit={handleSubmit}>
            <div>
                <label htmlFor="date">Date:</label>
                <input type="date" id="date" name="date" value={ formData.date } 
                    onChange={ e => setFormData(p => ({ ...p, date: e.target.value })) } required />
            </div>
            <div>
                <label htmlFor="meal">Meal:</label>
                <select id="meal" name="meal" value={ formData.meal } 
                    onChange={ e => setFormData(p => ({ ...p, meal: e.target.value })) } required>
                    <option value="breakfast">Breakfast</option>
                    <option value="lunch">Lunch</option>
                    <option value="dinner">Dinner</option>
                    <option value="snack">Snack</option>
                </select>
            </div>
            <div>
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
                <button onClick={ () => setFoodModalOpen(true) }>+Food</button>
                <ModalForm isOpen={ foodModalOpen }>
                    <AddFoodForm onSuccess={ () => { setFoodModalOpen(false); setLoadingFoods(true); setFoodUpdate(Date.now()); } } />
                </ModalForm>
            </div>
            <div>
                <label htmlFor="servings">Servings:</label>
                <input type="number" id="servings" name="servings" value={ formData.servings } 
                    onChange={ e => setFormData(p => ({ ...p, servings: +e.target.value }))} min="1"
                     required />
            </div>
            <button type="submit" disabled={ loadingFoods || foods.length===0 }>Add Entry</button>
        </form>
    )
}

function EntryList({ session }) {
    const [date, setDate] = useState(new Date().toLocaleString().split(',')[0]);
    const [entries, setEntries] = useState([]);
    const [loading, setLoading] = useState(true);
    const [entryUpdate, setEntryUpdate] = useState(null);
    const [EntryModalOpen, setEntryModalOpen] = useState(false);

    const addDays = n => {
        const d = new Date(date + " 00:00");
        d.setDate(d.getDate() + n);
        setDate(d.toLocaleString().split(',')[0]);
    }

    useEffect(() => {
        axios.get("/api/entry", { params: { UserID: +(session?.user_name), Date: date } })
            .then(res => {
                console.log(res.data);
                setEntries(res.data);
            })
            .catch(console.error)
            .finally(() => setLoading(false));
    }, [entryUpdate]);
    
    if(loading) return <p>Loading entries...</p>;

    return (
        <div>
            <h3>Your Food Entries</h3>
            <div>
                <button onClick={ () => { addDays(-1); setLoading(true); setEntryUpdate(Date.now()); } }>-1</button>
                <input type="date" value={ date } onChange={ e => setDate(e.target.value) } />
                <button onClick={ () => { addDays(1); setLoading(true); setEntryUpdate(Date.now()); } }>+1</button>
            </div>
            <div>
                <button onClick={ () => setEntryModalOpen(true) }>Insert New Entry</button>
            </div>
            <ModalForm isOpen={ EntryModalOpen }>
                <AddEntryForm onSuccess={ () => { setEntryModalOpen(false); setLoading(true); setEntryUpdate(Date.now()); } } />
            </ModalForm>
            <div>
                { !entries || entries.length === 0 ? (
                    <p>No entries for this date.</p>
                ) : (
                    <ul>
                        { entries.map(e => (
                            <div>
                                <li key={ e.id }>{ e.food.name } - { e.servings } serving(s) - { e.food.calories * e.servings } calories</li>
                            </div>
                        )) }
                    </ul>
                ) }
            </div>
        </div>
    )
}

function AddFoodForm({ onSuccess }) {
    const [name, setName] = useState('');
    const [brand, setBrand] = useState('');
    const [calories, setCalories] = useState(0);
    const [carbs, setCarbs] = useState(0);
    const [protein, setProtein] = useState(0);
    const [fat, setFat] = useState(0);
    const [servingSize, setServingSize] = useState(1);

    const handleSubmit = e => {
        e.preventDefault();
        e.stopPropagation();
        axios.put("/api/food", { name: name, brand: brand, calories: Number(calories), carbs: Number(carbs), 
            protein: Number(protein), fat: Number(fat), serving_size: Number(servingSize) })
            .then(res => {
                if(res.status === 201) {
                    setName('');
                    setBrand('');
                    setCalories(0);
                    setCarbs(0);
                    setProtein(0);
                    setFat(0);
                    setServingSize(1);
                    onSuccess?.();
                } else {
                    alert(`Failed to add food: ${res.data.message}`);
                }
            })
            .catch(console.error);
    };

    return (
        <form onSubmit={handleSubmit}>
            <div>
                <label htmlFor="name">Name:</label>
                <input type="text" id="name" name="name" value={ name } 
                    onChange={ e => setName(e.target.value) } required />
            </div>
            <div>
                <label htmlFor="brand">Brand:</label>
                <input type="text" id="brand" name="brand" value={ brand } 
                    onChange={ e => setBrand(e.target.value) } />
            </div>
            <div>
                <label htmlFor="calories">Calories:</label>
                <input type="number" id="calories" name="calories" value={ calories } 
                    onChange={ e => setCalories(e.target.value) } min="0" required />
            </div>
            <div>
                <label htmlFor="carbs">Carbs (g):</label>
                <input type="number" id="carbs" name="carbs" value={ carbs } 
                    onChange={ e => setCarbs(e.target.value) } min="0" required />
            </div>
            <div>
                <label htmlFor="protein">Protein (g):</label>
                <input type="number" id="protein" name="protein" value={ protein } 
                    onChange={ e => setProtein(e.target.value) } min="0" required />
            </div>
            <div>
                <label htmlFor="fat">Fat (g):</label>
                <input type="number" id="fat" name="fat" value={ fat } 
                    onChange={ e => setFat(e.target.value) } min="0" required />
            </div>
            <div>
                <label htmlFor="serving_size">Serving Size:</label>
                <input type="number" id="serving_size" name="serving_size" value={ servingSize } 
                    onChange={ e => setServingSize(e.target.value) } min="1" required />
            </div>
            <button type="submit">Add Food</button>
        </form>
    )
}

function ModalForm({ isOpen, children }) {
    if(!isOpen) return null;

    return createPortal(
        <div className="modal-form">
            <div className="modal-content" onClick={ e => e.stopPropagation() }>
                { children }
            </div>
        </div>, document.body
    );
}