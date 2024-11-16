import React, { useState } from "react"
import { useAtomFeed } from "@au5ton/use-atom-feed"
import MessagePanelMessage from "./MessagePanelMessage"
import Card from "./Card"

export default function AdminMessagePanel() {
    const { data } = useAtomFeed("/messages.atom")
    const entries = data?.entries

    const reactEntries = entries?.map(item => (
        <div className="my-4" key={item.id}>
            <MessagePanelMessage
                title={item.title.value}
                status={item.summary?.value}
                content={item.content?.value}
                updated={item.updated}
                withEditButtons
            />
        </div>
    ))

    return (
        <div>
            <h1 className="text-xl">Edit Messages</h1>
            {entries && entries.length > 0
                ? reactEntries
                : <div className="my-4">
                    <Card>
                        No messages posted
                    </Card>
                </div>
            }
        </div>
    )
}
