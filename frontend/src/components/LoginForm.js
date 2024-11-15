import { useState } from "react";
import Card from "./Card";
import WidthContainer from "./WidthContainer";

export default function LoginForm({ setApiKey }) {
    const [ input, setInput ] = useState("")
    const [ errorMessage, setErrorMessage ] = useState(null)

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
        <WidthContainer className="my-8">
            <Card>
                <form 
                    className="flex flex-col space-y-4"
                    onSubmit={e => {
                        e.preventDefault()
                        if (!input) return false
                        submit()
                    }}
                >
                    <label for="password" className="text-xl">Message API key</label>
                    <input 
                        type="password"
                        className="w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-0 focus:border-green-500"
                        placeholder="Enter your Message API key"
                        value={input}
                        onChange={e => setInput(e.target.value)}
                    />
                    <button 
                        type="submit"
                        className="w-full px-4 py-2 bg-green-500 text-white font-semibold rounded-md shadow-sm hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-1"
                    >
                        Login
                    </button>
                    { errorMessage && (
                        <p>Error: {errorMessage}</p>
                    )}
                </form>
            </Card>
        </WidthContainer>
    )
}
