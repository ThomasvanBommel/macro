import { useState, useContext } from 'react';

import { handleRegisterUser, handleLoginUser } from '../api'
import { NotificationContext, SessionContext } from '../Context';

// Renders a registration form and handles submission.
export default function RegisterForm() {
    const [loading, setLoading] = useState(false);

    const notifications = useContext(NotificationContext);
    const session = useContext(SessionContext);

    const handleSubmit = e => {
        e.preventDefault();

        const form = new FormData(e.currentTarget);
        const name = form.get('name');
        const password = form.get('password');

        setLoading(true);
        handleRegisterUser(name, password)
            .then(() => handleLoginUser(name, password))
            .then(session.refresh)
            .catch(error => notifications.add({ 
                heading: "Registration failed",
                details: error.message,
                type: "error"
            }))
            .finally(() => setLoading(false));
    }

    const handlePasswordChange = e => {
        const target = e.currentTarget;
        const p1 = target.form.password;
        const p2 = target.form.confirm;

        [p1, p2].forEach(e => {
            e.setAttribute("aria-invalid", p1.value !== p2.value);
            e.setCustomValidity(p1.value === p2.value ? "" : "Passwords do not match");
        });

        target.form.submit.disabled = p1.value !== p2.value;
    }

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
                <input type="password" 
                       name="password" 
                       onChange={handlePasswordChange} 
                       autoComplete="new-password"
                       required />
            </label>
            <label>
                Confirm :
                <input type="password" 
                       name="confirm" 
                       onChange={handlePasswordChange} 
                       autoComplete="new-password"
                       required />
            </label>
            {  loading ? 
                <div role="button" aria-busy="true" disabled>Registering...</div> : 
                <input type="submit" name="submit" value="Register" disabled={loading} />
            }
        </form>
    );
}