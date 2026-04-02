import { useState, useEffect, useContext } from 'react';

import { NotificationContext } from '../Context';

/** Displays notifications from the NotificationContext. */
export default function NotificationArea() {
    const { notifications, remove } = useContext(NotificationContext);

    const style = {
        position: 'fixed', bottom: 0, right: 0, zIndex: 1000,
        display: 'flex', flexDirection: 'column', gap: '0.5rem',
        maxWidth: '25rem', maxHeight: '100vh', padding: '0.5rem',
        overflow: 'auto', backgroundColor: 'rgba(5,5,25,0.75)', borderRadius: '0.5rem'
    };

    return (
        <div style={style}>
            { notifications.length > 1 && <>
            <button onClick={ () => notifications.forEach(n => remove(n.id)) }>Dismiss All</button>
            </> }
            { notifications.toReversed().map(n => (
                <article key={ n.id } style={{ margin: 0 }}>
                    <header style={{ display: 'flex', justifyContent: 'space-between' }}>
                        <div>
                            { n.heading }
                            <div style={{ fontSize: '0.75rem', opacity: 0.5 }}>{ n.id }</div>
                        </div>
                        <button style={{ 
                            aspectRatio: "1",
                            lineHeight: 0,
                            margin: "0.2rem",
                            padding: "0.75rem",
                        }} onClick={ () => remove(n.id) }>&times;</button>
                    </header>
                    <small>
                        Details:
                        { n.details && <pre style={{ whiteSpace: 'pre-wrap' }}>{ n.details }</pre> }    
                    </small>
                </article>
            )) }
        </div>
    );
}