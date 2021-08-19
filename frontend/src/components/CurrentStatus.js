import React from "react"

export default function CurrentStatus({ status }) {

	let color
	switch (status) {
		case "up":
			color = "bg-green-400"
			break
		case "problems":
			color = "bg-yellow-400"
			break
		case "down":
			color = "bg-red-400"
			break
		default: // or !status
			color = "bg-gray-200"
	}

	return (
		<div className={`w-20 md:w-32 ${color} rounded-lg text-center select-none`}>
			<p>{status}</p>
		</div>
	)
	
}
