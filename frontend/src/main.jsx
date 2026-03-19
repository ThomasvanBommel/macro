import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { Outlet, BrowserRouter, Routes, Route, Navigate } from 'react-router';
import { useState, useEffect, useRef } from 'react';

import { fetchSession } from './api';
import Home from './Home';
import { Login, Register } from './Authentication';
import Profile from './Profile';

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

    const authCheck = async () => {
        setSession(await fetchSession());
        setLoading(false);
    };

    useEffect(() => {
        if (!init.current) {
            init.current = true;
            authCheck();
        }
    }, []);

    if (loading) return <div>Loading...</div>;

    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Home session={ session } authCheck={ authCheck } />}>

                    {/* <Route element={<AuthLayout />}> */}
                    <Route element={<UnAuthorizedOnly session={ session } />}>
                        <Route path="login" element={<Login onSuccess={ authCheck } />} />
                        <Route path="register" element={<Register />} />
                    </Route>

                    <Route element={<AuthorizedOnly session={ session } />}>
                        <Route path="profile" element={<Profile session={ session } />} />
                    </Route>

                </Route>
            </Routes>
        </BrowserRouter>
    )
}

createRoot(document.getElementById('root')).render(
    <StrictMode>
        <App />
    </StrictMode>,
)
