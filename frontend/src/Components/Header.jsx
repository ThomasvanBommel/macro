import { NavLink, Link } from 'react-router';

import { useSession } from '../Context';
import { SessionTimer } from '../Components';

// Header component, shown on all pages. Displays different links based on session state.
export default function Header() {
    const { isValid, end } = useSession();

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
        <header style={{ display: 'flex', alignItems: 'center' }}>
            <h1 style={{ flex: "1 0", margin: 0 }}>
                <NavLink to="/" style={{ textDecoration: "none", color: "inherit" }}>
                    macro
                </NavLink>
            </h1>
            <nav style={{ display: 'flex', gap: '1rem' }}>
                { isValid() ? 
                    <>
                    <NavLink to="/profile">Profile</NavLink>
                    <Link to="/login" onClick={ end } style={{ position: "relative" }} >
                        Logout
                        <small style={ timerStyle }>
                            <SessionTimer />
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
    );
}