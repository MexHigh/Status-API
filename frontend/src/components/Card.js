import React from "react"

export default function Card({ cRef, children }) {
    return (
        <div 
            className="bg-white px-8 py-6 md:px-12 md:py-8 shadow-md rounded-lg" 
            ref={cRef}
        >
            { children }
        </div>
    )
}