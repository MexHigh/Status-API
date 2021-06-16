import React from "react"

export default function Header({ lastCheckTs }) {
	return (
		<div className="flex justify-between items-center py-2 px-10">
			{/* Left side */}
			<div className="flex items-center w-1/3">
				<div className="w-20 mr-4">
					<img src="https://cdn.leon.wtf/icon/logo-email-sig.png" alt="Header logo" />
				</div>
				<div>
					<h1 className="text-xl">
						Status API for <span className="font-bold">leon.wtf</span>
					</h1>
				</div>
			</div>
			{/* Right side */}
			<div>
				<h1 className="text-xl">
					Last update: <span className="font-bold">{new Date(lastCheckTs).toLocaleString("de-DE")}</span>
				</h1>
			</div>
		</div>
	)
}
