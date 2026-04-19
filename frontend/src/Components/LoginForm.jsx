import { useState, useEffect } from 'react';

import { handleLoginUser } from '../api'

// Renders a login form and handles submission.
export default function LoginForm({
    onError,
    onSuccess
}) {
    const [loading, setLoading] = useState(false);
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

    const handleSubmit = e => {
        e.preventDefault();

        setLoading(true);
        handleLoginUser(username, password)
            .then(res => {
                setUsername("");
                setPassword("");

                onSuccess(res);
            })
            .catch(onError)
            .finally(() => setLoading(false));
    };

    return (
        <form onSubmit={handleSubmit} style={{ maxWidth: "40rem", margin: "0 auto" }}>
            <label>
                Username:
                <input
                    type="text" 
                    name="name"
                    value={ username }
                    onChange={ e => setUsername(e.currentTarget.value) }
                    required
                    autoFocus />
            </label>
            <label>
                Password:
                <input
                    type="password"
                    name="password"
                    value={ password }
                    onChange={ e => setPassword(e.currentTarget.value) }
                    autoComplete="current-password"
                    required />
            </label>
            <button
                type="submit"
                disabled={ loading || !username || !password }
                aria-busy={ loading }>
                { loading ? "Logging in..." : "Login" }
            </button>
        </form>
    );
}