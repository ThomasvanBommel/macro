import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { Outlet, BrowserRouter, Routes, Route, Navigate } from 'react-router';
import { useState, useEffect, useRef } from 'react';

import { fetchSession,  } from './api';
import { Header } from './Components';
import { Home } from './Pages';
import { Login, Register } from './Authentication';
import Profile from './Profile';
import '@picocss/pico';
import './main.css';

function AuthorizedOnly({ session }) {
    return session ? <Outlet /> : <Navigate to="/login" replace />;
}

function UnAuthorizedOnly({ session }) {
    return !session ? <Outlet /> : <Navigate to="/profile" replace />;
}

function App() {
    const [session, setSession] = useState(null);
    const [loading, setLoading] = useState(true);
    const init = useRef(false);

    // Check for existing session. Also acts as a callback for successful login/registration
    const authCheck = async (session_model) => {
        setSession(null);
        if(session_model) {
            setSession(session_model);
            setLoading(false);
            return;
        }

        setSession(await fetchSession());
        setLoading(false);
    };

    // Run auth check on initial load
    useEffect(() => {
        if (!init.current) {
            init.current = true;
            authCheck();
        }
    }, []);

    if (loading) return <div>Loading...</div>;

    return (
        <BrowserRouter>
            <Header session={ session } authCheck={ authCheck } />

            <Routes>
                {/* <Route path="/"> */}

                <Route index element={ <Home /> } />

                <Route element={<UnAuthorizedOnly session={ session } />}>
                    <Route path="login" element={<Login onSuccess={ authCheck } />} />
                    <Route path="register" element={<Register onSuccess={ authCheck } />} />
                </Route>

                <Route element={<AuthorizedOnly session={ session } />}>
                    <Route path="profile" element={<Profile session={ session } />} />
                </Route>

                {/* </Route> */}
            </Routes>
        </BrowserRouter>
    )
}

createRoot(document.getElementById('root')).render(
    <StrictMode>
        <main className="container">
            <App />
        </main>
    </StrictMode>,
)
