// A simple modal component that can be used to display content in a dialog.
export default function Modal({ isOpen, onClose, children }) {
    if (!isOpen) return null;

    window.onkeydown = e => {
        e.stopImmediatePropagation();
        if (e.key === "Escape") onClose?.();
    };

    return (
        <dialog open={isOpen}>
            <article style={{ position: "relative" }}>
                { children }
            </article>
        </dialog>
    );
}

// A simple header component for modals, with a title and a close button.
export function ModalHeader({ title, onClose }) {
    return (
        <header style={{ display: "flex",
                         justifyContent: "space-between",
                         alignItems: "center" }}>
            <h3 style={{ margin: 0 }}>{ title }</h3>
            <button style={{ 
                aspectRatio: "1",
                lineHeight: 0,
                margin: "0.2rem",
                padding: "0.75rem",
            }} onClick={ () => onClose?.() }>&times;</button>
        </header>
    );
}