import { useState, useContext } from 'react';

import { handleLoginUser } from '../api'
import { NotificationContext, SessionContext } from '../Context';

// Renders a login form and handles submission.
export default function LoginForm() {
    const [loading, setLoading] = useState(false);

    const notifications = useContext(NotificationContext);
    const session = useContext(SessionContext);

    const handleSubmit = e => {
        e.preventDefault();

        const form = new FormData(e.currentTarget);
        setLoading(true);
        handleLoginUser(form.get('name'), form.get('password'))
            .then(session.refresh)
            .catch(err => notifications.add({ 
                heading: "Login failed",
                details: err.message,
                type: "error"
            }))
            .finally(() => setLoading(false));
    };

    if (session.changingState)
        return null;

    return (
        <form onSubmit={handleSubmit} style={{ maxWidth: "40rem", margin: "0 auto" }}>
            <label>
                Username:
                <input type="text" name="name" required />
            </label>
            <label>
                Password:
                <input type="password" name="password" autoComplete="current-password" required />
            </label>
            {  loading ? 
                <div role="button" aria-busy="true" disabled>Logging in...</div> : 
                <button type="submit">Login</button>
            }
        </form>
    );
}