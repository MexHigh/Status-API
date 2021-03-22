import React from "react"

export default function Footer() {
	return (
		<footer className="text-center p-16">
			<p className="text-gray-800">
				&copy; {new Date().getFullYear()} Leon Schmidt at <span className="font-bold">leon.wtf</span>
			</p>
		</footer>
	)
}
