import React from "react"

export default function CurrentStatus({ status }) {

	let color = status === "up" ? "green-400" : "red-400"

	return (
		<div className={`w-32 bg-${color} rounded-lg text-center`}>
			<p>{status}</p>
		</div>
	)
	
}
