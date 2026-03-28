import { NavLink, Link, useNavigate } from 'react-router';
import { clearSession } from '../api';
import { SessionTimer } from '../Components';

// Header component, shown on all pages. Displays different links based on session state.
export default function Header({ session, authCheck }) {
    const navigate = useNavigate();

    function Logout(){
        clearSession()
            .then(() => {
                authCheck().then(() => navigate("/login"));
            })
            .catch(alert);
    }

    const timerStyle = {
        width: "100%",
        position: "absolute",
        textAlign: "center",
        top: "1.2rem",
        left: "0",
        opacity: 0.5,
        whiteSpace: "nowrap"
    };

    return (
        <>
        <header style={{ display: 'flex', alignItems: 'center' }}>
            <h1 style={{ flex: "1 0", margin: 0 }}>macro</h1>
            <nav style={{ display: 'flex', gap: '1rem' }}>
                <NavLink to="/">Home</NavLink>
                { session ? 
                    <>
                    <NavLink to="/profile">Profile</NavLink>
                    <Link to="/logout" onClick={Logout} style={{ position: "relative" }} >
                        Logout
                        <small style={ timerStyle }>
                            <SessionTimer expiration={ session.expires } onExpire={ authCheck } />
                        </small>
                    </Link>
                    </>
                    : 
                    <>
                        <NavLink to="/login">Login</NavLink>
                        <NavLink to="/register">Register</NavLink>
                    </>}
            </nav>
        </header>

        <hr />

        <div className="warning-box" style={{ display: "flex" }}>
            <div style={{ marginRight: "0.5rem" }}>⚠️</div>
            <span>
                This is a development version of the app. It's subject to extreme change. 
            </span>
        </div>
        </>
    )
}