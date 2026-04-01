import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { Outlet, BrowserRouter, Routes, Route, Navigate } from 'react-router';
import { useState, useEffect } from 'react';

import { handleFetchSession } from './api';
import { Header } from './Components';
import { Home, Profile } from './Pages';
import { LoginForm, RegisterForm } from './Components';
import '@picocss/pico';
import './main.css';

/**
 * Component that renders its children only if the user is authenticated (i.e. has a valid session).
 * If the user is not authenticated, they are redirected to the login page.
 * @param { Object } props - The component props.
 * @param { Object } props.session - The user session object containing session information.
 * @returns { JSX.Element } - The rendered component based on the user's authentication status.
 */
function AuthorizedOnly({ session }) {
    return session ? <Outlet /> : <Navigate to="/login" replace />;
}

/**
 * Component that renders its children only if the user is not authenticated (i.e. does not have a 
 * valid session). If the user is authenticated, they are redirected to the profile page.
 * @param { Object } props - The component props.
 * @param { Object } props.session - The user session object containing session information.
 * @returns { JSX.Element } - The rendered component based on the user's authentication status.
 */
function UnAuthorizedOnly({ session }) {
    return !session ? <Outlet /> : <Navigate to="/profile" replace />;
}

/**
 * The main application component that manages the user session and renders the appropriate routes 
 * based on the user's authentication status. It uses the useState hook to manage the session and 
 * loading state, and the useEffect hook to fetch the session data when the component mounts. It 
 * also defines a nullSession function to reset the session state when the user logs out or when
 * the session becomes invalid. The component renders a loading indicator while the session data is 
 * being fetched, and once the data is available, it renders the appropriate routes based on the 
 * user's authentication status. The Header component is always rendered, and the Home, LoginForm,
 * RegisterForm, and Profile components are rendered based on the user's authentication status and 
 * the current session state.
 * @returns { JSX.Element } - The rendered application component with routing and session management
 */
function App() {
    const [updateSession, setUpdateSession] = useState(0);
    const [session, setSession] = useState(null);
    const [loading, setLoading] = useState(true);

    const nullSession = () => {
        setSession(null);
        setLoading(true);
        setUpdateSession(s => s + 1);
    };

    useEffect(() => {
        handleFetchSession()
            .then(setSession)
            .catch(() => setSession(null))
            .finally(() => setLoading(false));
    }, [updateSession])

    if (loading) return <article aria-busy="true" aria-label="Loading..." />;

    return (
        <BrowserRouter>
            <Header session={ session } />

            <Routes>
                <Route index element={ <Home /> } />

                <Route element={<UnAuthorizedOnly session={ session } />}>
                    <Route path="login" element={<LoginForm onSuccess={ nullSession } />} />
                    <Route path="register" element={<RegisterForm onSuccess={ nullSession } />} />
                </Route>

                <Route element={<AuthorizedOnly session={ session } />}>
                    <Route path="profile" element={<Profile session={ session } />} />
                </Route>
            </Routes>
        </BrowserRouter>
    )
}

/**
 * The main entry point of the application. It renders the App component inside a StrictMode wrapper
 * and a main container element. The StrictMode wrapper helps identify potential problems in the 
 * application by activating additional checks and warnings for its descendants.
 */
createRoot(document.getElementById('root')).render(
    <StrictMode>
        <main className="container">
            <App />
        </main>
    </StrictMode>,
)
