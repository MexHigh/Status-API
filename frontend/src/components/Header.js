import React from "react"

export default function Header({ lastCheckTs }) {

	let dateTimeString = new Date(lastCheckTs).toLocaleString(navigator.language, {
		year: "numeric", 
		month: "2-digit", 
		day: "2-digit", 
		hour: "2-digit", 
		minute: "2-digit", 
		second: "2-digit"
	})

	return (
		<div className="w-full text-center">
			<div className="mx-auto w-20 my-4">
				<a
					href="https://leon.wtf"
					target="_blank"
					rel="noreferrer"
				>
					<img 
						src="https://cdn.leon.wtf/icon/logo-email-sig.png" 
						alt="Header logo"
					/>
				</a>
			</div>
			<h1 className="text-xl">
				Status API for <span className="font-bold">leon.wtf</span>
			</h1>
			<h2 className="text-lg">
				Last update: <span className="font-bold">{dateTimeString}</span>
			</h2>
		</div>
	)
	
}
