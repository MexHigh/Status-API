import React from "react"
import { usePopperTooltip } from "react-popper-tooltip"
import StatusHoverMenu from "./StatusHoverMenu"

export default function StatusPill({ forDay, status, availability, downtimes }) {

	const { setTriggerRef, setTooltipRef, getTooltipProps, visible } = usePopperTooltip({
		offset: [0, 12],
		trigger: "hover",
		interactive: true,
		delayHide: 50,
		delayShow: 100,
	})

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
		<div>
			{/* Status Pill */}
			<div
				ref={setTriggerRef}
				className={`w-4 h-8 ${color} mx-1 rounded-xl`} 
			/>

			{/* Hover Menu */}
			{ visible && ( 
				<div 
					ref={setTooltipRef}
					{...getTooltipProps()}
				>
					<StatusHoverMenu
						forDay={forDay}
						status={status}
						availability={availability}
						// Just take the first 30 downtimes, if there is more, than this may be a bug anyway
						downtimes={downtimes && downtimes.slice(0, 30)}
					/>
				</div>
			)}
		</div>
	)

}
