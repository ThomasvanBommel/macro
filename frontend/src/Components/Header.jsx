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

    return <>
        <header id="main-header">
            <h1><NavLink to="/">macro</NavLink></h1>
            <nav>
                {
                    isValid() ? <>
                    <NavLink to="/profile">Profile</NavLink>
                    <Link onClick={ end }>
                        Logout
                        <small id="session-timer"><SessionTimer /></small>
                    </Link>
                    </> : <>
                    <NavLink to="/login">Login</NavLink>
                    <NavLink to="/register">Register</NavLink>
                    </>
                }


                {/* { isValid() ? 
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
                    </>} */}
            </nav>
        </header>

<style>{`
#main-header {
    display: flex;
    align-items: center;
    margin-bottom: var(--pico-typography-spacing-vertical);
    border-bottom: 1px solid var(--pico-muted-border-color);

    & > h1 {
        flex: 1 0;
        margin: 0;

        & > a {
            text-decoration: none;
            color: inherit;
        }
    }

    & > nav {
        display: flex;
        gap: 1px;

        & > a {
            position: relative;
            padding: 0 0.5rem 0 0.5rem;
            margin-top: 0.75rem;
            border-radius: 0.5rem 0.5rem 0 0;
            outline: 1px solid var(--pico-muted-border-color);
            border-bottom: none;
        }

        & > a.active {
            background-color: var(--pico-primary-background);
            color: var(--pico-primary-inverse);
        }

        #session-timer {
            width: 100%;
            text-align: center;
            position: absolute;
            left: 0;
            top: 1.4rem;
            opacity: 0.5;
        }
    }
}
`}</style>
    </>;
}