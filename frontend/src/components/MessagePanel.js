import React, { useState } from "react"
import { useAtomFeed } from "@au5ton/use-atom-feed"
import MessagePanelMessage from "./MessagePanelMessage"
import Card from "./Card"

export default function MessagePanel() {
    const { data } = useAtomFeed("/messages.atom")
    const [ showResolved, setShowResolved ] = useState(false)

    const toggleShowResolved = () => {
        setShowResolved(!showResolved)
    }

    const entries = data?.entries.filter(item => {
        if (showResolved) {
            return true
        } else {
            if (item.summary?.value === "Status: UNRESOLVED") {
                return true
            }
            return false
        }
    })

    const reactEntries = entries?.map(item => {
        return (
            <div className="mx-auto max-w-5xl my-4">
                <MessagePanelMessage 
                    key={item.id} 
                    title={item.title.value} 
                    status={item.summary?.value}
                    content={item.content?.value}
                    updated={item.updated}
                />
            </div>
        )
    })

    return (
        <>
            <div className="mx-auto max-w-5xl my-12 py-4">
                <div className="flex gap-4 items-baseline">
                    <p className="text-xl">
                        Messages
                    </p>
                    <button 
                        className="text-sm text-gray-400 focus:outline-none hover:text-green-300"
                        onClick={toggleShowResolved}
                    >
                        { showResolved ? "hide resolved" : "show resolved" }
                    </button>
                </div>
                <div>
                    { entries && entries.length > 0
                        ? reactEntries
                        : <div className="mx-auto max-w-5xl my-4">
                            <Card>
                                No unresolved{ showResolved && " or resolved" } messages posted
                            </Card>
                        </div>
                    }
                </div>
            </div>
        </>
    )
}