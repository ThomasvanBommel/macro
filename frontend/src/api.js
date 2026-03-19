import axios from 'axios';

export async function clearSession() {
    try {
        const res = await axios.delete("/api/session");
        if(res.status !== 200) {
            console.log("clearSession:", res.data);
            return false;
        }

        document.cookie = "macro_session=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
        return true;
    } catch(err) {
        console.error(err);
        return false;
    }
}

export async function fetchSession() {
    try {
        const res = await axios.post("/api/session");
        if(res.status !== 200)
            console.log("fetchSession:", res.data);
        return res.status === 200 ? res.data : null;
    } catch(err) {
        if(err.status !== 401)
            console.error("fetchSession:", err);
        return null;
    }
};