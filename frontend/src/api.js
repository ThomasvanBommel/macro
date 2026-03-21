import axios from 'axios';


export async function clearSession() {
    const res = await axios.post("/api/logout");

    if(res.status !== 200)
        throw new Error(`Logout failed: ${res.data.message}`);

    // document.cookie = "macro_session=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
}

export async function fetchSession() {
    try {
        const res = await axios.post("/api/login");

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