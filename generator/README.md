# Frequency Plan Generator

This generator generates a number of outputs based on a per-country definition of the frequency plan in the go package. Usually
there's no need to run the generator yourself as the outputs are committed in this repository too.

These outputs are described below:

### H3-index
In theory there should be 1 frequency-plan per location. However because of ill-defined borders multiple frequency-plans could be valid for a certain location. To H3 index allows for easy and quick look-up the valid frequency-plan(s) for a certain location. 

## How to run

### Download required tooling
Install [ogr2ogr](https://gdal.org/programs/ogr2ogr.html) of the gdal suite.
- For OSX: `brew install gdal`
- For Ubuntu: `apt install gdal-bin`

Also make sure Go is installed. 

### Download the EEZ definition file
Download the EEZ (Economic Exclusivity Zone) shape file from the download page the [Martine Regions](https://www.marineregions.org/downloads.php) website of the Flanders Marine Institute. You want the `World EEZ v11` shape-file. A form has to be filled. Extract the zip retrieved and copy all the `eez_v11.*` files into this folder.

### Convert the EEZ shape to GeoJSON
Run the following command in this folder to convert the shape-file into a GeoJSON file that can be read by this generator (this will take a while):
```
ogr2ogr -f GeoJSON -t_srs crs:84 eez_v11.geojson eez_v11.shp
```

Afterwards a `eez_v11.geojson` file should be created.

### Run the generator
Run the following command to generate all the outputs described above:
```
go run main.go
```

## Licenses
This generator uses EEZ map data from: 

Flanders Marine Institute (2019). Maritime Boundaries Geodatabase: Maritime Boundaries and Exclusive Economic Zones (200NM), version 11. Available online at https://www.marineregions.org/ https://doi.org/10.14284/386.

This data has a [CC BY 4.0 license](https://creativecommons.org/licenses/by/4.0/)

This generator uses country boundaries from OpenStreetMap:

OpenStreetMap® is open data, licensed under the [Open Data Commons Open Database License](https://opendatacommons.org/licenses/odbl/) (ODbL) by the [OpenStreetMap Foundation](https://osmfoundation.org/) (OSMF).

This data has a [CC BY-SA 2.0 license](https://creativecommons.org/licenses/by-sa/2.0/)