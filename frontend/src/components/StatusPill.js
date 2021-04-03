import React, { useState } from "react"
import StatusHoverMenu from "./StatusHoverMenu"

export default function StatusPill({ forDay, status, availability, downtimes }) {

	const [ hovering, setHovering ] = useState(false)
	
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
			color = "bg-gray-200"
	}

	return (

		<div 
			onMouseEnter={ () => { setHovering(true) } }
			onMouseLeave={ () => { setHovering(false) } }
			className="relative flex flex-col items-center"
		>

			{/* Status Pill */}
			<div className={`w-4 h-8 ${color} mx-1 rounded-xl`}></div>

			{/* Hover Menu */}
			{ hovering && 
				<StatusHoverMenu
					forDay={forDay}
					status={status}
					availability={availability}
					downtimes={downtimes}
				/>
			}
			
		</div>

	)

}
