import { useState, createContext } from 'react';

export function useNotifications() {
    const [notifications, setNotifications] = useState([]);

    const add = ({ heading, details, type }) => {
        const id = crypto.randomUUID();

        setNotifications(prev => {
            const nextNotifications = [...prev, { id, heading, details, type }];

            return nextNotifications;
        });
    };

    const remove = (id) => {
        setNotifications(notifications => notifications.filter(n => n.id !== id));
    };

    return { notifications, add, remove };
}

const NotificationContext = createContext(null);
export default NotificationContext;