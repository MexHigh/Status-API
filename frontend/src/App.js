import React, { useState, useEffect } from "react"
import Footer from "./components/Footer"
import Header from "./components/Header"
import Loading from "./components/Loading"
import StatusSummary from "./components/StatusSummary"
import ServiceContainer from "./components/ServiceContainer"
import MessagePanel from "./components/MessagePanel"

// Workaround for CSP errors with FontAwesome:
// 
// Usually, FA applies CSS inline which violates the style-src
// CSP policy. This workaround disables automatic CSS application
// and imports the CSS file manually to be purged by webpack.
import { config } from "@fortawesome/fontawesome-svg-core"
import "../node_modules/@fortawesome/fontawesome-svg-core/styles.css"
config.autoAddCss = false

export default function App() {

	const [latest, setLatest] = useState()
	const [timeline, setTimeline] = useState()

	const fetchApi = () => {
		fetch("/api/services/latest")
			.then(r => r.json())
			.then(r => {
				setLatest(r.response)
			})
		fetch("/api/services/timeline")
			.then(r => r.json())
			.then(r => {
				setTimeline(r.response)
			})
	}

	useEffect(() => {
		fetchApi()
		const interval = setInterval(() => {
			fetchApi()
		}, 30 * 1000 /* seconds */)
		return () => clearInterval(interval)
	}, [])

	if (!latest || !timeline) {
		return <Loading />
	} else {
		// restructure the response from /api/services/timeline
		// so that every ServiceContainer component receives only
		// it's own timeline array containing only one service
		let serviceTimeline = {}
		timeline.forEach(day => {
			// iterates over every service reported in one day
			for (const [name, status] of Object.entries(day.services)) {
				// create the timeline array, if it does not exist
				if (!serviceTimeline[name]) {
					serviceTimeline[name] = []
				}
				// append the status for a service and
				// slice in the "at" timestamp
				serviceTimeline[name].push({
					at: day.at,
					...status,
				})
			}
		})

		return (
			<>
				<header id="header" className="mx-auto max-w-5xl mb-8">
					<Header lastCheckTs={latest.at} />
				</header>
				<main className="w-11/12 md:w-5/6 mx-auto">
					<div id="status-summary" className="mx-auto w-max">
						<StatusSummary latest={latest} />
					</div>
					<div id="messages">
						<MessagePanel />
					</div>
					<div id="services">
						{
							// map over the latest service report to calculate
							// the number of ServiceContainer components
							Object.entries(latest.services).map(([serviceName, latestStatus]) => (
								<div 
									key={serviceName} 
									className="mx-auto max-w-5xl my-8"
								>
									<ServiceContainer
										name={serviceName}
										latest={latestStatus}
										timeline={serviceTimeline[serviceName] || []}
									/>
								</div>
							))
						}
					</div>
				</main>
				<footer id="footer">
					<Footer />
				</footer>
			</>
		)
	}
}