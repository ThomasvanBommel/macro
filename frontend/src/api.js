async function postJSON(url, body) {
    const res = await fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body)
    });

    const data = await res.json();
    if (data.error) throw new Error(data.error);

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

/** Session payload returned by session validation. */
/**
 * @typedef {Object} Session
 * @property {string} username
 * @property {string} token
 * @property {string} created
 * @property {string} expires
 */

/** Validates and returns the current session. */
export async function handleFetchSession() {
    return await postJSON("/api/validate-session");
}

/** Logs a user out and clears the local session cookie. */
export async function handleLogoutUser() {
    await postJSON("/api/logout");

    document.cookie = "macro_session=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
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