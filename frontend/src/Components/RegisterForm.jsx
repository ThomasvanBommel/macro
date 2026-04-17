import { useState } from 'react';

import { handleRegisterUser } from '../api'

// Renders a registration form and handles submission.
export default function RegisterForm({
    onSuccess,
    onError
}) {
    const [loading, setLoading] = useState(false);

    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");

    const handleSubmit = e => {
        e.preventDefault();

        setLoading(true);
        handleRegisterUser(username, password)
            .then(res => {
                setUsername("");
                setPassword("");
                setConfirmPassword("");

                onSuccess(res);
            })
            .catch(onError)
            .finally(() => setLoading(false));
    }

    const formChange = ({ currentTarget: { elements: { password, confirm } } }) => {
        const valid = password.value && password.value === confirm.value;

        [password, confirm].forEach(e => {
            if (!password.value && !confirm.value)
                return e.removeAttribute("aria-invalid");

            e.setAttribute("aria-invalid", !valid);
            e.setCustomValidity(valid ? "" : "Passwords do not match");
        });
    }

    return (
        <form 
            onSubmit={ handleSubmit } 
            onChange={ formChange }
            style={{ maxWidth: "40rem", margin: "0 auto" }}>
            <label>
                Username
                <input
                    type="text"
                    name="name"
                    value={ username }
                    onChange={ e => setUsername(e.currentTarget.value) } 
                    required
                    autoFocus />
            </label>
            <label>
                Password
                <input
                    type="password"
                    name="password"
                    value={ password }
                    onChange={ e => setPassword(e.currentTarget.value) }
                    autoComplete="new-password"
                    aria-invalid={ null }
                    required />
            </label>
            <label>
                Confirm Password
                <input
                    type="password"
                    name="confirm"
                    value={ confirmPassword }
                    onChange={ e => setConfirmPassword(e.currentTarget.value) }
                    autoComplete="new-password"
                    aria-invalid={ null }
                    required />
            </label>
            <button
                type="submit" 
                disabled={ loading || !password || password !== confirmPassword || !username } 
                aria-busy={ loading }>
                {  loading ? "Registering..." : "Register" }
            </button>
        </form>
    );
}