import { useState, useEffect, createContext, useContext } from 'react';
import { handleFetchSession, handleLogoutUser } from '../api';

const SessionContext = createContext();

export function SessionProvider({ children }) {
    const [username, setUsername] = useState(null);
    const [expires, setExpires] = useState(null);
    const [notifications, setNotifications] = useState([]);

    // On component mount, initialize the session by fetching the existing session data.
    useEffect(() => init(true), []);

    // Initializes the session by fetching the current session data from the server.
    const init = (ignoreError=false) => {
        handleFetchSession()
            .then(res => {
                setUsername(res.username);
                setExpires(res.expires);
            })
            .catch(err => {
                if (ignoreError) return;
                setUsername(null);
                setExpires(null);
                add({ 
                    heading: "Error fetching session",
                    details: err.message || "Please log in again.",
                    type: "error",
                    ttl: 0
                 });
            });
    };

    // Destroys the session by logging out the user and clearing session data.
    const end = () => {
        handleLogoutUser()
            .then(() => {
                setUsername(null);
                setExpires(null);
            })
            .catch(err => {
                add({ 
                    heading: "Error logging out",
                    details: err.message || "Please try again.",
                    type: "error",
                    ttl: 0
                 });
            });
    };

    // Returns the time-to-live (TTL) of the session in milliseconds.
    const ttl = () => {
        if (!expires) return 0;
        return new Date(expires) - new Date();
    }

    // Checks if the session is valid (i.e., has a username and has not expired).
    const isValid = () => username && ttl() > 0;

    // Adds a new notification with the given heading, details, type, and time-to-live (TTL) in ms.
    const add = ({ heading, details, type, ttl=5000 }) => {
        const id = crypto.randomUUID();
        const n = { id, heading, details, type, ttl, expires: new Date(Date.now() + ttl) };
        setNotifications(prev => [...prev, n]);

        if (ttl > 0) setTimeout(() => remove(id), ttl);
    };

    // Removes a notification by its ID.
    const remove = (id) => setNotifications(n => n.filter(n => n.id !== id));

    return <SessionContext.Provider value={{
        username, expires, init, end, ttl, isValid, 
        notifications: {
            list: notifications,
            add,
            remove
        }
    }}>{ children }</SessionContext.Provider>;
}

export const useSession = () => useContext(SessionContext);