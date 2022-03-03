const SmeeClient = require('smee-client')

console.log(process.env.GITHUB)
console.log(process.env.DOCKERHUB)

const smee1 = new SmeeClient({
    source: process.env.GITHUB,
    target: 'http://localhost:8080/api/v1/hook/github',
    logger: console
})
const smee2 = new SmeeClient({
    source: process.env.DOCKERHUB,
    target: 'http://localhost:8080/api/v1/hook/dockerhub',
    logger: console
})

const events1 = smee1.start()
const events2 = smee2.start()

// Stop forwarding events
// events.close()