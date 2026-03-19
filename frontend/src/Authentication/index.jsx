import { useNavigate, Navigate } from 'react-router';
import axios from 'axios';
import React, { Outlet, useState } from 'react';

export function Login({ onSuccess }) {
    const navigate = useNavigate();

    const handleSubmit = e => {
        e.preventDefault();
        axios.post("/api/login", e.target, { headers: { 'Content-Type': 'application/json' } })
            .then(res => {
                if(res.status === 200) {
                    // alert('Login successful!');
                    onSuccess?.()
                        .then(() => navigate('/profile'));
                } else {
                    alert(`Login failed: ${res.data.message}`);
                }
            })
            .catch(err => {
                console.log(err);
                alert("An error occurred");
                // alert(`An error occurred: ${err.message}${'\n' + err.response.data.error ?? ''}`);
            });
    };

    return (
        <form onSubmit={handleSubmit}>
            <div>
                <label htmlFor="username">Username:</label>
                <input type="text" id="username" name="username" required />
            </div>
            <div>
                <label htmlFor="password">Password:</label>
                <input type="password" id="password" name="password" required />
            </div>
            <button type="submit">Login</button>
        </form>
    )
}

export function Register() {
    const [password, setPassword] = useState('');
    const [confirm, setConfirm] = useState('');

    const navigate = useNavigate();
    const handleSubmit = e => {
        e.preventDefault();
        axios.post("/api/register", e.target, { headers: { 'Content-Type': 'application/json' } })
            .then(res => {
                if(res.status === 201) {
                    alert('Registration successful! Please log in.');
                    e.target.reset();
                    navigate('/login');
                } else {
                    alert(`Registration failed: ${res.data.message}`);
                }
            })
            .catch(err => {
                alert(`An error occurred: ${err.message}${'\n' + err.response.data.error ?? ''}`);
            });
    };

    return (
        <form onSubmit={handleSubmit}>
        <div>
            <label htmlFor="username">Username:</label>
            <input type="text" id="username" name="username" required />
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