import React from "react"
import CurrentStatus from "../components/CurrentStatus"
import StatusPill from "../components/StatusPill"

// ServiceContainer bundles the service name, StatusPill and
// CurrentStatus components to create an availability timeline.
export default function ServiceContainer({ name, latest, timeline }) {

	const makePills = () => {
		let pills = []
		for (let i = 30 - timeline.length; i > 0; i--) {
			pills.push(
				<StatusPill key={i} />
			)
		}
		timeline.forEach(day => {
			pills.push(
				<StatusPill key={day.at} status={day.status} />
			)
		})
		return pills
	}

	return (
		<div className="p-4 mx-auto my-8 border-4 rounded-lg w-5/6 max-w-7xl">
			{/* First line */}
			<div className="mb-4 flex justify-between bg-gray-100 rounded-lg">
				<h1 className="font-bold text-xl ml-2">{name}</h1>
				<CurrentStatus status={latest.status} />
			</div>
			{/* Second line */}
			<div className="flex justify-between">
				{makePills()}
			</div>
		</div>
	)

}
