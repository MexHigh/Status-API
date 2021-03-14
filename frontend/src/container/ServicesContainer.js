import React from "react"
import ServiceContainer from "./ServiceContainer"

// ServicesContainer bundles multiple ServiceContainer components.
// The whole fetch magic should propably happen in here.
export default function ServicesContainer() {
    return (
        <>
            <ServiceContainer />
            <p className="text-red-500">Nadine was here</p>
        </>

    )
}