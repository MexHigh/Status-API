import React, { useState } from "react"
import moment from "moment"
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faCheck, faTriangleExclamation } from '@fortawesome/free-solid-svg-icons'

export default function Message({ title, status, content, updated }) {
    const [ expanded, setExpanded ] = useState(false)
    
    const toggleExpanded = event => {
        event.preventDefault()
        setExpanded(!expanded)
    }

    return (
        <div className="bg-white my-4 py-6 px-12 rounded-lg shadow-md">
            <details open={expanded}>
                <summary 
                    onClick={toggleExpanded} 
                    className="cursor-pointer"
                >
                    <span className="mx-2">
                        { status === "Status: RESOLVED"
                            ? <FontAwesomeIcon 
                                icon={faCheck} 
                                className="text-green-400"
                                fixedWidth
                            />
                            : <FontAwesomeIcon 
                                icon={faTriangleExclamation} 
                                className="text-red-400"
                                fixedWidth
                            />
                        }
                    </span>
                    <span className="select-none text-lg">{title}</span>
                    { !expanded &&
                        <span className="text-gray-400 float-right">
                            {moment(updated).fromNow()}
                        </span>
                    }
                </summary>
                <div className="p-4 flex flex-col gap-2">
                    <p className="text-gray-400">
                        <span>Last update: </span> 
                        <span>{moment(updated).calendar()}</span>
                    </p>
                    <p>
                        { status === "Status: RESOLVED" && <span className="font-bold">[Resolved] </span> }
                        { content || "No content" }
                    </p>
                </div>
            </details>
        </div>
    )
}