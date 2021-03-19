import React from "react"

export default function Header({ lastCheckTs }) {
	return (
		<header className="bg-white shadow m-4 rounded-full">
			<div className="flex justify-between py-2 px-10">
                {/* Left side */}
                <div className="flex w-1/3">
                    {/* TODO */}
                    <img className="w-10" src="https://cdn.leon.wtf/icon/logo-header.svg" alt="Header logo" />
				    <p className="text-red-500">Status API</p>
                </div>
                {/* Right side */}
                <div>
                    <p className="text-red-500">
                        {new Date(lastCheckTs).toLocaleString("de-DE")}
                    </p>
                </div>
			</div>
		</header>
	)
}
