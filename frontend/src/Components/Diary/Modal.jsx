import { useEffect } from "react";
import { MdClose } from "react-icons/md";

export default function Modal({ 
    title, 
    isOpen, 
    onClose, 
    children 
}) {
    useEffect(() => {
        if (!isOpen) return;
        const handleKeyDown = e => {
            if (e.key === "Escape") {
                e.preventDefault();
                onClose?.();
            }
        };

        window.addEventListener("keydown", handleKeyDown);
        return () => window.removeEventListener("keydown", handleKeyDown);
    }, [isOpen, onClose]);

    if (!isOpen) return null;

    return <>
        <style>{`
            .modal .close-button:hover {
                cursor: pointer;
                background-color: hsl(from var(--pico-color) h s l / 5%);
                border-radius: 0.25rem;
            }

            .modal .close-button:active {
                background-color: hsl(from var(--pico-color) h s l / 10%);
            }

            .modal > article {
                display: flex;
                flex-direction: column;
                padding: 0;
            }

            .modal > article > header {
                display: flex;
                justify-content: space-between;
                padding: var(--pico-block-spacing-vertical) var(--pico-block-spacing-horizontal);
                margin: 0;
            }
        `}</style>
        <dialog open onClick={ e => e.target === e.currentTarget && onClose?.() } className="modal">
            <article onClick={ e => e.stopPropagation() }>
                <header>
                    <strong>{ title }</strong>
                    <div onClick={ onClose }
                         title="Close"
                         style={{ lineHeight: 0, padding: "0.25rem" }}
                         className="close-button">
                        <MdClose />
                    </div>
                </header>
                <div className="content">{ children }</div>
            </article>
        </dialog>
    </>;
}