import axios from 'axios';
import { DateString } from './Classes'

/**
 * Handles the creation of a user via the registration endpoint.
 * @param { string } name - The name of the user.
 * @param { string } password - The password of the user.
 * @returns { Promise<string> } - A promise that resolves to the success message from the server.
 * @throws { Error } - Throws an error if the registration fails or if the server returns an error 
 *                     message.
 */
export async function handleRegisterUser(name, password) {
    const { message, error } = await axios.post("/api/register", { name, password });

    if (error) throw new Error(error);
    return message;
}

/**
 * Handles the login of a user via the login endpoint.
 * @param { string } name - The name of the user.
 * @param { string } password - The password of the user.
 * @returns { Promise<string> } - A promise that resolves to the success message from the server.
 * @throws { Error } - Throws an error if the login fails or if the server returns an error message.
 */
export async function handleLoginUser(name, password) {
    const { message, error } = await axios.post("/api/login", { name, password });

    if (error) throw new Error(error);
    return message;
}

/**
 * Session object containing user session information.
 * @typedef { Object } Session
 * @property { string } username - The name of the user associated with the session.
 * @property { string } token - The unique identifier for the session.
 * @property { string } created - The timestamp when the session was created.
 * @property { string } expires - The timestamp when the session expires.
 */

/**
 * Handles the validation of a user session via the validate-session endpoint.
 * @returns { Promise<Session> } - A promise that resolves to the session data if the session is 
 *                                valid.
 * @throws { Error } - Throws an error if the session validation fails or if the server returns an 
 *                     error message.
 */
export async function handleFetchSession() {
    const { error, ...session } = await axios.post("/api/validate-session");

    if (error) throw new Error(error);
    return session;
}

/**
 * Handles the logout of a user via the logout endpoint.
 * @returns { Promise<string> } - A promise that resolves to the success message from the server.
 * @throws { Error } - Throws an error if the logout fails or if the server returns an error message
 */
export async function handleLogoutUser() {
    const { message, error } = await axios.post("/api/logout");

    if (error) throw new Error(error);
    document.cookie = "macro_session=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";

    return message;
}

/**
 * Food object containing information about a food item.
 * @typedef { Object } Food
 * @property { number } id - The unique identifier for the food item.
 * @property { string } name - The name of the food item.
 * @property { string } brand - The brand of the food item.
 * @property { string } created - The timestamp when the food item was created.
 * @property { string } username - The name of the user who created the food item.
 * @property { number } calories - The calorie content of the food item.
 * @property { number } carbs - The carbohydrate content of the food in grams.
 * @property { number } protein - The protein content of the food in grams.
 * @property { number } fat - The fat content of the food in grams.
 * @property { string } serving_size - The serving size description for the food item.
 * @property { number } serving_count - The number of servings in the food item.
 */

/**
 * Entry object containing information about an entry.
 * @typedef { Object } Entry
 * @property { number } id - The unique identifier for the entry.
 * @property { string } username - The name of the user associated with the entry.
 * @property { Food } food - The food associated with the entry.
 * @property { string } mealname - The name of the meal associated with the entry.
 * @property { string } date - The date associated with the entry.
 * @property { number } servings - The number of servings associated with the entry.
 * @property { string } created - The timestamp when the entry was created.
 */

/**
 * Handles the fetching of the entry list for a specific user and date via the entries endpoint.
 * @param { string } name - The name of the user whose entries are to be fetched.
 * @param { DateString } datestr - The date for which the entries are to be fetched.
 * @returns { Promise<Entry[]> } - A promise that resolves to an array of entry objects for the 
 *                                 specified user and date.
 * @throws { Error } - Throws an error if the fetching of entries fails or if the server returns an 
 *                     error message.
 */
export async function handleFetchEntryList(name, datestr) {
    const { error, ...entries } = await axios.post("/api/entries", { 
        name: name, date: datestr.value 
    });

    if (error) throw new Error(error);
    return entries;
}
