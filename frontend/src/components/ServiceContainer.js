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
		// calulate number of pills
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
			let tempNumOfPills = Math.round(widthRef.current.clientWidth / 35)
			console.log("Rendering", tempNumOfPills, "pills for", name)
			setNumOfPills(tempNumOfPills)
		} else {
			console.error("widthRef is", widthRef)
		}
	}

	const makePills = () => {
		let pills = []
		
		let trimmedTimeline = timeline.slice(
			timeline.length - (numOfPills > 30 ? 30 : numOfPills), 
			timeline.length
		)

		trimmedTimeline.forEach((day, i) => {
			pills.push(
				<StatusPill 
					key={i + 30}
					forDay={day.at}
					status={day.status}
					availability={day.availability}
					downtimes={day.downtimes}
				/>
			)
		})

		return pills // should have a length of 30

	}

	return (
		<div className="px-12 py-8 shadow-lg rounded-lg" ref={widthRef}>
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
