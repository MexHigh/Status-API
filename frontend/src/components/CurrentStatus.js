import React from "react"
import { usePopperTooltip } from "react-popper-tooltip"
import StatusHoverMenu from "./StatusHoverMenu"

/**
 * Shows the current status of a service as a colored bar. Optionally,
 * hovering the bar shows the same hover menu the status pills use,
 * e.g. to observe todays downtimes.
 * @param {string} status One of "up"|"problems"|"down"
 * @param {string} forDay Timestamp of the day the optional hover menu shows details for (no hover menu, if undefined)
 * @param {float} availability The availabilty of the service in percent for that day
 * @param {Array} downtimes Array of downtimes for that day
 */
export default function CurrentStatus({ status, forDay, availability, downtimes }) {

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
			color = "bg-green-400"
			break
		case "problems":
			color = "bg-yellow-400"
			break
		case "down":
			color = "bg-red-400"
			break
		default: // or !status
			color = "bg-gray-200"
	}

	return (
		<div>
			<div ref={setTriggerRef} className={`w-20 md:w-32 h-full ${color} rounded-lg text-center select-none`}>
				<p>{status}</p>
			</div>
			{/* Hover Menu */}
			{ (forDay && visible) && (
				<div
					ref={setTooltipRef}
					{...getTooltipProps()}
				>
					<StatusHoverMenu
						forDay={forDay}
						status={status}
						availability={availability}
						downtimes={downtimes && downtimes.slice(0, 30)}
					/>
				</div>
			)}
		</div>
	)

}