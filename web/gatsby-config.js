/**
 * Configure your Gatsby site with this file.
 *
 * See: https://www.gatsbyjs.org/docs/gatsby-config/
 */

var cors = require('cors')
module.exports = {
  /* Your site config here */
    developMiddleware: app => {
        app.use(cors())
    },
    plugins: [`gatsby-plugin-react-helmet`],
    proxy: {
        prefix: "/api",
        url: "http://localhost:8080",
    },
}
