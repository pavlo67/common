exports.correctPath = correctPath;

// -----------------------------------------------------------------------------

const caller = require('caller');
const f      = require('./basis');

function correctPath(path, caller__) {
  let [l, caller_] = f.array(path);

  l = f.str(l);
  l = l.match(/^([\/\\]|\w\:)/)
    ? l
    : (f.str(caller_ || caller__ || caller()).replace(/[\/\\][^\/\\]*$/, '/') + l);

  return l.replace(/\\/g, '/').replace(/\/\.\//g, '/');
}
