import React from "react"

export default function Header({lastCheckTs}) {
    return (
        <header className="bg-black m-4 rounded-full">
            <div className="flex justify-between py-2 px-10">
                <p className="text-red-500">Status API</p>
                <p className="text-red-500">{new Date(lastCheckTs).toLocaleString("de-DE")}</p>
            </div>
        </header>
    )
}