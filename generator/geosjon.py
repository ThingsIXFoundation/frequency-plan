import json
import geojson
import h3.api.numpy_int as h3
import numpy as np

def fix_transmeridian(loop):
    isTransmerdian = False
    for i in range(len(loop)):
        if loop[i][0]-loop[(i+1) % len(loop)][0] > 180.0:
            isTransmerdian = True
            break

    if not isTransmerdian:
        return loop
    
    for i in range(len(loop)):
        if loop[i][0] < 0.0:
            loop[i] = (loop[i][0]+360.0, loop[i][1])
    return loop

def fix_transmeridian_multipoly(multipoly):
    for i in range(len(multipoly)):
        for j in range(len(multipoly[i])):
            multipoly[i][j] = fix_transmeridian(multipoly[i][j])

features = []
for plan in ["EU868", "AU915"]:
    cells = np.fromfile(file=f"../go/frequency_plan/{plan}.h3", dtype=np.int64, sep="", count=-1, offset=0)
    multipoly=h3.h3_set_to_multi_polygon(h3.uncompact(cells, 6), geo_json=True)
    fix_transmeridian_multipoly(multipoly)
    gmp = geojson.MultiPolygon(multipoly, validate=True)
    feature = geojson.Feature(id=plan, geometry=gmp, properties={"Plan": plan})
    features.append(feature)

fc = geojson.FeatureCollection(features)
f = open("plans.geojson", "w")
f.write(json.dumps(fc))
