import { Outlet, NavLink } from 'react-router';
import './index.css';

export default function Home() {
    return (
        <>
            <header id="main-header">
                <h1 className="logo">macro</h1>
                <nav>
                    <NavLink to="/">Home</NavLink>
                    <NavLink to="/login">Login</NavLink>
                    <NavLink to="/register">Register</NavLink>
                </nav>
            </header>

            <br />
    
            <Outlet />
        </>
    )
}