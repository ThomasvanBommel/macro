import { useContext } from 'react';

import { handleRegisterUser, handleLoginUser } from '../api'
import { NotificationContext, SessionContext } from '../Context';

//
export default function RegisterForm() {
    const notifications = useContext(NotificationContext);
    const session = useContext(SessionContext);

    const handleSubmit = e => {
        e.preventDefault();

        const form = new FormData(e.currentTarget);
        const name = form.get('name');
        const password = form.get('password');

        handleRegisterUser(name, password)
            .then(() => handleLoginUser(name, password))
            .then(session.refresh)
            .catch(error => notifications.add({ 
                heading: "Registration failed",
                details: error.message,
                type: "error"
            }));
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

    return (
        <form onSubmit={handleSubmit}>
            <label>
                Username:
                <input type="text" name="name" required />
            </label>
            <label>
                Password:
                <input type="password" name="password" onChange={handlePasswordChange} required />
            </label>
            <label>
                Confirm :
                <input type="password" name="confirm" onChange={handlePasswordChange} required />
            </label>
            <input type="submit" name="submit" value="Register" disabled />
        </form>
    );
}