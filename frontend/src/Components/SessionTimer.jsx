import { useState, useEffect } from "react";

/**
 * Callback for when the session expires. Should trigger a re-check of the session, and redirect to login if expired.
 * @callback onExpire
 * @returns { void }
 */

/**
 * Component to display a countdown timer until session expiration. Calls onExpire callback when 
 * expired.
 * @param { Object } props - The component props
 * @param { string|Date } props.expiration - The expiration time, as a Date object or a string 
 *                                           parseable by Date
 * @param { onExpire } [props.onExpire] - Optional callback to call when the timer expires
 * @returns { JSX.Element } The rendered component
 */
export default function SessionTimer({ expiration, onExpire }) {
    const [timeLeft, setTimeLeft] = useState("");

    useEffect(() => {
        const interval = setInterval(() => {
            const diff = new Date(expiration) - new Date();
            const min = Math.floor(diff/60000);
            const sec = Math.floor((diff%60000)/1000);
            
            setTimeLeft(`${min}'${sec}`);
            if (diff <= 0) onExpire?.();
        }, 1000);

        return () => clearInterval(interval);
    }, [expiration]);

    return <>{ timeLeft }</>;
}