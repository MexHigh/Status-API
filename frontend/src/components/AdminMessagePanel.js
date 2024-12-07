import React, { useEffect, useState } from "react"
import MessagePanelMessage from "./MessagePanelMessage"
import Card from "./Card"
import EditMessageModal from "./EditMessageModal"

export default function AdminMessagePanel({ apiKey }) {
    const [ messages, setMessages ] = useState([])
    const [ createModalOpen, setCreateModalOpen ] = useState(false)

    const fetchMessages = () => {
        fetch("/api/messages", {
            headers: {
                "X-Api-Key": apiKey
            }
        })
            .then(r => r.json())
            .then(r => {
                setMessages(r.response)
            })
    }

    const createMessage = (payload) => {
        return fetch(`/api/message`, {
            method: "POST",
            headers: {
                "X-Api-Key": apiKey
            },
            body: JSON.stringify(payload)
        })
    }

    const patchMessage = (id, change) => {
        return fetch(`/api/message/${id}`, {
            method: "PATCH",
            headers: {
                "X-Api-Key": apiKey
            },
            body: JSON.stringify(change)
        })
    }

    const deleteMessage = id => {
        return fetch(`/api/message/${id}`, {
            method: "DELETE",
            headers: {
                "X-Api-Key": apiKey
            }
        })
    }

    useEffect(() => {
        fetchMessages()
    // eslint-disable-next-line
    }, [apiKey])

    const reactEntries = messages.toReversed().map(item => (
        <div className="my-4" key={item.id || item.Db_Id}>
            <MessagePanelMessage
                id={item.Db_Id}
                title={item.Title}
                status={item.Description}
                content={item.Content}
                updated={item.Updated}
                withEditButtons
                doPatch={(change) => patchMessage(item.Db_Id, change)}
                doDelete={() => deleteMessage(item.Db_Id)}
                doRefetch={fetchMessages}
            />
        </div>
    ))

    return (
        <div>
            <EditMessageModal
                isVisible={createModalOpen}
                setIsVisible={setCreateModalOpen}
                forCreation
                doCreate={payload => createMessage(payload)}
                doRefetch={fetchMessages}
            />
            <div className="flex gap-4 items-baseline ml-2">
                <h1 className="text-xl">Edit Messages</h1>
                <button 
                    className="text-sm text-gray-400 focus:outline-none select-none hover:text-green-300"
                    onClick={() => setCreateModalOpen(true)}
                >
                    Create new message
                </button>
            </div>

            {messages && messages.length > 0
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
