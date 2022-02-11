const SmeeClient = require('smee-client')

const smee = new SmeeClient({
    source: 'https://smee.io/xzMYO1EMlJciCLhw',
    target: 'http://localhost:8080/api/v1/hook/github',
    logger: console
})

const events = smee.start()

// Stop forwarding events
// events.close()