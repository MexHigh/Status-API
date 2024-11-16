import { useState } from "react";
import Card from "./Card";
import Button from "./Button";
import TextInput from "./TextInput";

export default function LoginForm({ setApiKey }) {
    const [input, setInput] = useState("")
    const [errorMessage, setErrorMessage] = useState(null)

    const submit = () => {
        fetch("/api/auth/api-key/test", {
            method: "POST",
            headers: {
                "X-API-Key": input
            }
        })
            .then(r => r.json())
            .then(r => {
                if (r.error !== null) {
                    setInput("")
                    setErrorMessage(r.error)
                    setTimeout(() => {
                        setErrorMessage(null)
                    }, 2000)
                    return
                }
                if (r.response === "ok") {
                    setApiKey(input)
                }
            })
    }

    return (
        <Card>
            <form
                className="flex flex-col space-y-4"
                onSubmit={e => {
                    e.preventDefault()
                    if (!input) return false
                    submit()
                }}
            >
                <TextInput 
                    label="Message API key" 
                    placeholder="Enter your Message API key"
                    value={input}
                    onChange={e => setInput(e.target.value)}
                />
                <Button type="submit" text="Login" primary wFull />
                {errorMessage && (
                    <p>Error: {errorMessage}</p>
                )}
            </form>
        </Card>
    )
}
