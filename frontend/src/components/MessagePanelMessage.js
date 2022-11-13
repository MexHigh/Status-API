import React from "react"
import moment from "moment"

export default function Message({ title, status, content, updated }) {
    return (
        <div className="bg-white my-4 py-6 px-12 rounded-lg shadow-md">
            <details>
                <summary className="cursor-pointer">
                    <span className="select-none text-lg">{title}</span>
                    <span className="text-gray-400 float-right">
                        {moment(updated).fromNow()}
                    </span>
                </summary>
                <div className="p-4 flex flex-col gap-2">
                    <p className="font-bold ">
                        {status}
                    </p>
                    <p>Last update: {moment(updated).calendar()}</p>
                    <p>{content || "This entry has not content"}</p>
                </div>
            </details>
        </div>
    )
}