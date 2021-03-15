import React from "react"
import CurrentStatus from "../components/CurrentStatus"
import StatusPill from "../components/StatusPill"

// ServiceContainer bundles the service name, StatusPill and
// CurrentStatus components to create an availability timeline.
export default function ServiceContainer({service}) {
    
    return(            
        <div className="p-4 mx-4 my-8 border-4 rounded-lg"> 
            {/* First line */}
            <div className="mb-4 flex justify-between bg-gray-100 rounded-lg">
                <h1 className="font-bold text-xl ml-2">{service}</h1>
                <CurrentStatus of={service}/>
            </div>
            {/* Second line */}
            <div className="flex">
                <StatusPill status="up"/>
                <StatusPill />
                <StatusPill status="down"/>
                <StatusPill status="problems"/>
            </div>
        </div>
    )

}