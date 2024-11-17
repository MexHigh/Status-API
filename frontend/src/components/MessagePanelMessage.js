import React, { useState } from "react"
import moment from "moment"
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faCheck, faTriangleExclamation, faAngleRight } from '@fortawesome/free-solid-svg-icons'
import Card from "./Card"
import Button from "./Button"
import EditMessageModal from "./EditMessageModal"

// This component is compatible for both Atom Feed entries and API entries!
export default function MessagePanelMessage({ id, title, status, content, updated, withEditButtons = false, doPatch, doDelete, doRefetch }) {
    if (withEditButtons === true && (!id || !doPatch || !doDelete || !doRefetch)) {
        throw new Error("if 'withEditButtons' is defined, you also need to define 'doPatch(change)', 'doDelete()' and 'doRefetch()'")
    }

    const [expanded, setExpanded] = useState(false)
    const [modalVisible, setModalVisible] = useState(false)

    const toggleExpanded = event => {
        event.preventDefault()
        setExpanded(!expanded)
    }

    const toggleResolved = () => {
        doPatch({
            Resolved: status === "Status: RESOLVED" ? false : true
        })
            .then(() => {
                doRefetch()
            })
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
                                icon={faAngleRight}
                                className="text-gray-300"
                                rotation={expanded ? 90 : 0}
                            />
                        </span>
                        <span className="mx-2">
                            {status === "Status: RESOLVED"
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
                    {!expanded &&
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
                        {status === "Status: RESOLVED" && <span className="font-bold">[Resolved] </span>}
                        {content || "No content"}
                    </p>
                    {withEditButtons && (
                        <>
                            <EditMessageModal 
                                isVisible={modalVisible} 
                                setIsVisible={setModalVisible}
                                initTitle={title}
                                initContent={content}
                                doPatch={(change) => doPatch(change)}
                                doRefetch={() => doRefetch()}
                            />
                            <div className="mt-2 flex gap-4">
                                <Button 
                                    text={ status === "Status: RESOLVED" ? "Unresolve" : "Resolve" } 
                                    onClick={() => {
                                        toggleResolved()
                                    }}
                                    wFull={false}
                                />
                                <Button 
                                    text="Edit" 
                                    onClick={() => {
                                        setModalVisible(true)
                                    }}
                                    wFull={false}
                                />
                                <Button
                                    text="Delete"
                                    onClick={() => {
                                        doDelete().then(() => {
                                            doRefetch()
                                        })
                                    }}
                                    wFull={false}
                                />
                            </div>
                        </>
                    )}
                </div>
            </details>
        </Card>
    )
}