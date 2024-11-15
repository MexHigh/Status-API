import React, { useState, useEffect } from "react"
import Footer from "./components/Footer"
import Header from "./components/Header"
import Loading from "./components/Loading"
import StatusSummary from "./components/StatusSummary"
import ServiceContainer from "./components/ServiceContainer"
import MessagePanel from "./components/MessagePanel"

export default function App() {
	const [title, setTitle] = useState()
	const [logoURL, setLogoURL] = useState()
	const [latest, setLatest] = useState()
	const [timeline, setTimeline] = useState()

	const getTitleAndLogo = () => {
		fetch("/api/dashboard/title")
			.then(r => r.json())
			.then(r => {
				setTitle(r.response)
			})
		fetch("/api/dashboard/logo")
			.then(r => r.blob())
			.then(r => {
				let url = URL.createObjectURL(r)
				// Preload an invisible image to prevent flickering.
				// Instead, let the image preload and set the state afterwards
				// (which triggers the conditional re-render of the main page).
				let img = new Image()
				img.src = url
				img.onload = () => {
					setLogoURL(url)
				}
			})
	}

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
		getTitleAndLogo()
		fetchApi()
		const interval = setInterval(() => {
			fetchApi()
		}, 30 * 1000 /* seconds */)
		return () => clearInterval(interval)
	}, [])

	if (!title || !logoURL || !latest || !timeline) {
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
					<Header 
						title={title}
						logoURL={logoURL}
						lastCheckTs={latest.at}
					/>
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