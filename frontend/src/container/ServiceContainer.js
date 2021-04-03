import React from "react"
import CurrentStatus from "../components/CurrentStatus"
import MiscEntries from "../components/MiscEntries"
import StatusPill from "../components/StatusPill"

// ServiceContainer bundles the service name, StatusPill and
// CurrentStatus components to create an availability timeline.
export default function ServiceContainer({ name, latest, timeline }) {

	const makePills = () => {

		let pills = []

		// count the timeline entries and add as many grey
		// StatusPills so that there are 30 in total
		for (let i = 30 - timeline.length; i > 0; i--) {
			pills.push(
				<StatusPill key={i} />
			)
		}

		// add the actual status entry pills
		timeline.forEach((day, i) => {
			pills.push(
				<StatusPill 
					key={i + 30} 
					status={day.status}
					downtimes={day.downtimes}
				/>
			)
		})

		return pills // should have a length of 30

	}

	return (
		<div className="px-12 py-8 mx-auto my-12 w-5/6 max-w-5xl shadow-lg rounded-lg">
			{/* First line */}
			<div className="mb-4 flex justify-between bg-gray-100 rounded-lg">
				<a 
					className="font-bold text-xl ml-2" 
					href={latest.url} 
					target="_blank" 
					rel="noreferrer"
				>
					{name}
				</a>
				<CurrentStatus status={latest.status} />
			</div>
			{/* Second line */}
			<div className="flex justify-between">
				{makePills()}
			</div>
			{/* Third line (optional) */}
			<MiscEntries misc={latest.misc} />
		</div>
	)
}
