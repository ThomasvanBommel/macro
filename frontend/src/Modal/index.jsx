import { useEffect } from 'react';
import './index.css';

export default function Modal({ isOpen, onClose, children }) {
    useEffect(() => {
        function onKeyDown(e) {
            if (e.key === "Escape") {
                onClose();
            }
        }

        if (isOpen) {
            window.addEventListener("keydown", onKeyDown);
        } else {
            window.removeEventListener("keydown", onKeyDown);
        }

        return () => window.removeEventListener("keydown", onKeyDown);
    }, [isOpen, onClose]);

    if (!isOpen) return null;

    return (
        <div className="modal-container" onClick={onClose}>
            <div className="modal-content" onClick={e => e.stopPropagation()}>
                {children}
            </div>
        </div>
    );
}

export function ModalForm({ isOpen, onClose, children }) {
    return (
        <Modal isOpen={isOpen} onClose={onClose}>
            <div style={{ padding: '1rem' }}>
                {children}
            </div>
        </Modal>
    );
}