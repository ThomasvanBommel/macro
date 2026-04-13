async function postJSON(url, body) {
    const res = await fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body)
    });

    const text = await res.text();
    const data = text ? JSON.parse(text) : null;
    if (!res.ok) {
        const error = new Error(data || "Unknown error");
        error.status = res.status;
        throw error;
    }

    return data;
}

/** Registers a user and returns the API message. */
export async function handleRegisterUser(name, password) {
    return await postJSON("/api/register", { name, password });
}

/** Logs a user in and returns the API message. */
export async function handleLoginUser(name, password) {
    return await postJSON("/api/login", { name, password });
}

/** Validates and returns the current session. */
export async function handleFetchSession() {
    return await postJSON("/api/validate-session");
}

/** Logs a user out and clears the local session cookie. */
export async function handleLogoutUser() {
    await postJSON("/api/logout");
}

/** Fetches entries for a user on a specific date. */
export async function handleFetchEntryList(name, datestr) {
    return await postJSON("/api/entries", { name: name, date: datestr.value });
}

/** Creates an entry for the active user. */
export async function handleCreateEntry(id, meal, date, servings) {
    return await postJSON("/api/entry", { 
        food_id: id, meal_name: meal, date: date, servings 
    });
}

/** Fetches the available food list. */
export async function handleFetchFoodList() {
    return await postJSON("/api/foods");
}

// Creates a food item for the active user. Returns the created food data.
export async function handleCreateFood(food) {
    return await postJSON("/api/food", food);
}

// Fetches the diary data for a user on a specific date, including entries and totals.
export async function handleGetDiary(name, date) {
    return await postJSON("/api/diary", { name: name, date: date });
}

// Fetches a list of foods matching the search query. Returns an array of food data.
export async function handleSearchFoods(query) {
    return await postJSON("/api/food/search", { query });
}

// Edits an existing entry by ID for the active user. Returns the updated entry data.
export async function handleEditEntry(id, foodID, mealName, date, servings) {
    return await postJSON("/api/entry/edit", { 
        id: id, food_id: foodID, meal_name: mealName, date, servings 
    });
}

// Deletes an entry by ID for the active user. Returns a success message.
export async function handleDeleteEntry(id) {
    return await postJSON("/api/entry/delete", { id });
}