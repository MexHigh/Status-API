import React from "react"

export default function StatusSummary({ latest }) {

    let counter = {
        ups: 0,
        downs: 0
    }
    Object.entries(latest.services).forEach(([serviceName, latestStatus]) => {
        latestStatus.status === "up" ? counter.ups++ : counter.downs++
    })

    let bgColor = "bg-gray-200"
    let text = ""
    if (counter.downs === 0) {
        bgColor = "bg-green-400"
        text = "All services operational"
    } else if (counter.downs < counter.ups) {
        bgColor = "bg-yellow-400"
        text = "Some services are unavailable"
    } else {
        bgColor = "bg-red-400"
        text = "Most services are unavailable"
    }

    return (
        <div className={`rounded-xl shadow-lg w-max py-4 px-10 flex items-center space-x-4`}>
            <div className={`${bgColor} w-5 h-5 rounded-full`}>
                <div className={`${bgColor} w-full h-full rounded-lg animate-ping`}></div>
            </div>
            <h1 className="text-center font-semibold">{text}</h1>
        </div>
    )
}