let env = process.env.ENV || 'local';

module.exports = require(`../environments/${env}.json`);

console.log(44444444, module.exports);