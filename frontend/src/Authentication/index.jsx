import { useNavigate } from 'react-router';
import axios from 'axios';
import React, { useState } from 'react';

export function Login({ onSuccess }) {
    const navigate = useNavigate();

    const handleSubmit = e => {
        e.preventDefault();

        const form = new FormData(e.currentTarget);
        const payload = {
            name: String(form.get('name') || ''),
            password: String(form.get('password') || ''),
        };

        axios.post("/api/login", payload)
            .then(res => {
                if(res.status !== 200)
                    return alert(`Login failed: ${res.data.message}`);

                onSuccess?.(res.data).then(() => navigate('/profile'));
            })
            .catch(err => {
                console.log(err);
                alert(`An error occurred: ${err.message}${'\n' + err.response.data.error ?? ''}`);
            });
    };

    return (
        <form onSubmit={handleSubmit}>
            <div>
                <label htmlFor="name">Username:</label>
                <input type="text" id="name" name="name" required />
            </div>
            <div>
                <label htmlFor="password">Password:</label>
                <input type="password" id="password" name="password" required />
            </div>
            <button type="submit">Login</button>
        </form>
    )
}

export function Register({ onSuccess }) {
    const [password, setPassword] = useState('');
    const [confirm, setConfirm] = useState('');

    const navigate = useNavigate();
    const handleSubmit = e => {
        e.preventDefault();

        const form = new FormData(e.currentTarget);
        const payload = {
            name: String(form.get('name') || ''),
            password: String(form.get('password') || ''),
        };

        axios.post("/api/register", payload)
            .then(res => {
                if(res.status !== 200)
                    return alert(`Registration failed: ${res.data.message}`);
                
                onSuccess?.(res.data).then(() => navigate('/profile'));
            })
            .catch(err => {
                alert(`An error occurred: ${err.message}${'\n' + err.response.data.error ?? ''}`);
            });
    };

    return (
        <form onSubmit={handleSubmit}>
            <div>
                <label htmlFor="name">Username:</label>
                <input type="text" id="name" name="name" required />
            </div>
            <div>
                <label htmlFor="password">Password:</label>
                <input type="password" id="password" name="password" required 
                onChange={e => setPassword(e.target.value)} />
            </div>
            <div>
                <label htmlFor="confirm">Confirm :</label>
                <input type="password" id="confirm" name="confirm" required 
                onChange={e => setConfirm(e.target.value)} />
            </div>
            <button type="submit" disabled={password !== confirm}>Register</button>
            {password !== confirm && <span style={{ color: 'red' }}> Passwords do not match!</span>}
        </form>
    )
}