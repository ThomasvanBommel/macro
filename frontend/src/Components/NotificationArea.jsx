import { useState, useEffect, useContext, useRef } from 'react';
import { PiWarningFill } from "react-icons/pi";

import { SessionContext } from '../Context';

/** Displays notifications from the NotificationContext. */
export default function NotificationArea() {
    const { notifications: { list, remove } } = useContext(SessionContext);

    return (
        <div>
        { list.map(n => {
            return (
                <article 
                    key={ n.id }
                    className="warning-box notification"
                    style={{ position: "relative" }}
                    onClick={ () => remove(n.id) }>
                    <h5>
                        <span style={{ marginRight: "0.5rem" }}>
                            <PiWarningFill color="yellow" />
                        </span>
                        { n.heading }
                    </h5>
                    <div style={{ }}>
                        { n.details }
                    </div>
                    <span style={{ 
                        position: "absolute", 
                        bottom: "0.2rem", 
                        right: "0.5rem", 
                        opacity: 0.5,
                        fontSize: "0.75rem",
                     }}>Click to dismiss</span>
                     <span style={{
                        position: "absolute",
                        top: "0.2rem",
                        right: "0.5rem",
                        opacity: 0.5,
                        fontSize: "0.75rem"
                     }}>
                        <Countdown expires={ n.expires } />
                     </span>
                </article>
            );
        }) }
        </div>
    );
}

function Countdown({ expires }) {
    const [timeLeft, setTimeLeft] = useState(expires - new Date());
    const aref = useRef();
    const bref = useRef();

    const animate = t => {
        if (!aref.current) aref.current = t;

        const remaining = expires - new Date();
        setTimeLeft(remaining);

        if (remaining > 0)
            bref.current = requestAnimationFrame(animate);
    }

    useEffect(() => {
        bref.current = requestAnimationFrame(animate);
        return () => cancelAnimationFrame(bref.current);
    }, []);

    if (timeLeft <= 0) return null;

    return (
        <span style={{ fontVariantNumeric: "tabular-nums" }}>
            { (timeLeft / 1000).toFixed(1) }s
        </span>
    );

}