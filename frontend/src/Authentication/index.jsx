import { Outlet } from 'react-router';
import React, { useState } from 'react'

export default function AuthLayout() {
    return (
        <div>
            <Outlet />
        </div>
    )
}

export function Login({ setIsAuthenticated }) {
  return (
    <form action="/api/login" method="POST">
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

export function Register({ setIsAuthenticated }) {
  const [password, setPassword] = useState('');
  const [confirm, setConfirm] = useState('');

  const handleSubmit = e => {
    e.preventDefault();

    const formData = new FormData(e.target);
    const data = Object.fromEntries(formData.entries());
    
    fetch('/api/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    }).then(response => {
      if (response.ok) {
        setPassword('');
        setConfirm('');
        e.target.reset();
        alert('Registration successful!');
        setIsAuthenticated(true);
      } else {
        response.text().then(text => alert(`Registration failed: ${text}`));
      }
    }).catch(error => {
      alert(`An error occurred: ${error.message}`);
    });
  }

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