import { Outlet, BrowserRouter, Routes, Route, Navigate, useNavigate } from 'react-router';
import { createRoot } from 'react-dom/client';
import { StrictMode, useEffect, useCallback} from 'react';

import { SessionProvider, useSession } from './Context';
import { Header, LoginForm, RegisterForm, NotificationArea } from './Components';
import { Home, Profile } from './Pages';

import '@picocss/pico';
import './main.css';

const UnauthorizedOnly = () => {
    const { isValid } = useSession();
    return isValid() ? <Navigate to="/profile" replace /> : <Outlet />;
};

const AuthorizedOnly = () => {
    const { isValid } = useSession();
    return isValid() ? <Outlet /> : <Navigate to="/login" replace />;
};

// The main application component that sets up routing and session management.
function App() {
    const { init, isValid, notifications } = useSession();
    const navigate = useNavigate();

    const handleAuthError = err => {
        console.error(err);
        notifications.add({
            heading: "Failed to authenticate",
            details: err.message,
            type: "error",
            ttl: 0
        });
    };

    const handleRegisterSuccess = res => {
        navigate("/login", { replace: true });
    };

    const handleLoginSuccess = res => {
        init();
    };

    return <>
        <Header />
        <NotificationArea />

        <Routes>
            <Route index element={ <Home /> } />

            <Route element={<UnauthorizedOnly />}>
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
    </>;
}

// Sets the document title and renders the application within the session provider and router.
function Root() {
    useEffect(() => { document.title = "macro"; }, []);

    return (
        <BrowserRouter>
            <SessionProvider>
                <App />
            </SessionProvider>
        </BrowserRouter>
    );
}

// Renders the root component into the DOM.
createRoot(document.getElementById('root')).render(
    <StrictMode>
        <main className="container">
            <Root />
        </main>
    </StrictMode>,
)
