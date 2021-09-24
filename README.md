# Hygieia

Hygieia is the Ancient Egyptian god of cleanliness.  This project helps enable VATSIM Facility Engineers to easily filter out unwanted lines within their diagrams.

## Configuration

The configuration is a YAML file.

```yaml
filter:
  type: "radius"
  direction: "inside"
radius:
  center:
    lat: 40.7128
    lon: -74.0059
  radius: 
    distance: 60
    unit: nm
points:
- lat: 43.434109
  lon: -88.890204
- lat: 44.375058
  lon: -88.628186
- lat: 44.411986
  lon: -89.561127
- lat: 43.425246
  lon: -89.966063
- lat: 43.434109
  lon: -88.890204
```

### Filter

This defines the filter to be used on the map

filter.type accepted values:
* radius
* polygon

filter.direction accepted values:
* inside
* outside

### Radius

Define the radius filter

radius.center:
* lat defines the latitude of the center of the circle (decimal degrees represented as a float)
* lon defines the longitude of the center of the circle

radius.radius:
* distance
* unit (expected: km, nm, sm, or mi [kilometer, nautical mile, statute mile, and statute mile alias respectively])

### Points
* An object of lat and lons that create the polygon to use for filtering

## Usage

To run (using Linux as an example), run the following. In this example I am using Linux and will be cleaning "Chicago ARTCC Combined_1617077493.sct2" and want to output to cleaned.sct2.

```
./hygieia_linux_amd64 clean Chicago\ ARTCC\ Combined_1617077493.sct2 output.sct2
```

The command syntax is the same, replace the executable as necessary for your OS.

Run the command help for more information.
```
./hygieia_linux_amd64 help
```

## License

This project is licensed under GPL 3. Please see LICENSE.md for more information.