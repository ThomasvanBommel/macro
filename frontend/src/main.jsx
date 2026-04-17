import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { Outlet, BrowserRouter, Routes, Route, Navigate, useNavigate } from 'react-router';
import { useState, useEffect, useContext } from 'react';

import { SessionContext, useSession } from './Context';
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
    const session = useSession();
    const navigate = useNavigate();

    const [loading, setLoading] = useState(true);

    useEffect(() => { document.title = "macro"; }, []);

    useEffect(() => {
        if (!session.isValid())
            return setLoading(false);

        const i = setInterval(() => {
            console.log("Checking session validity...");
            session.isValid();
        }, 1000);

        return () => clearInterval(i);
    }, [session]);

    const handleAuthError = err => {
        session.notifications.add({
            heading: "Failed to register",
            details: err.message,
            type: "error",
            ttl: 0
        });
    };

    const handleRegisterSuccess = res => {
        navigate("/login", { replace: true });
    };

    const handleLoginSuccess = res => {
        session.refresh();
    };

    if (loading || session.changingState)
        return <article aria-busy="true" aria-label="Loading..." />;

    return (
        <SessionContext value={ session }>
            <Header />
            <hr />
            <NotificationArea />

            <Routes>
                <Route index element={ <Home /> } />

                <Route element={<UnAuthorizedOnly />}>
                    <Route path="login" element={
                        <LoginForm
                            onError={ handleAuthError }
                            onSuccess={ handleLoginSuccess } />
                    } />
                    <Route path="register" element={
                        <RegisterForm
                            onError={ handleAuthError }
                            onSuccess={ handleRegisterSuccess } />
                    } />
                </Route>

                <Route element={<AuthorizedOnly />}>
                    <Route path="profile" element={<Profile />} />
                </Route>
            </Routes>

        </SessionContext>
    )
}

function Root() {
    return (
        <BrowserRouter>
            <App />
        </BrowserRouter>
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
