import { useState, useEffect } from "react";

import { useSession } from '../Context';

/** Displays a countdown to expiration and triggers onExpire when time runs out. */
export default function SessionTimer() {
    const { isValid, expires, ttl } = useSession();
    if (!isValid) return null;

    const [timeLeft, setTimeLeft] = useState("");

    useEffect(() => {
        const interval = setInterval(() => {
            const diff = new Date(expires) - new Date();
            const min = Math.floor(diff/60000);
            const sec = Math.floor((diff%60000)/1000);
            
            setTimeLeft(`${min}'${sec}`);
        }, 1000);

        return () => clearInterval(interval);
    }, [expires]);

    return ttl() < 0 ? null : <>{ timeLeft }</>;
}