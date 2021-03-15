import React, { useState, useEffect } from "react"

export default function CurrentStatus({of}) {

    const [ color, setColor ] = useState("gray-200")
    const [ status, setStatus ] = useState()

    useEffect(() => {

        // TODO wait for the backend to support single service status
        fetch("https://status.leon.wtf/api/services/latest")
            .then(r => r.json())
            .then(r => r.services[of].status)
            .then(r => {
                setColor(r === "up" ? "green-400" : "red-500")
                setStatus(r)
            })

    })

    return (
        <div className={`w-32 bg-${color} rounded-lg text-center`}>
            <p>{ status ? status : "loading" }</p>
        </div>
    )
}