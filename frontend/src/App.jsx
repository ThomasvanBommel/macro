import React, { useState } from 'react'

import './App.css'
import Header from './Header'

function AuthenticatedPage() {
  const [isLogin, setIsLogin] = useState(true);
  const toggleForm = () => setIsLogin(!isLogin);

  return (
    <div>
      {isLogin ? <LoginForm /> : <RegisterForm />}
      <p>
          {isLogin ? "Don't already have an account? " : "Already have an account? "} 
          <a href="#" onClick={toggleForm}>{isLogin ? "Sign up" : "Login"}</a>
      </p>
      <p style={{ color: 'red' }}>
        This is a public project. 
        Please assume all data is publicly accessible, including passwords.
      </p>
    </div>
  );
}

function LoginForm() {
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

function RegisterForm() {
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
        alert('Registration successful! Please log in.');
        setPassword('');
        setConfirm('');
        e.target.reset();
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

export default function App() {

  return (
    <>
      <Header />
      <main>
        <AuthenticatedPage />
      </main>
    </>
  )
}
