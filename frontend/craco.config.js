const HtmlWebpackPlugin = require('html-webpack-plugin')
const CspHtmlWebpackPlugin = require('csp-html-webpack-plugin')

let myPlugins = []

/**
 * Add the CSP Plugin only when creating production builds
 */
if (process.env.NODE_ENV === "production") {
	plugins.push(
		new HtmlWebpackPlugin(),
		new CspHtmlWebpackPlugin({
			"default-src": "'self'",
			"style-src": "'self' 'unsafe-inline' data:",
			"img-src": "*"
		})
	)
}

module.exports = {
	style: {
		postcss: {
			plugins: [require("tailwindcss"), require("autoprefixer")],
		},
	},
	webpack: {
		plugins: myPlugins
	}
}
