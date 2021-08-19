import React from "react"

export default function Header({ lastCheckTs }) {
	/*return (
		<div className="flex justify-between items-center py-2 px-10">
			<div className="flex items-center w-2/3">
				<div className="w-20 mr-4">
					<img src="https://cdn.leon.wtf/icon/logo-email-sig.png" alt="Header logo" />
				</div>
				<div>
					<h1 className="text-xl">
						Status API for <span className="font-bold">leon.wtf</span>
					</h1>
				</div>
			</div>
			<div>
				<h1 className="text-xl">
					Last update: <span className="font-bold">{new Date(lastCheckTs).toLocaleString("de-DE")}</span>
				</h1>
			</div>
		</div>
	)*/
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
				Last update: <span className="font-bold">{new Date(lastCheckTs).toLocaleString("de-DE")}</span>
			</h2>
		</div>
	)
}
