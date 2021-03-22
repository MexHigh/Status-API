import React from "react"

export default function CurrentStatus({ status }) {

	let color = status === "up" ? "bg-green-400" : "bg-red-400"

	return (
		<div className={`w-32 ${color} rounded-lg text-center`}>
			<p>{status}</p>
		</div>
	)
	
}
