module.exports = {
    dev: {},
    prod: {
        "default-src": "'self'",
        "style-src": [
            "'self'",
            "'unsafe-inline'",
            "data:"
        ],
        "image-src": [
            "*"
        ]
    }
}