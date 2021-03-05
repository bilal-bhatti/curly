import collections

def update(d, u, extend_array=False):
  for k, v in u.items():
    if isinstance(v, collections.abc.Mapping):
      d[k] = update(d.get(k, {}), v)
    elif isinstance(v, list):
      if k in d and extend_array:
        d[k].extend(v)
      else:
        d[k] = v
    else:
      d[k] = v
  return d
