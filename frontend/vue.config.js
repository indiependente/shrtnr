module.exports = {
    publicPath: process.env.NODE_ENV === 'production'
        ? 'https://shr.tn/'
        : 'http://localhost:7000/'
}
