import React from "react"
import { createRoot } from 'react-dom/client'
import reportWebVitals from "./reportWebVitals"
import { createBrowserRouter, RouterProvider } from "react-router-dom"
import "./index.css" // this is where Tailwind CSS will be injected
import App from "./App"
import Admin from "./Admin"

// Workaround for CSP errors with FontAwesome:
// 
// Usually, FA applies CSS inline which violates the style-src
// CSP policy. This workaround disables automatic CSS application
// and imports the CSS file manually to be purged by webpack.
import { config } from "@fortawesome/fontawesome-svg-core"
import "../node_modules/@fortawesome/fontawesome-svg-core/styles.css"
config.autoAddCss = false

const root = createRoot(document.getElementById("root"))
const router = createBrowserRouter([
	{
		path: "/",
		element: <App />
	},
	{
		path: "/admin",
		element: <Admin />
	}
])

root.render(
	<React.StrictMode>
		<RouterProvider router={router} />
	</React.StrictMode>,
)

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals()
