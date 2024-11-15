import React, { useState, useEffect } from "react"
import Footer from "./components/Footer"
import Header from "./components/Header"
import Loading from "./components/Loading"
import LoginForm from "./components/LoginForm"
import WidthContainer from "./components/WidthContainer"

export default function Admin() {
	const [title, setTitle] = useState()
	const [logoURL, setLogoURL] = useState()
    const [apiToken, setApiToken] = useState(undefined) // once loaded either "null" or "string"

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

    const tryGetApiToken = () => {
        let key = localStorage.getItem('status_api_key')
        setApiToken(key) // key is null if not found
    }

	useEffect(() => {
		getTitleAndLogo()
        tryGetApiToken()
	}, [])

	if (!title || !logoURL || apiToken === undefined) {
		return <Loading />
	} else {
		return (
			<>
				<header id="header">
					<WidthContainer
						className="mb-8"
					>
						<Header 
							title={title}
							logoURL={logoURL}
							supplement="Admin Dashboard"
						/>
					</WidthContainer>
				</header>
				<main className="w-11/12 md:w-5/6 mx-auto">
					{ apiToken === null ? (
                        <LoginForm />
                    ) : (
                        <div>
                            <p>Logged in!</p>
                        </div>
                    )}
				</main>
				<footer id="footer">
					<Footer />
				</footer>
			</>
		)
	}
}