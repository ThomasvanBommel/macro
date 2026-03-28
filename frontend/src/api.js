import axios from 'axios';


export async function clearSession() {
    const res = await axios.post("/api/logout");

    if(res.status !== 200)
        throw new Error(`Logout failed: ${res.data.message}`);
}

export async function fetchSession() {
    try {
        const res = await axios.post("/api/login", {});

        if(res.status !== 200)
            console.error(`fetchSession !== 200: ${res.data.message}`);

        return res.status === 200 ? res.data : null;
    } catch (err) {
        console.error(err);
        return null;
    }
}

export async function fetchEntries(user_name, date) {
    try {
        const res = await axios.post("/api/entries", { user_name, date });

        if(res.status !== 200)
            console.error(`fetchEntries !== 200: ${res.data.message}`);

        return res.status === 200 ? res.data : [];
    } catch (err) {
        console.error(err);
        return [];
    }
}

// fetchFoods
export async function fetchFoods() {
    const scale = foods => foods.map(food => {
        food.calories /= 100;
        food.protein /= 100;
        food.carbs /= 100;
        food.fat /= 100;
        food.serving_count /= 100;

        return food;
    });

    try {
        const res = await axios.get("/api/foods");

        if(res.status !== 200)
            console.error(`fetchFoods !== 200: ${res.data.message}`);

        return res.status === 200 ? scale(res.data) : [];
    } catch (err) {
        console.error(err);
        return [];
    }
}

export async function createFood(foodData) {
    const data = { ...foodData };
    const scale = key => data[key] = Math.floor((data[key] ?? 0) * 100);

    scale("calories");
    scale("protein");
    scale("carbs");
    scale("fat");
    scale("serving_count");

    try {
        const res = await axios.post("/api/food", data);

        if(res.status !== 200)
            console.error(`createFood !== 200: ${res.data.message}`);

        return res.status === 200 ? true : false;
    } catch (err) {
        alert(`Failed to create food: ${err.response?.data?.message || err.message}`);
        return false;
    }
}