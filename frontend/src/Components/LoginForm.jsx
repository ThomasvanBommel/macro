import { useContext } from 'react';
import { handleLoginUser } from '../api'

import { NotificationContext, SessionContext } from '../Context';

// Renders a login form and handles submission.
export default function LoginForm() {
    const notifications = useContext(NotificationContext);
    const session = useContext(SessionContext);

    const handleSubmit = e => {
        e.preventDefault();

        const form = new FormData(e.currentTarget);
        
        handleLoginUser(form.get('name'), form.get('password'))
            .then(session.refresh)
            .catch(err => notifications.add({ 
                heading: "Login failed",
                details: err.message,
                type: "error"
            }));
    };

    return (
        <form onSubmit={handleSubmit}>
            <label>
                Username:
                <input type="text" id="name" name="name" required />
            </label>
            <label>
                Password:
                <input type="password" id="password" name="password" required />
            </label>
            <button type="submit">Login</button>
        </form>
    );
}