// exports.int        = int;
//
// exports.nn         = nn;
// exports.str        = str;
// exports.cut        = cut;
// exports.escapeXml  = escapeXml;  // npm i html-entities
//
// exports.array      = array;
// exports.first      = first;
// exports.second     = second;
// exports.order      = order;
//
// exports.object     = object;
// exports.copy       = copy;

export default {
  int, nn, str, cut, escapeXml, array, first, second, order, object, copy, dateStr, i, ii, iii,
}

function dateStr(d) {
    if (typeof d === "string") {
        d = new Date(d);
    }

    if (d instanceof Date) {
        let month = d.getMonth(); if (month < 10) month = "0" + month;
        let day = d.getDate(); if (day < 10) day = "0" + day;
        let hours = d.getHours(); if (hours == 0) { hours = "00"; } else if (hours < 10 ) hours = "0" + hours;
        let minutes = d.getMinutes(); if (minutes == 0) { minutes = "00"; } else if (minutes < 10 ) minutes = "0" + minutes;

        return d.getFullYear() + "-" + month + "-" + day + " " + hours + ":" + minutes
    }

    // if (typeof d === "string") {
    //     return d.substr(0, 16).replace("T", " ");
    // }

    return "";
}

function int(A) {
  if (typeof A === 'string') {
    A = A.replace(/^\s+/, '').replace(/\s+$/, '');
    if (A.match(/\d\D/)) {
      return 0;
    }
  } else if (A instanceof Object) {
    return 0;
  }
  return isNaN(A = parseInt(A)) ? 0 : A;
}

function nn(A) {
  return A === undefined || A === null ? '' : A;
}

function str(A) {
  return A === undefined || A === null ? '' : '' + A;
}

function cut(A, maxLength, rest) {
  A    = str(A);
  rest = str(rest) || '...';
  if (!(maxLength >= rest.length)) {
    maxLength = rest.length;
  }

  return A.length > maxLength ? A.substr(0, maxLength - rest.length) + rest : A;
}

function escapeXml(t) {
  return str(t)
         .replace(/"/g, '&quot;')
         .replace(/</g, '&lt;')
         .replace(/>/g, '&gt;');
}

function array(A) {
  return A instanceof Array              ? A
         : A === undefined || A === null ? []
         : [A];
}

function first(A) {
  return A instanceof Array ? A[0] : A;
}

function second(A) {
  return A instanceof Array ? A[1] : undefined;
}

function order(list, order, addRest) {
  list = array(list);
  let res = addRest
             ? array(order)
             : array(order).filter(e => list.includes(e));
  return res.concat(list.filter(e => !res.includes(e)));
}

function object(A) {
  return A instanceof Object ? A
       : A === undefined || A === null ? {}
       : {'': A};
}

function copy(A) {
  if (A instanceof Array) {
    return A.map(e => copy(e));
  } else if (A instanceof Object) {
    let A_ = {};
    
    for (let key of Object.keys(A)) {
      // for (let key in A) {
      //   if (A.hasOwnProperty(key)) {

      A_[key] = copy(A[key]);
    }
    return A_;
  }
  return A;
}

// //////////////////////////////////////////////

const i =  (...data) => i_(1111111111111,  ...data);
const ii =  (...data) => i_(222222222222,  ...data);
const iii =  (...data) => i_(33333333333,  ...data);
const iiii =  (...data) => i_(4444444444,  ...data);
const iiiii =  (...data) => i_(555555555,  ...data);
const iiiiii =  (...data) => i_(66666666,  ...data);
const iiiiiii =  (...data) => i_(7777777,  ...data);
const iiiiiiii =  (...data) => i_(888888,  ...data);
const iiiiiiiii =  (...data) => i_(99999,  ...data);

function i_(...data_) {
    let data = [];
    data_.map(d => data.push(...['//', d]));
    console.log(...data);
}

