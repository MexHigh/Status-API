import moment from "moment"
import React from "react"
import CurrentStatus from "./CurrentStatus"

export default function StatusHoverMenu({ forDay, status, availability, downtimes }) {

	let forDaydateTimeString = moment(forDay).format("LL")

    return (
		<div className="w-80 px-4 py-2 bg-gray-100 rounded-lg shadow-md">
			<div className="m-3">
				{/* First line (status and date) */}
				<div className="flex justify-evenly">
					<CurrentStatus status={status}/>
					<p>on</p>
					<h1 className="font-bold text-center">
						{ forDaydateTimeString }
					</h1>
				</div>
				{/* Second line (availability) */}
				<div className="flex justify-evenly mt-2">
					<p>Availability: <span>{(availability*100).toFixed(2)}%</span></p>
				</div>
			</div>
			{/* downtime list */}
			{ downtimes && downtimes.length > 0 && 
				<h2>Downtimes: </h2> 
			}
			{ downtimes && 
				downtimes.map(downtime => 
					<Downtime
						key={downtime.from}
						from={downtime.from}
						to={downtime.to}
						reason={downtime.reason}
					/>
				)
			}
		</div> 
    )
}

function Downtime({ from, to, reason }) {

	const fromTime = new Date(from).toLocaleTimeString("de-DE")
	const toTime = new Date(to).toLocaleTimeString("de-DE")

	return (
		<div className="mt-2 whitespace-nowrap overflow-hidden">
			<h1 className="italic">
				{ fromTime === toTime ? <>{fromTime}</> : <>{fromTime} - {toTime}</> }
			</h1>
			<h2 className="text-gray-600">{reason || "Reason unknown"}</h2>
		</div>
	)
}