import React from 'react';

export default function Modal({ children, isVisible, setIsVisible }) {
    if (!isVisible) return null;

    return (
        <div className="fixed inset-0 bg-gray-800 bg-opacity-75 flex items-center justify-center">
            <div className="bg-white rounded-lg p-8 w-96 shadow-md relative">
                <button
                    className="absolute top-2 right-3 text-gray-500 hover:text-gray-800"
                    onClick={() => {
                        setIsVisible(false);
                    }}
                >
                    &times;
                </button>
                {children}
            </div>
        </div>
    )
}