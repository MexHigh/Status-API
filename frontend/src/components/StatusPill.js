import React from "react"

export default function StatusPill({ status }) {
	let color = "gray"

	if (status) {
		switch (status) {
			case "up":
				color = "green"
				break
			case "problems":
				color = "yellow"
				break
			case "down":
				color = "red"
				break
			default:
				color = "gray"
		}
	}

	return <div className={`w-4 h-8 bg-${color}-400 mr-4`}></div>
}
