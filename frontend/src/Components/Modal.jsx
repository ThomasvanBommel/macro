/**
 * A simple modal component that uses the HTML <dialog> element. It also listens for the Escape key 
 * to close the modal.
 * @param { Object } props - The component props
 * @param { boolean } props.isOpen - Whether the modal is open or not
 * @param { function } props.onClose - Callback to call when the modal should be closed (e.g. when 
 *                                     Escape is pressed)
 * @param { React.ReactNode } props.children - The content of the modal
 * @returns { JSX.Element|null } The rendered component or null if not open
 */
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

/**
 * A simple header component for modals, with a title and a close button.
 * @param { Object } props - The component props
 * @param { string } props.title - The title to display in the header
 * @param { function } [props.onClose] - Optional callback to call when the close button is clicked
 * @returns { JSX.Element } The rendered component
 */
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