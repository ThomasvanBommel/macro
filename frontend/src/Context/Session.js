import { useState, useEffect, createContext, useContext } from 'react';
import { handleLoginUser, handleLogoutUser, handleFetchSession } from '../api';
import NotificationContext from './Notification';

export function useSession() {
    const [session, setSession] = useState(null);
    const notifications = useContext(NotificationContext);

    const create = (username, password) => {
        if (isValid()) return;

        handleLoginUser(username, password)
            .then(setSession)
            .catch(err => notifications.add({
                heading: "Failed to login",
                details: err.message,
                type: "error"
            }));
    };

    const refresh = () => {
        handleFetchSession()
            .then(setSession)
            .catch(err => notifications.add({
                heading: "Failed to fetch session",
                details: err.message,
                type: "error"
            }));
    }

    const isValid = () => {
        if (!session) return false;

        if (new Date(session.expires) - new Date() > 0)
            return true;

        setSession(null);
        return false;
    }

    const destroy = () => {
        if (!isValid()) return;

        handleLogoutUser()
            .then(() => setSession(null))
            .catch(err => notifications.add({
                heading: "Failed to logout",
                details: err.message,
                type: "error"
            }));
    };

    useEffect(refresh, []);

    return { ...session, create, refresh, isValid, destroy };
}

const SessionContext = createContext(null);
export default SessionContext;