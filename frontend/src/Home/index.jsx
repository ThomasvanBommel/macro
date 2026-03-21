import { Outlet, NavLink, Link, useNavigate } from 'react-router';
import { clearSession } from '../api';
import './index.css';

export default function Home({ session, authCheck }) {
    const navigate = useNavigate();

    function Logout(){
        clearSession()
            .then(() => {
                authCheck().then(() => navigate("/login"));
            })
            .catch(alert);
    }

    return (
        <>
            <header id="main-header">
                <h1 className="logo">macro</h1>
                <nav>
                    <NavLink to="/">Home</NavLink>
                    { session ? 
                        <>
                        <NavLink to="/profile">Profile</NavLink> 
                        <Link to="/logout" onClick={Logout}>Logout</Link>
                        </>
                        : 
                        <>
                            <NavLink to="/login">Login</NavLink>
                            <NavLink to="/register">Register</NavLink>
                        </>}
                </nav>
            </header>

            <br />

            <p style={{ textAlign: 'center', color: 'red' }}>
                This is a development version of the app. It's subject to extreme change and
                frequent database resets.
            </p>
    
            <Outlet />
        </>
    )
}