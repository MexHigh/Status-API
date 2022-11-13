import React from "react"

export default function Card({ cRef, children }) {
    return (
        <div 
            className="bg-white px-6 md:px-10 py-6 shadow-md rounded-lg" 
            ref={cRef}
        >
            { children }
        </div>
    )
}