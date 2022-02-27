const SmeeClient = require('smee-client')

const smee1 = new SmeeClient({
    source: 'https://smee.io/xzMYO1EMlJciCLhw',
    target: 'http://localhost:8080/api/v1/hook/github',
    logger: console
})
const smee2 = new SmeeClient({
    source: 'https://smee.io/VfwJDdPvLsJZjr5o',
    target: 'http://localhost:8080/api/v1/hook/dockerhub',
    logger: console
})

const events1 = smee1.start()
const events2 = smee2.start()

// Stop forwarding events
// events.close()