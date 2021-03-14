import React from "react"
import CurrentStatus from "../components/CurrentStatus"
import StatusPill from "../components/StatusPill"

// ServiceContainer bundles the serice name, StatusPill and
// CurrentStatus components to create an availability timeline.
export default function ServiceContainer({service}) {
    return(

            
        <div className="p-4 m-12 border-4 rounded-lg"> 
            <div className="mb-4 flex justify-between">
                <h1>Service {service}</h1>
                <CurrentStatus />
            </div>
            <div className="flex">
                <StatusPill status="up"/>
                <StatusPill />
                <StatusPill status="down"/>
                <StatusPill status="problems"/>
            </div>
        </div>


    )
}