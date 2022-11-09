import React from "react"
import { useAtomFeed } from "@au5ton/use-atom-feed";

export default function RSSPanel() {
    const { data, error, isValidating } = useAtomFeed("/messages.atom")

    return (
        <>
            <div className="mx-auto max-w-4xl border-4 border-red-300 my-12">
                <ul>
                    {data?.entries.map(item => {
                        return <li key={item.id}>{item.title.value}</li>
                    })}
                </ul>
            </div>
            <pre>
                {JSON.stringify(data?.entries[0], undefined, 2)}
            </pre>
        </>
    )
}