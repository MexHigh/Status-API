import React, { useState } from "react"

export default function StatusPill({ status, downtimes }) {

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
			<div className={`w-4 h-8 ${color} mx-1 rounded-lg`}></div>

			{/* Hover Menu */}
			{ hovering && 
				<div className="absolute top-10 w-32 h-16 shadow-xl bg-gray-100 rounded-lg">
					
				</div> 
			}
			
		</div>

	)

}
