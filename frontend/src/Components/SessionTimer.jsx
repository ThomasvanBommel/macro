import { useState, useEffect, useContext } from "react";

import { SessionContext } from '../Context';

/** Displays a countdown to expiration and triggers onExpire when time runs out. */
export default function SessionTimer() {
    const session = useContext(SessionContext);
    if (!session.isValid()) return null;

    const [timeLeft, setTimeLeft] = useState("");

    useEffect(() => {
        const interval = setInterval(() => {
            const diff = new Date(session.expires) - new Date();
            const min = Math.floor(diff/60000);
            const sec = Math.floor((diff%60000)/1000);
            
            setTimeLeft(`${min}'${sec}`);
        }, 1000);

        return () => clearInterval(interval);
    }, [session]);

    return <>{ timeLeft }</>;
}