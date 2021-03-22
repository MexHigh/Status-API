import React from "react"

export default function StatusPill({ status }) {
	
	let color

	switch (status) {
		case "up":
			color = "bg-green-300"
			break
		case "problems":
			color = "bg-yellow-300"
			break
		case "down":
			color = "bg-red-300"
			break
		default: // or !status
			color = "bg-gray-300"
	}

	return <div className={`w-4 h-8 ${color} mx-1 rounded-lg`}></div>

}
