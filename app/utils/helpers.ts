// $roots keeps previous parent properties as they will be added as a prefix for each prop.
// $sep is just a preference if you want to seperate nested paths other than dot.
const flatten = (
  obj: { [key: string]: any },
  roots: Array<string> = [],
  sep = '.'
): { [key: string]: any } =>
  Object
    // find props of given object
    .keys(obj)
    // return an object by iterating props
    .reduce(
      (memo, prop: string) =>
        Object.assign(
          // create a new object
          {},
          // include previously returned object
          memo,
          Object.prototype.toString.call(obj[prop]) === '[object Object]'
            ? // keep working if value is an object
              flatten(obj[prop], roots.concat([prop]), sep)
            : // include current prop and value and prefix prop with the roots
              { [roots.concat([prop]).join(sep)]: obj[prop] }
        ),
      {}
    );

const nestedValue = (obj: { [field: string]: any }, key: string): any => {
  if (!key || !obj) {
    return;
  }
  const splitted = key.split('.');
  if (splitted.length === 1) {
    return obj?.[key];
  }
  return nestedValue(obj?.[splitted[0]], splitted.slice(1).join('.'));
};

const setNestedValue = (key: string, values: { [field: string]: any }, value: any) => {
  const splitted = key.split('.');
  if (splitted.length > 1) {
    let parent = values[splitted[0]];
    if (!parent) {
      parent = {};
      values[splitted[0]] = parent;
    }
    for (let i = 1; i < splitted.length; i++) {
      if (i + 1 === splitted.length) {
        parent[splitted[i]] = value;
      } else {
        let v = parent[splitted[i] as any];
        if (!v) {
          v = {};
        }
        parent = v;
      }
    }
  } else {
    values[key] = value;
  }
};

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

export { flatten, nestedValue, setNestedValue, sleep };
