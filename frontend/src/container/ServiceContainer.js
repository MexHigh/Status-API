import React from "react"
import CurrentStatus from "../components/CurrentStatus"
import StatusPill from "../components/StatusPill"

// ServiceContainer bundles the serice name, StatusPill and
// CurrentStatus components to create an availability timeline.
export default function ServiceContainer({service}) {
    
    return(            
        <div className="p-4 mx-4 my-8 border-4 rounded-lg"> 
            <div className="mb-4 flex justify-between">
                <h1 className="font-bold text-xl">{service}</h1>
                <CurrentStatus of={service}/>
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