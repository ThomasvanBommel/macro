import { Outlet, NavLink, Link, useNavigate } from 'react-router';
import { clearSession } from '../api';
import './index.css';

export default function Home({ session, authCheck }) {
    const navigate = useNavigate();

    function Logout(){
        clearSession();
        authCheck().then(() => navigate("/login"));
    }

    return (
        <>
            <header id="main-header">
                <h1 className="logo">macro</h1>
                <nav>
                    <NavLink to="/">Home</NavLink>
                    { session ? 
                        <Link onClick={ Logout }>Logout</Link>: 
                        <>
                            <NavLink to="/login">Login</NavLink>
                            <NavLink to="/register">Register</NavLink>
                        </>}
                </nav>
            </header>

            <br />
    
            <Outlet />
        </>
    )
}