import React, { useState, useEffect } from "react"
import Loading from "../components/Loading"
import ServiceContainer from "./ServiceContainer"

// ServicesContainer bundles multiple ServiceContainer components.
// The whole fetch magic should propably happen in here.
export default function ServicesContainer() {

    const [ data, setData ] = useState()

    useEffect(() => {
/*
        fetch("https://status.leon.wtf/api/services/latest")
            .then(r => r.json())
            .then(r => {
                setData(r)
            })
*/
    }, [])

    if (!data) {
        return <Loading />
    }
    else {
        return (
            <div>
                <h1>Last check: {data.at}</h1>
                {
                    Object.entries(data.services).map(([key, value]) => 
                        <ServiceContainer key={key} service={key} />
                    )
                }
            </div>
        )
    }

}