import React, { useState, useEffect } from "react"
import Header from "../components/Header"
import Loading from "../components/Loading"
import ServiceContainer from "./ServiceContainer"

// ServicesContainer bundles multiple ServiceContainer components.
// The whole fetch magic should propably happen in here.
export default function ServicesContainer() {

	const [latest, setLatest] = useState()
	const [timeline, setTimeline] = useState()

	useEffect(() => {

		fetch("https://status.leon.wtf/api/services/latest")
			.then(r => r.json())
			.then(r => {
				setLatest(r)
			})

		fetch("https://status.leon.wtf/api/services/timeline")
			.then(r => r.json())
			.then(r => {
				setTimeline(r)
			})

	}, [])

	if (!latest || !timeline) {
		return <Loading />
	} else {

		let serviceTimeline = {}
		timeline.forEach(day => {
			for (const [name, status] of Object.entries(day.services)) {
				if (!serviceTimeline[name]) {
					serviceTimeline[name] = []
				}
				serviceTimeline[name].push({
					at: day.at,
					...status
				})
			}
		});

		return (
			<div>
				<Header lastCheckTs={latest.at} />
				{Object.entries(latest.services).map(
					([serviceName, latestStatus]) => (
						<ServiceContainer
							key={serviceName}
							name={serviceName}
							latest={latestStatus}
							timeline={serviceTimeline[serviceName]}
						/>
					)
				)}
			</div>
		)

	}

}
