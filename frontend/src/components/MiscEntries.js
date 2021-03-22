import React from "react"

export default function MiscEntries({ misc }) {

    if (!misc) return null

    return (
        <div className="mx-2 flex flex-row text-gray-700">
            {
                Object.entries(misc).map(([k, v]) => (
                    <div key={k} className="mt-6 w-1/6">
                        <h3 className="font-bold capitalize">
                            {k.replace("_", " ")}
                        </h3>
                        <p>{v}</p>
                    </div>
                ))
            }
        </div>
    )
}