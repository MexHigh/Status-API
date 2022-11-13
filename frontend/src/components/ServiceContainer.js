import React, { useState, useRef, useEffect } from "react"
import CurrentStatus from "./CurrentStatus"
import MiscEntries from "./MiscEntries"
import StatusPill from "./StatusPill"

// ServiceContainer bundles the service name, StatusPill and
// CurrentStatus components to create an availability timeline.
export default function ServiceContainer({ name, latest, timeline }) {

	const [ numOfPills, setNumOfPills ] = useState()
	const widthRef = useRef()

	useEffect(() => {
		// calulate number of pills on mount
		calculateNumberOfPills()
		// add event listener on mount
		window.addEventListener("resize", calculateNumberOfPills)
		return () => { // remove event listener on unmount
			window.removeEventListener("resize", calculateNumberOfPills)
		}
	// eslint-disable-next-line
	}, [])

	const calculateNumberOfPills = () => {
		if (widthRef) {
			setNumOfPills(
				Math.round(widthRef.current.clientWidth / 35)
			)
		} else {
			console.error("widthRef is", widthRef)
		}
	}

	const makePills = () => {
		let pills = []
		for (let i = 0; i < numOfPills; i++) {
			let status = timeline[i + (timeline.length - (numOfPills > 30 ? 30 : numOfPills))]
			if (status === undefined) {
				pills.push(
					<StatusPill 
						key={i+30}
						inactive={true} // do not show hover menu
					/>
				)
			} else {
				pills.push(
					<StatusPill 
						key={i + 30}
						forDay={status.at}
						status={status.status}
						availability={status.availability}
						downtimes={status.downtimes}
					/>
				)
			}
		}
		return pills
	}

	return (
		<div className="bg-white px-8 py-6 md:px-12 md:py-8 shadow-md rounded-lg" ref={widthRef}>
			{/* First line */}
			<div className="mb-4 flex justify-between bg-gray-100 rounded-lg">
				<a 
					className="text-xl ml-2 truncate" 
					href={latest.url} 
					target="_blank" 
					rel="noreferrer"
				>
					{name}
				</a>
				<CurrentStatus status={latest.status} />
			</div>
			{/* Second line */}
			<div className="flex justify-between h-8">
				{makePills()}
			</div>
			{/* Third line (optional) */}
			<MiscEntries misc={latest.misc} />
		</div>
	)
}
