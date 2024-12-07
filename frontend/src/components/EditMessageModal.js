import React, { useState } from "react"
import Button from "./Button"
import Modal from "./Modal"
import TextInput from "./TextInput"
import Textarea from "./Textarea"

export default function EditMessageModal({ isVisible, setIsVisible, initTitle = "", initContent = "", forCreation = false, doCreate, doPatch, doRefetch }) {
    if (forCreation === true && !doCreate) {
        throw new Error("if 'forCreation' is true, define 'doCreate(obj)'")
    }
    if (!forCreation && !doPatch) {
        throw new Error("if 'forCreation' is false, define 'doPatch(change)'")
    }
    if (!doRefetch) {
        throw new Error("define 'doRefetch'")
    }

    const [ title, setTitle ] = useState(initTitle)
    const [ content, setContent ] = useState(initContent)
    
    const handleSubmit = e => {
        e.preventDefault()

        if (forCreation === true) {
            let payload = {
                Title: title,
                Content: content
            }
            doCreate(payload).then(() => {
                doRefetch()
                setIsVisible(false)
            })
        } else {
            let change = {}
            if (initTitle !== title) {
                change.Title = title
            }
            if (initContent !== content) {
                change.Content = content
            }
            if (Object.keys(change).length > 0) {
                doPatch(change).then(() => {
                    doRefetch()
                    setIsVisible(false)
                })
            }
        }
    }

    return (
        <div>
            <Modal
                isVisible={isVisible}
                setIsVisible={setIsVisible}
            >
                <form 
                    className="flex flex-col space-y-4"
                    onSubmit={handleSubmit}
                >
                    <TextInput
                        label="Title"
                        value={title}
                        minlength={5}
                        onChange={e => setTitle(e.target.value)}
                    />
                    <Textarea
                        label="Content"
                        value={content}
                        minlength={5} 
                        rows={10}
                        onChange={e => setContent(e.target.value)}
                    />
                    <Button
                        primary
                        type="submit"
                        text={ forCreation === true ? "Create" : "Update" }
                    />
                </form>
            </Modal> 
        </div>
    )
}
