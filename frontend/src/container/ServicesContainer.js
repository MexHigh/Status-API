import React, { useState, useEffect } from "react"
import Footer from "../components/Footer"
import Header from "../components/Header"
import Loading from "../components/Loading"
import ServiceContainer from "./ServiceContainer"

// ServicesContainer bundles multiple ServiceContainer components.
// The whole fetch magic should propably happen in here.
export default function ServicesContainer() {
	
	const [latest, setLatest] = useState()
	const [timeline, setTimeline] = useState()

	useEffect(() => {

		fetch("/api/services/latest")
			.then(r => r.json())
			.then(r => {
				setLatest(r)
			})

		fetch("/api/services/timeline")
			.then(r => r.json())
			.then(r => {
				setTimeline(r)
			})

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
			<div>
				<Header lastCheckTs={latest.at} />
				{
					// map over the latest service report to calculate
					// the number of ServiceContainer components
					Object.entries(latest.services).map(([serviceName, latestStatus]) => (
						<ServiceContainer
							key={serviceName}
							name={serviceName}
							latest={latestStatus}
							timeline={serviceTimeline[serviceName]}
						/>
					))
				}
				<Footer />
			</div>
		)
	}
}
