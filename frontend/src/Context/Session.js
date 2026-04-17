import { useState, useEffect, createContext, useContext } from 'react';
import { handleLoginUser, handleLogoutUser, handleFetchSession } from '../api';
import NotificationContext, { useNotifications } from './Notification';

export function useSession() {
    const [session, setSession] = useState(null);
    const [changingState, setChangingState] = useState(false);
    const notifications = useNotifications();

    const create = (username, password) => {
        if (isValid()) return;

        setChangingState(true);
        handleLoginUser(username, password)
            .then(setSession)
            .catch(err => notifications.add({
                heading: "Failed to login",
                details: err.message,
                type: "error"
            }))
            .finally(() => setChangingState(false));
    };

    const refresh = () => {
        setChangingState(true);
        handleFetchSession()
            .then(setSession)
            .catch(err => notifications.add({
                heading: "Failed to fetch session",
                details: err.message,
                type: "error"
            }))
            .finally(() => setChangingState(false));
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

        setChangingState(true);
        handleLogoutUser()
            .then(() => setSession(null))
            .catch(err => notifications.add({
                heading: "Failed to logout",
                details: err.message,
                type: "error"
            }))
            .finally(() => setChangingState(false));
    };

    useEffect(refresh, []);

    return { ...session, create, refresh, isValid, destroy, changingState, notifications };
}

const SessionContext = createContext(null);
export default SessionContext;