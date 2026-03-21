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
            <header style={{ display: 'flex', alignItems: 'center' }}>
                <h1 style={{ flex: "1 0", margin: 0 }}>macro</h1>
                <nav style={{ display: 'flex', gap: '1rem' }}>
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

            <div className="warning-box" style={{ display: "flex" }}>
                <div style={{ marginRight: "0.5rem" }}>⚠️</div>
                <span>
                    This is a development version of the app. It's subject to extreme change and
                    frequent database resets. 
                    <br />
                    Assume all data is publically accessible, including 
                    passwords. Don't use real passwords here.
                </span>
            </div>
    
            <Outlet />
        </>
    )
}