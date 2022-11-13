import React, { useState } from "react"
import moment from "moment"
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faCheck, faTriangleExclamation, faAngleDown } from '@fortawesome/free-solid-svg-icons'
import Card from "./Card"

export default function Message({ title, status, content, updated }) {
    const [ expanded, setExpanded ] = useState(false)
    
    const toggleExpanded = event => {
        event.preventDefault()
        setExpanded(!expanded)
    }

    return (
        <Card>
            <details open={expanded}>
                <summary 
                    onClick={toggleExpanded} 
                    className="cursor-pointer flex justify-between items-baseline"
                >
                    <span>
                        <span>
                            <FontAwesomeIcon
                                icon={faAngleDown}
                                className="text-gray-300"
                                rotation={ expanded ? 80 : 270}
                            />
                        </span>
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
                        <span className="select-none text-lg">
                            {title}
                        </span>
                    </span>
                    { !expanded &&
                        <span className="text-gray-400 hidden md:block flex-none">
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
        </Card>
    )
}