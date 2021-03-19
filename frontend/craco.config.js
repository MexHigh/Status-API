const HtmlWebpackPlugin = require('html-webpack-plugin')
const CspHtmlWebpackPlugin = require('csp-html-webpack-plugin')

module.exports = {
	style: {
		postcss: {
			plugins: [require("tailwindcss"), require("autoprefixer")],
		},
	},
	webpack: {
		plugins: [
			new HtmlWebpackPlugin(),
			new CspHtmlWebpackPlugin({
					"default-src": "'self'",
					"img-src": "*"
			})
		]
	}
}
