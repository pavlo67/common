function convert(sw) {
  if (!(sw instanceof Object && sw.paths instanceof Object)) return {};

  let swc = {
    host: sw.port, // "http://localhost" +
    endpoints: {},
  };


  for (let path in sw.paths) {
    let p = sw.paths[path];

    if (!(p instanceof Object)) {
      // TODO: signal the error
      continue;
    }

    // console.log(path, p[path], p[path] instanceof Object);

    for (let method in p) {
      if (!(p[method] instanceof Object)) {
        // TODO: signal the error
        continue;
      } else if (!p[method].operationId) {
        // TODO: signal the error
        continue;
      }
      let operationId = p[method].operationId;
      delete p[method].operationId;

      p[method].path   = path;
      p[method].method = method;

      swc.endpoints[operationId] = p[method];
    }
  }

  return swc;
}

function ep(swc, key) {
  if (!(swc instanceof Object && swc.endpoints instanceof Object && swc.endpoints[key] instanceof Object)) {
    return "";
  }

  return window.location.protocol + "//" + window.location.hostname + swc.host + swc.endpoints[key].path;
}

export { convert, ep } ;
