import moment from "moment"
import React from "react"

export default function Header({ title, logoURL, lastCheckTs }) {
	const dateTimeString = moment(lastCheckTs).calendar()

	return (
		<div className="w-full text-center">
			<div className="mx-auto w-20 my-4">
				<img 
					src={logoURL}
					alt="Header logo"
				/>
			</div>
			<h1 className="text-xl">
				{title}
			</h1>
			<h2 className="text-lg">
				Last update: <span className="font-bold">{dateTimeString}</span>
			</h2>
		</div>
	)
	
}
