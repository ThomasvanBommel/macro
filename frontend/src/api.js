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

        return res.status === 200 ? res.data.map(scaleEntryDown) : [];
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

    try {
        const res = await axios.post("/api/food", scaleFoodUp(data));

        if(res.status !== 200)
            console.error(`createFood !== 200: ${res.data.message}`);

        return res.status === 200 ? scaleFoodDown(res.data) : false;
    } catch (err) {
        alert(`Failed to create food: ${err.response?.data?.message || err.message}`);
        return false;
    }
}

export async function createEntry(entryData) {
    const data = { ...entryData };
    data.servings = Math.floor((data.servings ?? 0) * 100);

    try {
        const res = await axios.post("/api/entry", data);

        if(res.status !== 200)
            console.error(`createEntry !== 200: ${res.data.message}`);

        return res.status === 200 ? scaleEntryDown(res.data) : false;
    } catch (err) {
        alert(`Failed to create entry: ${err.response?.data?.message || err.message}`);
        return false;
    }
}

function scaleFoodUp(food) {
    food.calories = Math.floor(food.calories * 100);
    food.protein = Math.floor(food.protein * 100);
    food.carbs = Math.floor(food.carbs * 100);
    food.fat = Math.floor(food.fat * 100);
    food.serving_count = Math.floor(food.serving_count * 100);

    return food;
}

function scaleFoodDown(food) {
    food.calories = food.calories / 100;
    food.protein = food.protein / 100;
    food.carbs = food.carbs / 100;
    food.fat = food.fat / 100;
    food.serving_count = food.serving_count / 100;

    return food;
}

// function scaleEntryUp(entry) {
//     if (entry.food)
//         entry.food = scaleFoodUp(entry.food);
//     entry.servings = Math.floor(entry.servings * 100);

//     return entry;
// }

function scaleEntryDown(entry) {
    entry.food = scaleFoodDown(entry.food);
    entry.servings = entry.servings / 100;

    return entry;
}