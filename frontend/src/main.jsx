import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { Outlet, BrowserRouter, Routes, Route, Navigate } from 'react-router';
import { useState, useEffect, useContext } from 'react';

import { SessionContext, useSession, NotificationContext, useNotifications } from './Context';
import { Header, LoginForm, RegisterForm, NotificationArea } from './Components';
import { Home, Profile } from './Pages';

import '@picocss/pico';
import './main.css';

/** Renders nested routes only when a valid session exists. */
function AuthorizedOnly() {
    const session = useContext(SessionContext);

    if (session.isValid()) 
        return <Outlet />;

    return <Navigate to="/login" replace />;
}

/** Renders nested routes only when no session exists. */
function UnAuthorizedOnly() {
    const session = useContext(SessionContext);

    if (!session.isValid())
        return <Outlet />;

    return <Navigate to="/profile" replace />;
}

/** Manages session state and renders auth-gated routes. */
function App() {
    const [loading, setLoading] = useState(true);
    const session = useSession();

    useEffect(() => {
        if (!session.isValid())
            return setLoading(false);

        const i = setInterval(() => {
            console.log("Checking session validity...");
            session.isValid();
        }, 1000);

        return () => clearInterval(i);
    }, [session]);

    if (loading) return <article aria-busy="true" aria-label="Loading..." />;

    return (
        <BrowserRouter>
        <SessionContext value={ session }>
            <Header />

            <Routes>
                <Route index element={ <Home /> } />

                <Route element={<UnAuthorizedOnly />}>
                    <Route path="login" element={<LoginForm />} />
                    <Route path="register" element={<RegisterForm />} />
                </Route>

                <Route element={<AuthorizedOnly />}>
                    <Route path="profile" element={<Profile />} />
                </Route>
            </Routes>

            <NotificationArea />
        </SessionContext>
        </BrowserRouter>
    )
}

function Root() {
    const notifications = useNotifications();

    useEffect(() => {
        document.title = "macro";
    }, []);

    return (
        <NotificationContext value={ notifications }>
            <App />
        </NotificationContext>
    );
}

/** Application entry point. */
createRoot(document.getElementById('root')).render(
    <StrictMode>
        <main className="container">
            <Root />
        </main>
    </StrictMode>,
)
