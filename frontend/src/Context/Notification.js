import { useState, createContext } from 'react';

export function useNotifications() {
    const [notifications, setNotifications] = useState([]);

    const add = ({ heading, details, type, ttl=5000 }) => {
        const id = crypto.randomUUID();

        setNotifications(prev => {
            const n = {
                id, heading, details, type, ttl,
                expires: new Date(Date.now() + ttl)
            };

            const nextNotifications = [...prev, n];
            return nextNotifications;
        });

        if (ttl > 0) {
            setTimeout(() => remove(id), ttl);
        }
    };

    const remove = (id) => {
        setNotifications(notifications => notifications.filter(n => n.id !== id));
    };

    return { notifications, add, remove };
}

const NotificationContext = createContext(null);
export default NotificationContext;