import React from "react"

export default function MiscEntries({ misc }) {

    if (!misc) return null

    return (
        <div className="mx-2 text-gray-700 flex gap-4 md:gap-10">
            {
                // runs over the misc object an creates an own
                // div for every entry
                Object.entries(misc).map(([key, value]) => (
                    <div key={key} className="mt-6">
                        <h3 className="font-bold capitalize">
                            {key.replace("_", " ")}
                        </h3>
                        <p>{value}</p>
                    </div>
                ))
            }
        </div>
    )
}